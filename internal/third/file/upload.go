package file

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/log"
)

type UploadFileReq struct {
	Filepath    string `json:"filepath"`
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	Cause       string `json:"cause"`
	Uuid        string `json:"uuid"`
}

type UploadFileResp struct {
	URL string `json:"url"`
}

type partInfo struct {
	ContentType string
	PartSize    int64
	PartNum     int
	FileMd5     string
	PartMd5     string
	PartSizes   []int64
	PartMd5s    []string
}

func NewFile(database db_interface.DataBase, loginUserID string) *File {
	return &File{database: database, loginUserID: loginUserID, confLock: &sync.Mutex{}, mapLocker: &sync.Mutex{}, uploading: make(map[string]*lockInfo)}
}

type File struct {
	database    db_interface.DataBase
	loginUserID string
	confLock    sync.Locker
	partLimit   *third.PartLimitResp
	mapLocker   sync.Locker
	uploading   map[string]*lockInfo
}

type lockInfo struct {
	count  int32
	locker sync.Locker
}

func (f *File) lockHash(hash string) {
	f.mapLocker.Lock()
	locker, ok := f.uploading[hash]
	if !ok {
		locker = &lockInfo{count: 0, locker: &sync.Mutex{}}
		f.uploading[hash] = locker
	}
	atomic.AddInt32(&locker.count, 1)
	f.mapLocker.Unlock()
	locker.locker.Lock()
}

func (f *File) unlockHash(hash string) {
	f.mapLocker.Lock()
	locker, ok := f.uploading[hash]
	if !ok {
		f.mapLocker.Unlock()
		return
	}
	if atomic.AddInt32(&locker.count, -1) == 0 {
		delete(f.uploading, hash)
	}
	f.mapLocker.Unlock()
	locker.locker.Unlock()
}

func (f *File) UploadFile(ctx context.Context, req *UploadFileReq, cb UploadFileCallback) (*UploadFileResp, error) {
	if cb == nil {
		cb = emptyUploadCallback{}
	}
	if req.Name == "" {
		return nil, errors.New("name is empty")
	}
	if req.Name[0] == '/' {
		req.Name = req.Name[1:]
	}
	if prefix := f.loginUserID + "/"; !strings.HasPrefix(req.Name, prefix) {
		req.Name = prefix + req.Name
	}
	file, err := Open(req)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileSize := file.Size()
	cb.Open(fileSize)
	info, err := f.getPartInfo(ctx, file, fileSize, cb)
	if err != nil {
		return nil, err
	}
	if req.ContentType == "" {
		req.ContentType = info.ContentType
	}
	partSize := info.PartSize
	partSizes := info.PartSizes
	partMd5s := info.PartMd5s
	partMd5Val := info.PartMd5
	if err := file.StartSeek(0); err != nil {
		return nil, err
	}
	f.lockHash(partMd5Val)
	defer f.unlockHash(partMd5Val)
	maxParts := 20
	if maxParts > len(partSizes) {
		maxParts = len(partSizes)
	}
	uploadInfo, err := f.getUpload(ctx, &third.InitiateMultipartUploadReq{
		Hash:        partMd5Val,
		Size:        fileSize,
		PartSize:    partSize,
		MaxParts:    int32(maxParts), // retrieve the number of signatures in one go
		Cause:       req.Cause,
		Name:        req.Name,
		ContentType: req.ContentType,
	})
	if err != nil {
		return nil, err
	}
	if uploadInfo.Resp.Upload == nil {
		cb.Complete(fileSize, uploadInfo.Resp.Url, 0)
		return &UploadFileResp{
			URL: uploadInfo.Resp.Url,
		}, nil
	}
	if uploadInfo.Resp.Upload.PartSize != partSize {
		f.cleanPartLimit()
		return nil, fmt.Errorf("part fileSize not match, expect %d, got %d", partSize, uploadInfo.Resp.Upload.PartSize)
	}
	cb.UploadID(uploadInfo.Resp.Upload.UploadID)
	uploadedSize := fileSize
	for i := 0; i < len(partSizes); i++ {
		if !uploadInfo.Bitmap.Get(i) {
			uploadedSize -= partSizes[i]
		}
	}
	continueUpload := uploadedSize > 0
	for i, currentPartSize := range partSizes {
		partNumber := int32(i + 1)
		md5Reader := NewMd5Reader(io.LimitReader(file, currentPartSize))
		if uploadInfo.Bitmap.Get(i) {
			if _, err := io.Copy(io.Discard, md5Reader); err != nil {
				return nil, err
			}
		} else {
			reader := NewProgressReader(md5Reader, func(current int64) {
				cb.UploadComplete(fileSize, uploadedSize+current, uploadedSize)
			})
			urlval, header, err := uploadInfo.GetPartSign(ctx, partNumber)
			if err != nil {
				return nil, err
			}
			if err := f.doPut(ctx, http.DefaultClient, urlval, header, reader, currentPartSize); err != nil {
				log.ZError(ctx, "doPut", err, "partMd5Val", partMd5Val, "name", req.Name, "partNumber", partNumber)
				return nil, err
			}
			uploadedSize += currentPartSize
			if uploadInfo.DBInfo != nil && uploadInfo.Bitmap != nil {
				uploadInfo.Bitmap.Set(i)
				uploadInfo.DBInfo.UploadInfo = base64.StdEncoding.EncodeToString(uploadInfo.Bitmap.Serialize())
				if err := f.database.UpdateUpload(ctx, uploadInfo.DBInfo); err != nil {
					log.ZError(ctx, "SetUploadPartPush", err, "partMd5Val", partMd5Val, "name", req.Name, "partNumber", partNumber)
				}
			}
		}
		md5val := md5Reader.Md5()
		if md5val != partMd5s[i] {
			return nil, fmt.Errorf("upload part %d failed, md5 not match, expect %s, got %s", i, partMd5s[i], md5val)
		}
		cb.UploadPartComplete(i, currentPartSize, partMd5s[i])
		log.ZDebug(ctx, "upload part success", "partMd5Val", md5val, "name", req.Name, "partNumber", partNumber)
	}
	log.ZDebug(ctx, "upload all part success", "partHash", partMd5Val, "name", req.Name)
	resp, err := f.completeMultipartUpload(ctx, &third.CompleteMultipartUploadReq{
		UploadID:    uploadInfo.Resp.Upload.UploadID,
		Parts:       partMd5s,
		Name:        req.Name,
		ContentType: req.ContentType,
		Cause:       req.Cause,
	})
	if err != nil {
		return nil, err
	}
	typ := 1
	if continueUpload {
		typ++
	}
	cb.Complete(fileSize, resp.Url, typ)
	if uploadInfo.DBInfo != nil {
		if err := f.database.DeleteUpload(ctx, info.PartMd5); err != nil {
			log.ZError(ctx, "DeleteUpload", err, "partMd5Val", info.PartMd5, "name", req.Name)
		}
	}
	return &UploadFileResp{
		URL: resp.Url,
	}, nil
}

