package conversation_msg

import (
	"context"
	"errors"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/tools/log"

	"os"
	"path/filepath"
	"strings"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

func (c *Conversation) CreateTextMessage(ctx context.Context, text string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Text)
	if err != nil {
		return nil, err
	}
	s.TextElem = &sdk_struct.TextElem{Content: text}
	return &s, nil
}

func (c *Conversation) CreateAdvancedTextMessage(ctx context.Context, text string, messageEntities []*sdk_struct.MessageEntity) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.AdvancedText)
	if err != nil {
		return nil, err
	}
	s.AdvancedTextElem = &sdk_struct.AdvancedTextElem{
		Text:              text,
		MessageEntityList: messageEntities,
	}
	return &s, nil
}

func (c *Conversation) CreateTextAtMessage(ctx context.Context, text string, userIDList []string, usersInfo []*sdk_struct.AtInfo, qs *sdk_struct.MsgStruct) (*sdk_struct.MsgStruct, error) {
	if text == "" {
		return nil, errors.New("text can not be empty")
	}
	if len(userIDList) > 10 {
		return nil, sdkerrs.ErrArgs
	}
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.AtText)
	if err != nil {
		return nil, err
	}
	//Avoid nested references
	if qs != nil {
		if qs.ContentType == constant.Quote {
			qs.ContentType = constant.Text
			qs.TextElem = &sdk_struct.TextElem{Content: qs.QuoteElem.Text}
		}
	}
	s.AtTextElem = &sdk_struct.AtTextElem{
		Text:         text,
		AtUserList:   userIDList,
		AtUsersInfo:  usersInfo,
		QuoteMessage: qs,
	}
	return &s, nil
}

func (c *Conversation) CreateLocationMessage(ctx context.Context, description string, longitude, latitude float64) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Location)
	if err != nil {
		return nil, err
	}
	s.LocationElem = &sdk_struct.LocationElem{
		Description: description,
		Longitude:   longitude,
		Latitude:    latitude,
	}
	return &s, nil
}

func (c *Conversation) CreateCustomMessage(ctx context.Context, data, extension string, description string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Custom)
	if err != nil {
		return nil, err
	}
	s.CustomElem = &sdk_struct.CustomElem{
		Data:        data,
		Extension:   extension,
		Description: description,
	}
	return &s, nil
}

func (c *Conversation) CreateQuoteMessage(ctx context.Context, text string, qs *sdk_struct.MsgStruct) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Quote)
	if err != nil {
		return nil, err
	}
	//Avoid nested references
	if qs.ContentType == constant.Quote {
		qs.ContentType = constant.Text
		qs.TextElem = &sdk_struct.TextElem{Content: qs.QuoteElem.Text}
	}
	s.QuoteElem = &sdk_struct.QuoteElem{
		Text:         text,
		QuoteMessage: qs,
	}
	return &s, nil
}

func (c *Conversation) CreateAdvancedQuoteMessage(ctx context.Context, text string, qs *sdk_struct.MsgStruct, messageEntities []*sdk_struct.MessageEntity) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Quote)
	if err != nil {
		return nil, err
	}
	//Avoid nested references
	if qs.ContentType == constant.Quote {
		//qs.Content = qs.QuoteElem.Text
		qs.ContentType = constant.Text
		qs.TextElem = &sdk_struct.TextElem{Content: qs.QuoteElem.Text}
	}
	s.QuoteElem = &sdk_struct.QuoteElem{
		Text:              text,
		QuoteMessage:      qs,
		MessageEntityList: messageEntities,
	}
	return &s, nil
}

func (c *Conversation) CreateCardMessage(ctx context.Context, card *sdk_struct.CardElem) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Card)
	if err != nil {
		return nil, err
	}
	s.CardElem = card
	return &s, nil
}

func (c *Conversation) CreateVideoMessageFromFullPath(ctx context.Context, videoFullPath string, videoType string, duration int64, snapshotFullPath string) (*sdk_struct.MsgStruct, error) {
	dstFile := utils.FileTmpPath(videoFullPath, c.DataDir) //a->b
	written, err := utils.CopyFile(videoFullPath, dstFile)
	if err != nil {
		//log.Error("internal", "open file failed: ", err, videoFullPath)
		return nil, err
	}
	log.ZDebug(ctx, "videoFullPath dstFile", "videoFullPath", videoFullPath,
		"dstFile", dstFile, "written", written)

	dstFile = utils.FileTmpPath(snapshotFullPath, c.DataDir) //a->b
	sWritten, err := utils.CopyFile(snapshotFullPath, dstFile)
	if err != nil {
		//log.Error("internal", "open file failed: ", err, snapshotFullPath)
		return nil, err
	}
	log.ZDebug(ctx, "snapshotFullPath dstFile", "snapshotFullPath", snapshotFullPath,
		"dstFile", dstFile, "sWritten", sWritten)

	s := sdk_struct.MsgStruct{}
	err = c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Video)
	if err != nil {
		return nil, err
	}
	s.VideoElem = &sdk_struct.VideoElem{
		VideoPath: videoFullPath,
		VideoType: videoType,
		Duration:  duration,
	}
	if snapshotFullPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = snapshotFullPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		//log.Error("internal", "get file Attributes error", err.Error())
		return nil, err
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotFullPath != "" {
		imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			log.ZError(ctx, "getImageInfo err:", err, "snapshotFullPath", snapshotFullPath)
			return nil, err
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	return &s, nil
}

