package interaction

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"strings"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

const writeTimeoutSeconds = 30

type WsConn struct {
	stateMutex     sync.Mutex
	conn           *websocket.Conn
	loginStatus    int32
	listener       open_im_sdk_callback.OnConnListener
	token          string
	loginUserID    string
	IsCompression  bool
	ConversationCh chan common.Cmd2Value
	IsConnected    bool
}

func NewWsConn(listener open_im_sdk_callback.OnConnListener, token string, loginUserID string, isCompression bool, conversationCh chan common.Cmd2Value) *WsConn {
	p := WsConn{listener: listener, token: token, loginUserID: loginUserID, IsCompression: isCompression, ConversationCh: conversationCh}
	//	go func() {
	p.conn, _, _ = p.ReConn("init:" + utils.OperationIDGenerator())
	//	}()
	return &p
}

func (u *WsConn) CloseConn(operationID string) error {
	u.Lock()
	defer u.Unlock()
	if u.conn != nil {
		err := u.conn.Close(websocket.StatusNormalClosure, "Actively close the connection")
		log.NewWarn(operationID, "close conn, ", u.conn)
		//u.conn = nil
		return utils.Wrap(err, "")
	}
	return nil
}

func (u *WsConn) LoginStatus() int32 {
	return u.loginStatus
}

func (u *WsConn) SetLoginStatus(loginState int32) {
	u.loginStatus = loginState
}

func (u *WsConn) Lock() {
	u.stateMutex.Lock()
}

func (u *WsConn) Unlock() {
	u.stateMutex.Unlock()
}

func (u *WsConn) SendPingMsg() error {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	if u.conn == nil {
		return utils.Wrap(errors.New("conn == nil"), "")
	}
	//ping := "try ping"
	err := u.SetWriteTimeout(writeTimeoutSeconds)
	if err != nil {
		return utils.Wrap(err, "SetWriteDeadline failed")
	}
	err = u.conn.Ping(context.Background())
	if err != nil {
		return utils.Wrap(err, "WriteMessage failed")
	}
	return nil
}

func (u *WsConn) SetWriteTimeout(timeout int) error {
	//return u.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	return nil
}

func (u *WsConn) SetReadTimeout(timeout int) {
	//u.conn.SetReadLimit(time.Now().Add(time.Duration(timeout) * time.Second).Unix())
}

func (u *WsConn) writeBinaryMsg(msg GeneralWsReq) (*websocket.Conn, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(msg)
	if err != nil {
		return nil, utils.Wrap(err, "Encode error")
	}

	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	if u.conn != nil {
		err := u.SetWriteTimeout(writeTimeoutSeconds)
		if err != nil {
			return nil, utils.Wrap(err, "SetWriteTimeout")
		}
		log.Debug("this msg length is :", float32(len(buff.Bytes()))/float32(1024), "kb")
		if len(buff.Bytes()) > constant.MaxTotalMsgLen {
			return nil, utils.Wrap(errors.New("msg too long"), utils.IntToString(len(buff.Bytes())))
		}
		var data []byte
		if u.IsCompression {
			var gzipBuffer bytes.Buffer
			gz := gzip.NewWriter(&gzipBuffer)
			if _, err := gz.Write(buff.Bytes()); err != nil {
				return nil, utils.Wrap(err, "")
			}
			if err := gz.Close(); err != nil {
				return nil, utils.Wrap(err, "")
			}
			data = gzipBuffer.Bytes()
		} else {
			data = buff.Bytes()
		}
		return u.conn, utils.Wrap(u.conn.Write(context.Background(), websocket.MessageBinary, data), "")
	} else {
		return nil, utils.Wrap(errors.New("conn==nil"), "")
	}
}

func (u *WsConn) decodeBinaryWs(message []byte) (*GeneralWsResp, error) {
	buff := bytes.NewBuffer(message)
	dec := gob.NewDecoder(buff)
	var data GeneralWsResp
	err := dec.Decode(&data)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return &data, nil
}

func (u *WsConn) IsReadTimeout(err error) bool {
	if strings.Contains(err.Error(), "timeout") {
		return true
	}
	return false
}