func (f *File) cleanPartLimit() {
	f.confLock.Lock()
	defer f.confLock.Unlock()
	f.partLimit = nil
}

func (f *File) initiateMultipartUploadResp(ctx context.Context, req *third.InitiateMultipartUploadReq) (*third.InitiateMultipartUploadResp, error) {
	return api.ObjectInitiateMultipartUpload.Invoke(ctx, req)
}

func (f *File) authSign(ctx context.Context, req *third.AuthSignReq) (*third.AuthSignResp, error) {
	if len(req.PartNumbers) == 0 {
		return nil, errs.ErrArgs.WrapMsg("partNumbers is empty")
	}
	return api.ObjectAuthSign.Invoke(ctx, req)
}

func (f *File) completeMultipartUpload(ctx context.Context, req *third.CompleteMultipartUploadReq) (*third.CompleteMultipartUploadResp, error) {
	return api.ObjectCompleteMultipartUpload.Invoke(ctx, req)
}

func (f *File) getObjectPartLimit(ctx context.Context) (*third.PartLimitResp, error) {
	return api.ObjectPartLimit.Invoke(ctx, &third.PartLimitReq{})
}

func (f *File) getPartNum(fileSize int64, partSize int64) int {
	partNum := fileSize / partSize
	if fileSize%partSize != 0 {
		partNum++
	}
	return int(partNum)
}

func (f *File) partSize(ctx context.Context, size int64) (int64, error) {
	f.confLock.Lock()
	defer f.confLock.Unlock()
	if f.partLimit == nil {
		var err error
		f.partLimit, err = f.getObjectPartLimit(ctx)
		if err != nil {
			return 0, err
		}
	}
	if size <= 0 {
		return 0, errors.New("size must be greater than 0")
	}
	if size > f.partLimit.MaxPartSize*int64(f.partLimit.MaxNumSize) {
		return 0, fmt.Errorf("size must be less than %db", f.partLimit.MaxPartSize*int64(f.partLimit.MaxNumSize))
	}
	if size <= f.partLimit.MinPartSize*int64(f.partLimit.MaxNumSize) {
		return f.partLimit.MinPartSize, nil
	}
	partSize := size / int64(f.partLimit.MaxNumSize)
	if size%int64(f.partLimit.MaxNumSize) != 0 {
		partSize++
	}
	return partSize, nil
}