func (c *Conversation) CreateFileMessageFromFullPath(ctx context.Context, fileFullPath string, fileName string) (*sdk_struct.MsgStruct, error) {
	dstFile := utils.FileTmpPath(fileFullPath, c.DataDir)
	_, err := utils.CopyFile(fileFullPath, dstFile)
	if err != nil {
		//log.Error("internal", "open file failed: ", err.Error(), fileFullPath)
		return nil, err

	}
	s := sdk_struct.MsgStruct{}
	err = c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.File)
	if err != nil {
		return nil, err
	}
	fi, err := os.Stat(fileFullPath)
	if err != nil {
		//log.Error("internal", "get file Attributes error", err.Error())
		return nil, err
	}
	s.FileElem = &sdk_struct.FileElem{
		FilePath: fileFullPath,
		FileName: fileName,
		FileSize: fi.Size(),
	}
	return &s, nil
}

func (c *Conversation) CreateImageMessageFromFullPath(ctx context.Context, imageFullPath string) (*sdk_struct.MsgStruct, error) {
	dstFile := utils.FileTmpPath(imageFullPath, c.DataDir) //a->b
	_, err := utils.CopyFile(imageFullPath, dstFile)
	if err != nil {
		//log.Error(operationID, "open file failed: ", err, imageFullPath)
		return nil, err
	}
	s := sdk_struct.MsgStruct{}
	err = c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Picture)
	if err != nil {
		return nil, err
	}
	imageInfo, err := getImageInfo(imageFullPath)
	if err != nil {
		//log.Error(operationID, "getImageInfo err:", err.Error())
		return nil, err
	}
	s.PictureElem = &sdk_struct.PictureElem{
		SourcePath: imageFullPath,
		SourcePicture: &sdk_struct.PictureBaseInfo{
			Width:  imageInfo.Width,
			Height: imageInfo.Height,
			Type:   imageInfo.Type,
		},
	}
	return &s, nil
}

func (c *Conversation) CreateSoundMessageFromFullPath(ctx context.Context, soundPath string, duration int64) (*sdk_struct.MsgStruct, error) {
	dstFile := utils.FileTmpPath(soundPath, c.DataDir) //a->b
	_, err := utils.CopyFile(soundPath, dstFile)
	if err != nil {
		//log.Error("internal", "open file failed: ", err, soundPath)
		return nil, err
	}

	s := sdk_struct.MsgStruct{}
	err = c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Sound)
	if err != nil {
		return nil, err
	}
	fi, err := os.Stat(soundPath)
	if err != nil {
		//log.Error("internal", "getSoundInfo err:", err.Error(), s.SoundElem.SoundPath)
		return nil, err
	}
	s.SoundElem = &sdk_struct.SoundElem{
		SoundPath: soundPath,
		Duration:  duration,
		DataSize:  fi.Size(),
		SoundType: strings.Replace(filepath.Ext(fi.Name()), ".", "", 1),
	}
	return &s, nil
}

func (c *Conversation) CreateImageMessage(ctx context.Context, imagePath string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Picture)
	if err != nil {
		return nil, err
	}
	path := c.DataDir + imagePath
	//path := imagePath
	imageInfo, err := getImageInfo(path)
	if err != nil {
		//log.Error("internal", "get imageInfo err", err.Error())
		return nil, err
	}
	s.PictureElem = &sdk_struct.PictureElem{
		SourcePath: path,
		SourcePicture: &sdk_struct.PictureBaseInfo{
			Width:  imageInfo.Width,
			Height: imageInfo.Height,
			Type:   imageInfo.Type,
		},
	}
	return &s, nil
}

func (c *Conversation) CreateImageMessageByURL(ctx context.Context, sourcePath string, sourcePicture, bigPicture, snapshotPicture sdk_struct.PictureBaseInfo) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Picture)
	if err != nil {
		return nil, err
	}
	s.PictureElem = &sdk_struct.PictureElem{
		SourcePath:      sourcePath,
		SourcePicture:   &sourcePicture,
		BigPicture:      &bigPicture,
		SnapshotPicture: &snapshotPicture,
	}
	return &s, nil
}

func (c *Conversation) CreateSoundMessageByURL(ctx context.Context, soundElem *sdk_struct.SoundBaseInfo) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Sound)
	if err != nil {
		return nil, err
	}
	s.SoundElem = &sdk_struct.SoundElem{
		UUID:      soundElem.UUID,
		SoundPath: soundElem.SoundPath,
		SourceURL: soundElem.SourceURL,
		DataSize:  soundElem.DataSize,
		Duration:  soundElem.Duration,
		SoundType: soundElem.SoundType,
	}
	return &s, nil
}