func (u *WsConn) IsWriteTimeout(err error) bool {
	if strings.Contains(err.Error(), "timeout") {
		return true
	}
	return false
}

func (u *WsConn) IsFatalError(err error) bool {
	if strings.Contains(err.Error(), "timeout") {
		return false
	}
	return true
}

func (u *WsConn) ReConn(operationID string) (*websocket.Conn, error, bool) {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	if u.conn != nil {
		log.NewWarn(operationID, "close conn, ", u.conn)
		_ = u.conn.Close(websocket.StatusNormalClosure, "Actively close the connection")
		u.conn = nil
	}
	if u.loginStatus == constant.TokenFailedKickedOffline {
		return nil, utils.Wrap(errors.New("don't re conn"), "TokenFailedKickedOffline"), false
	}
	u.listener.OnConnecting()

	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d&operationID=%s", sdk_struct.SvrConf.WsAddr, u.loginUserID, u.token, sdk_struct.SvrConf.Platform, operationID)
	log.Info(operationID, "ws connect begin, dail: ", url)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if u.IsCompression {
		url += fmt.Sprintf("&compression=%s", "gzip")
	}
	log.Debug(operationID, "last url:", url, u.IsCompression)
	conn, httpResp, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Error(operationID, "ws connect failed ", url, err.Error())
		u.loginStatus = constant.LoginFailed
		if httpResp != nil {
			errMsg := httpResp.Header.Get("ws_err_msg") + " operationID " + operationID + err.Error()
			log.Error(operationID, "websocket.DefaultDialer.Dial failed ", errMsg, httpResp.StatusCode)
			u.listener.OnConnectFailed(int32(httpResp.StatusCode), errMsg)
			u.IsConnected = false

			switch int32(httpResp.StatusCode) {
			case constant.ErrTokenExpired.ErrCode:
				u.listener.OnUserTokenExpired()
				return nil, utils.Wrap(err, errMsg), false
			case constant.ErrTokenInvalid.ErrCode:
				return nil, utils.Wrap(err, errMsg), false
			case constant.ErrTokenMalformed.ErrCode:
				return nil, utils.Wrap(err, errMsg), false
			case constant.ErrTokenNotValidYet.ErrCode:
				return nil, utils.Wrap(err, errMsg), false
			case constant.ErrTokenUnknown.ErrCode:
				return nil, utils.Wrap(err, errMsg), false
			case constant.ErrTokenDifferentPlatformID.ErrCode:
				return nil, utils.Wrap(err, errMsg), false
			case constant.ErrTokenDifferentUserID.ErrCode:
				return nil, utils.Wrap(err, errMsg), false
			case constant.ErrTokenKicked.ErrCode:
				u.listener.OnKickedOffline()
				return nil, utils.Wrap(err, errMsg), false
			default:
				errMsg = err.Error() + " operationID " + operationID
				u.listener.OnConnectFailed(1001, errMsg)
				u.IsConnected = false

				return nil, utils.Wrap(err, errMsg), true
			}
		} else {
			errMsg := err.Error() + " operationID " + operationID
			u.listener.OnConnectFailed(1001, errMsg)
			u.IsConnected = false
			if u.ConversationCh != nil {
				common.TriggerCmdSuperGroupMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: nil, OperationID: operationID, SyncFlag: constant.MsgSyncBegin}, u.ConversationCh)
				common.TriggerCmdSuperGroupMsgCome(sdk_struct.CmdNewMsgComeToConversation{MsgList: nil, OperationID: operationID, SyncFlag: constant.MsgSyncFailed}, u.ConversationCh)
			}
			log.Error(operationID, "websocket.DefaultDialer.Dial failed ", errMsg, "url ", url)
			return nil, utils.Wrap(err, errMsg), true
		}
	}
	log.Info(operationID, "ws connect end, dail : ", url, conn)
	u.listener.OnConnectSuccess()
	u.IsConnected = true
	u.loginStatus = constant.LoginSuccess
	u.conn = conn
	u.conn.SetReadLimit(1024 * 1024 * 30)
	return conn, nil, true
}