func (f *File) accessURL(ctx context.Context, req *third.AccessURLReq) (*third.AccessURLResp, error) {
	return api.ObjectAccessURL.Invoke(ctx, req)
}

func (f *File) doHttpReq(req *http.Request) ([]byte, *http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return data, resp, nil
}

func (f *File) partMD5(parts []string) string {
	s := strings.Join(parts, ",")
	md5Sum := md5.Sum([]byte(s))
	return hex.EncodeToString(md5Sum[:])
}

type AuthSignParts struct {
	Sign  *third.SignPart
	Times []time.Time
}

type UploadInfo struct {
	PartNum      int
	Bitmap       *Bitmap
	DBInfo       *model_struct.LocalUpload
	Resp         *third.InitiateMultipartUploadResp
	CreateTime   time.Time
	BatchSignNum int32
	f            *File
}

func (u *UploadInfo) getIndex(partNumber int32) int {
	if u.Resp.Upload.Sign == nil {
		return -1
	} else {
		if u.CreateTime.IsZero() {
			return -1
		} else {
			if time.Since(u.CreateTime) > time.Minute {
				return -1
			}
		}
	}
	for i, part := range u.Resp.Upload.Sign.Parts {
		if part.PartNumber == partNumber {
			return i
		}
	}
	return -1
}

func (u *UploadInfo) buildRequest(i int) (*url.URL, http.Header, error) {
	sign := u.Resp.Upload.Sign
	part := sign.Parts[i]
	rawURL := sign.Url
	if part.Url != "" {
		rawURL = part.Url
	}
	urlval, err := url.Parse(rawURL)
	if err != nil {
		return nil, nil, err
	}
	if len(sign.Query)+len(part.Query) > 0 {
		query := urlval.Query()
		for i := range sign.Query {
			v := sign.Query[i]
			query[v.Key] = v.Values
		}
		for i := range part.Query {
			v := part.Query[i]
			query[v.Key] = v.Values
		}
		urlval.RawQuery = query.Encode()
	}
	header := make(http.Header)
	for i := range sign.Header {
		v := sign.Header[i]
		header[v.Key] = v.Values
	}
	for i := range part.Header {
		v := part.Header[i]
		header[v.Key] = v.Values
	}
	return urlval, header, nil
}

func (u *UploadInfo) GetPartSign(ctx context.Context, partNumber int32) (*url.URL, http.Header, error) {
	if partNumber < 1 || int(partNumber) > u.PartNum {
		return nil, nil, errors.New("invalid partNumber")
	}
	if index := u.getIndex(partNumber); index >= 0 {
		return u.buildRequest(index)
	}
	partNumbers := make([]int32, 0, u.BatchSignNum)
	for i := int32(0); i < u.BatchSignNum; i++ {
		if int(partNumber+i) > u.PartNum {
			break
		}
		partNumbers = append(partNumbers, partNumber+i)
	}
	authSignResp, err := u.f.authSign(ctx, &third.AuthSignReq{
		UploadID:    u.Resp.Upload.UploadID,
		PartNumbers: partNumbers,
	})
	if err != nil {
		return nil, nil, err
	}
	u.Resp.Upload.Sign.Url = authSignResp.Url
	u.Resp.Upload.Sign.Query = authSignResp.Query
	u.Resp.Upload.Sign.Header = authSignResp.Header
	u.Resp.Upload.Sign.Parts = authSignResp.Parts
	u.CreateTime = time.Now()
	index := u.getIndex(partNumber)
	if index < 0 {
		return nil, nil, errs.ErrInternalServer.WrapMsg("server part sign invalid")
	}
	return u.buildRequest(index)
}

func (f *File) getLocalUploadInfo(ctx context.Context, req *third.InitiateMultipartUploadReq) (info *UploadInfo) {
	partNum := f.getPartNum(req.Size, req.PartSize)
	if partNum <= 1 {
		return nil
	}
	dbUpload, err := f.database.GetUpload(ctx, req.Hash)
	if err != nil {
		return nil
	}
	defer func() {
		if info == nil {
			if err := f.database.DeleteUpload(ctx, req.Hash); err != nil {
				log.ZError(ctx, "delete upload db", err, "partHash", req.Hash)
			}
		}
	}()
	if dbUpload.UploadID == "" || dbUpload.ExpireTime-3600*1000 < time.Now().UnixMilli() {
		return nil
	}
	bitmapBytes, err := base64.StdEncoding.DecodeString(dbUpload.UploadInfo)
	if err != nil {
		log.ZError(ctx, "decode upload info", err, "partHash", req.Hash)
		return nil
	}
	return &UploadInfo{
		PartNum: partNum,
		Bitmap:  ParseBitmap(bitmapBytes, partNum),
		DBInfo:  dbUpload,
		Resp: &third.InitiateMultipartUploadResp{
			Upload: &third.UploadInfo{
				PartSize:   req.PartSize,
				Sign:       &third.AuthSignParts{},
				UploadID:   dbUpload.UploadID,
				ExpireTime: dbUpload.ExpireTime,
			},
		},
		BatchSignNum: req.MaxParts,
		f:            f,
	}
}