func (c *Conversation) CreateSoundMessage(ctx context.Context, soundPath string, duration int64) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Sound)
	if err != nil {
		return nil, err
	}
	path := c.DataDir + soundPath
	fi, err := os.Stat(path)
	if err != nil {
		//log.Error("internal", "get sound info err", err.Error())
		return nil, err
	}
	s.SoundElem = &sdk_struct.SoundElem{
		SoundPath: path,
		Duration:  duration,
		DataSize:  fi.Size(),
	}
	if typ := strings.Replace(filepath.Ext(fi.Name()), ".", "", 1); typ != "" {
		s.SoundElem.SoundType = "audio/" + strings.ToLower(typ)
	}
	return &s, nil
}

func (c *Conversation) CreateVideoMessageByURL(ctx context.Context, videoElem sdk_struct.VideoBaseInfo) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Video)
	if err != nil {
		return nil, err
	}
	s.VideoElem = &sdk_struct.VideoElem{
		VideoPath:      videoElem.VideoPath,
		VideoUUID:      videoElem.VideoUUID,
		VideoURL:       videoElem.VideoURL,
		VideoType:      videoElem.VideoType,
		VideoSize:      videoElem.VideoSize,
		Duration:       videoElem.Duration,
		SnapshotPath:   videoElem.SnapshotPath,
		SnapshotUUID:   videoElem.SnapshotUUID,
		SnapshotSize:   videoElem.SnapshotSize,
		SnapshotURL:    videoElem.SnapshotURL,
		SnapshotWidth:  videoElem.SnapshotWidth,
		SnapshotHeight: videoElem.SnapshotHeight,
		SnapshotType:   videoElem.SnapshotType,
	}
	return &s, nil
}

func (c *Conversation) CreateVideoMessage(ctx context.Context, videoPath string, videoType string, duration int64, snapshotPath string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Video)
	if err != nil {
		return nil, err
	}
	s.VideoElem = &sdk_struct.VideoElem{}
	s.VideoElem.VideoPath = c.DataDir + videoPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = c.DataDir + snapshotPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		log.ZDebug(ctx, "get video file error", "videoPath", videoPath, "snapshotPath", snapshotPath)
		return nil, err
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotPath != "" {
		imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			//log.Error("internal", "get snapshot info ", err.Error())
			return nil, err
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	return &s, nil
}

func (c *Conversation) CreateFileMessageByURL(ctx context.Context, fileElem sdk_struct.FileBaseInfo) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.File)
	if err != nil {
		return nil, err
	}
	s.FileElem = &sdk_struct.FileElem{
		FilePath:  fileElem.FilePath,
		UUID:      fileElem.UUID,
		SourceURL: fileElem.SourceURL,
		FileName:  fileElem.FileName,
		FileSize:  fileElem.FileSize,
		FileType:  fileElem.FileType,
	}
	return &s, nil
}

func (c *Conversation) CreateFileMessage(ctx context.Context, filePath string, fileName string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{FileElem: &sdk_struct.FileElem{}}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.File)
	if err != nil {
		return nil, err
	}
	s.FileElem.FilePath = c.DataDir + filePath
	s.FileElem.FileName = fileName
	fi, err := os.Stat(s.FileElem.FilePath)
	if err != nil {
		//log.Error("internal", "get file message err", err.Error())
		return nil, err
	}
	s.FileElem.FileSize = fi.Size()
	s.Content = utils.StructToJsonString(s.FileElem)
	return &s, nil
}

func (c *Conversation) CreateMergerMessage(ctx context.Context, messages []*sdk_struct.MsgStruct, title string, summaries []string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{MergeElem: &sdk_struct.MergeElem{}}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Merger)
	if err != nil {
		return nil, err
	}
	s.MergeElem.AbstractList = summaries
	s.MergeElem.Title = title
	s.MergeElem.MultiMessage = messages
	s.Content = utils.StructToJsonString(s.MergeElem)
	return &s, nil
}

func (c *Conversation) CreateFaceMessage(ctx context.Context, index int, data string) (*sdk_struct.MsgStruct, error) {
	s := sdk_struct.MsgStruct{FaceElem: &sdk_struct.FaceElem{}}
	err := c.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Face)
	if err != nil {
		return nil, err
	}
	s.FaceElem.Data = data
	s.FaceElem.Index = index
	s.Content = utils.StructToJsonString(s.FaceElem)
	return &s, nil
}

func (c *Conversation) CreateForwardMessage(ctx context.Context, s *sdk_struct.MsgStruct) (*sdk_struct.MsgStruct, error) {
	if s.Status != constant.MsgStatusSendSuccess {
		log.ZError(ctx, "only send success message can be Forward",
			errors.New("only send success message can be Forward"))
		return nil, errors.New("only send success message can be Forward")
	}
	err := c.initBasicInfo(ctx, s, constant.UserMsgType, s.ContentType)
	if err != nil {
		return nil, err
	}
	//Forward message seq is set to 0
	s.Seq = 0
	s.Status = constant.MsgStatusSendSuccess
	return s, nil
}
