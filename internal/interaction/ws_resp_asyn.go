package interaction

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/tools/log"
)

type GeneralWsResp struct {
	ReqIdentifier int    `json:"reqIdentifier"`
	ErrCode       int    `json:"errCode"`
	ErrMsg        string `json:"errMsg"`
	MsgIncr       string `json:"msgIncr"`
	OperationID   string `json:"operationID"`
	Data          []byte `json:"data"`
}

type GeneralWsReq struct {
	ReqIdentifier int    `json:"reqIdentifier"`
	Token         string `json:"token"`
	SendID        string `json:"sendID"`
	OperationID   string `json:"operationID"`
	MsgIncr       string `json:"msgIncr"`
	Data          []byte `json:"data"`
}

type WsRespAsyn struct {
	wsNotification map[string]chan *GeneralWsResp
	wsMutex        sync.RWMutex
}

func NewWsRespAsyn() *WsRespAsyn {
	return &WsRespAsyn{wsNotification: make(map[string]chan *GeneralWsResp, 10)}
}

func GenMsgIncr(userID string) string {
	return userID + "_" + utils.OperationIDGenerator()
}

func (u *WsRespAsyn) AddCh(userID string) (string, chan *GeneralWsResp) {
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()
	for {
		msgIncr := GenMsgIncr(userID)
		ch := make(chan *GeneralWsResp, 1)
		if _, ok := u.wsNotification[msgIncr]; ok {
			continue
		}
		u.wsNotification[msgIncr] = ch
		return msgIncr, ch
	}
}

func (u *WsRespAsyn) AddChByIncr(msgIncr string) chan *GeneralWsResp {
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()
	ch := make(chan *GeneralWsResp, 1)
	_, ok := u.wsNotification[msgIncr]
	if ok {
		log.ZError(context.Background(), "Repeat failed", nil, msgIncr)
	}
	u.wsNotification[msgIncr] = ch
	return ch
}

func (u *WsRespAsyn) GetCh(msgIncr string) chan *GeneralWsResp {
	ch, ok := u.wsNotification[msgIncr]
	if ok {
		return ch
	}
	return nil
}

func (u *WsRespAsyn) DelCh(msgIncr string) {
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()
	ch, ok := u.wsNotification[msgIncr]
	if ok {
		close(ch)
		delete(u.wsNotification, msgIncr)
	}
}

func (u *WsRespAsyn) notifyCh(ch chan *GeneralWsResp, value *GeneralWsResp, timeout int64) error {
	var flag = 0
	select {
	case ch <- value:
		flag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		return errors.New("send cmd timeout")
	}
}

// write a unit test for this function
func (u *WsRespAsyn) NotifyResp(ctx context.Context, wsResp GeneralWsResp) error {
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()

	ch := u.GetCh(wsResp.MsgIncr)
	if ch == nil {
		return errs.WrapMsg(errors.New("no ch"), "GetCh failed "+wsResp.MsgIncr)
	}
	for {
		err := u.notifyCh(ch, &wsResp, 1)
		if err != nil {
			log.ZWarn(ctx, "TriggerCmdNewMsgCome failed ", err, "ch", ch, "wsResp", wsResp)
			continue

		}
		return nil
	}
}
func (u *WsRespAsyn) WaitResp(ctx context.Context, ch chan *GeneralWsResp, timeout int) (*GeneralWsResp, error) {
	select {
	case r, ok := <-ch:
		if !ok { //ch has been closed
			//log.Debug(operationID, "ws network has been changed ")
			return nil, nil
		}
		//log.Debug(operationID, "ws ch recvMsg success, code ", r.ErrCode)
		if r.ErrCode != 0 {
			//log.Error(operationID, "ws ch recvMsg failed, code, err msg: ", r.ErrCode, r.ErrMsg)
			//switch r.ErrCode {
			//case int(constant.ErrInBlackList.ErrCode):
			//	return nil, &constant.ErrInBlackList
			//case int(constant.ErrNotFriend.ErrCode):
			//	return nil, &constant.ErrNotFriend
			//}
			//return nil, errors.New(utils.IntToString(r.ErrCode) + ":" + r.ErrMsg)
		} else {
			return r, nil
		}

	case <-time.After(time.Second * time.Duration(timeout)):
		//log.Error(operationID, "ws ch recvMsg err, timeout")
		//if w.conn.IsNil() {
		//	return nil, errors.New("ws ch recvMsg err, timeout,conn is nil")
		//}
		//if w.conn.CheckSendConnDiffNow() {
		//	return nil, constant.WsRecvConnDiff
		//} else {
		//	return nil, constant.WsRecvConnSame
		//}
	}
	return nil, nil
}