func (f *File) getUpload(ctx context.Context, req *third.InitiateMultipartUploadReq) (*UploadInfo, error) {
	if info := f.getLocalUploadInfo(ctx, req); info != nil {
		return info, nil
	}
	partNum := f.getPartNum(req.Size, req.PartSize)
	resp, err := f.initiateMultipartUploadResp(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Upload == nil {
		return &UploadInfo{
			Resp: resp,
		}, nil
	}
	bitmap := NewBitmap(partNum)
	var dbUpload *model_struct.LocalUpload
	if partNum > 1 {
		dbUpload = &model_struct.LocalUpload{
			PartHash:   req.Hash,
			UploadID:   resp.Upload.UploadID,
			UploadInfo: base64.StdEncoding.EncodeToString(bitmap.Serialize()),
			ExpireTime: resp.Upload.ExpireTime,
			CreateTime: time.Now().UnixMilli(),
		}
		if err := f.database.DeleteUpload(ctx, req.Hash); err != nil {
			log.ZError(ctx, "delete upload db", err, "partHash", req.Hash)
		}
		if err := f.database.InsertUpload(ctx, dbUpload); err != nil {
			log.ZError(ctx, "insert upload db", err, "pratsHash", req.Hash, "name", req.Name)
		}
	}
	return &UploadInfo{
		PartNum:      partNum,
		Bitmap:       bitmap,
		DBInfo:       dbUpload,
		Resp:         resp,
		CreateTime:   time.Now(),
		BatchSignNum: req.MaxParts,
		f:            f,
	}, nil
}

func (f *File) doPut(ctx context.Context, client *http.Client, url *url.URL, header http.Header, reader io.Reader, size int64) error {
	rawURL := url.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, rawURL, reader)
	if err != nil {
		return err
	}
	for key := range header {
		req.Header[key] = header[key]
	}
	req.ContentLength = size
	log.ZDebug(ctx, "do put req", "url", rawURL, "contentLength", size, "header", req.Header)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	log.ZDebug(ctx, "do put resp status", "url", rawURL, "status", resp.Status, "contentLength", resp.ContentLength, "header", resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "do put resp body", "url", rawURL, "body", string(body))
	if resp.StatusCode/200 != 1 {
		return fmt.Errorf("PUT %s failed, status code %d, body %s", rawURL, resp.StatusCode, string(body))
	}
	return nil
}

func (f *File) getPartInfo(ctx context.Context, r io.Reader, fileSize int64, cb UploadFileCallback) (*partInfo, error) {
	partSize, err := f.partSize(ctx, fileSize)
	if err != nil {
		return nil, err
	}
	partNum := int(fileSize / partSize)
	if fileSize%partSize != 0 {
		partNum++
	}
	cb.PartSize(partSize, partNum)
	partSizes := make([]int64, partNum)
	for i := 0; i < partNum; i++ {
		partSizes[i] = partSize
	}
	partSizes[partNum-1] = fileSize - partSize*(int64(partNum)-1)
	partMd5s := make([]string, partNum)
	buf := make([]byte, 1024*8)
	fileMd5 := md5.New()
	var contentType string
	for i := 0; i < partNum; i++ {
		h := md5.New()
		r := io.LimitReader(r, partSize)
		for {
			if n, err := r.Read(buf); err == nil {
				if contentType == "" {
					contentType = http.DetectContentType(buf[:n])
				}
				h.Write(buf[:n])
				fileMd5.Write(buf[:n])
			} else if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		partMd5s[i] = hex.EncodeToString(h.Sum(nil))
		cb.HashPartProgress(i, partSizes[i], partMd5s[i])
	}
	partMd5Val := f.partMD5(partMd5s)
	fileMd5val := hex.EncodeToString(fileMd5.Sum(nil))
	cb.HashPartComplete(f.partMD5(partMd5s), hex.EncodeToString(fileMd5.Sum(nil)))
	return &partInfo{
		ContentType: contentType,
		PartSize:    partSize,
		PartNum:     partNum,
		FileMd5:     fileMd5val,
		PartMd5:     partMd5Val,
		PartSizes:   partSizes,
		PartMd5s:    partMd5s,
	}, nil
}
