package interaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 8) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024

	//Maximum number of reconnection attempts
	maxReconnectAttempts = 300

	sendAndWaitTime = time.Second * 10
)

const (
	DefaultNotConnect = iota
	Closed            = iota + 1
	Connecting
	Connected
)

var (
	ErrChanClosed                = errors.New("send channel closed")
	ErrConnClosed                = errors.New("conn has closed")
	ErrNotSupportMessageProtocol = errors.New("not support message protocol")
	ErrClientClosed              = errors.New("client actively close the connection")
	ErrPanic                     = errors.New("panic error")
)

type LongConnMgr struct {
	//conn status mutex
	w          sync.Mutex
	connStatus int
	// The long connection,can be set tcp or websocket.
	conn       LongConn
	listener   open_im_sdk_callback.OnConnListener
	userOnline func(map[string][]int32)
	// Buffered channel of outbound messages.
	send               chan Message
	pushMsgAndMaxSeqCh chan common.Cmd2Value
	conversationCh     chan common.Cmd2Value
	loginMgrCh         chan common.Cmd2Value
	closedErr          error
	ctx                context.Context
	IsCompression      bool
	Syncer             *WsRespAsyn
	encoder            Encoder
	compressor         Compressor
	reconnectStrategy  ReconnectStrategy

	mutex        sync.Mutex
	IsBackground bool
	// write conn lock
	connWrite *sync.Mutex

	sub *subscription
}

type Message struct {
	Message GeneralWsReq
	Resp    chan *GeneralWsResp
}

func NewLongConnMgr(ctx context.Context, listener open_im_sdk_callback.OnConnListener, userOnline func(map[string][]int32), pushMsgAndMaxSeqCh, loginMgrCh chan common.Cmd2Value) *LongConnMgr {
	l := &LongConnMgr{
		listener:           listener,
		userOnline:         userOnline,
		pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh,
		loginMgrCh:         loginMgrCh,
		IsCompression:      true,
		Syncer:             NewWsRespAsyn(),
		encoder:            NewGobEncoder(),
		compressor:         NewGzipCompressor(),
		reconnectStrategy:  NewExponentialRetry(),
		sub:                newSubscription(),
	}
	l.send = make(chan Message, 10)
	l.conn = NewWebSocket(WebSocket)
	l.connWrite = new(sync.Mutex)
	l.ctx = ctx
	return l
}
func (c *LongConnMgr) Run(ctx context.Context) {
	go c.readPump(ctx)
	go c.writePump(ctx)
	go c.heartbeat(ctx)
}

func (c *LongConnMgr) SendReqWaitResp(ctx context.Context, m proto.Message, reqIdentifier int, resp proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return sdkerrs.ErrArgs
	}
	msg := Message{
		Message: GeneralWsReq{
			ReqIdentifier: reqIdentifier,
			SendID:        ccontext.Info(ctx).UserID(),
			OperationID:   ccontext.Info(ctx).OperationID(),
			Data:          data,
		},
		Resp: make(chan *GeneralWsResp, 1),
	}
	c.send <- msg
	log.ZDebug(ctx, "send message to send channel success", "msg", m, "reqIdentifier", reqIdentifier)
	select {
	case <-ctx.Done():
		return sdkerrs.ErrCtxDeadline
	case v, ok := <-msg.Resp:
		if !ok {
			return errors.New("response channel closed")
		}
		if v.ErrCode != 0 {
			return errs.NewCodeError(v.ErrCode, v.ErrMsg)
		}
		if err := proto.Unmarshal(v.Data, resp); err != nil {
			return sdkerrs.ErrArgs
		}
		return nil
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.

func (c *LongConnMgr) readPump(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())

			log.ZWarn(ctx, "readPump panic", nil, "panic info", err)
		}
	}()

	log.ZDebug(ctx, "readPump start", "goroutine ID:", getGoroutineID())
	defer func() {
		_ = c.close()
		log.ZWarn(c.ctx, "readPump closed", c.closedErr)
	}()
	connNum := 0
	for {
		select {
		case <-ctx.Done():
			c.closedErr = ctx.Err()
			log.ZInfo(c.ctx, "readPump done, sdk logout.....")
			return
		default:
		}
		ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
		needRecon, err := c.reConn(ctx, &connNum)
		if !needRecon {
			c.closedErr = err
			return
		}
		if err != nil {
			log.ZWarn(c.ctx, "reConn", err)
			time.Sleep(c.reconnectStrategy.GetSleepInterval())
			continue
		}
		c.conn.SetReadLimit(maxMessageSize)
		_ = c.conn.SetReadDeadline(pongWait)
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.ZError(c.ctx, "readMessage err", err, "goroutine ID:", getGoroutineID())
			_ = c.close()
			c.sub.onConnClosed(err)
			continue
		}
		switch messageType {
		case MessageBinary:
			err := c.handleMessage(message)
			if err != nil {
				c.closedErr = err
				return
			}
		case MessageText:
			c.closedErr = ErrNotSupportMessageProtocol
			return
		case CloseMessage:
			c.closedErr = ErrClientClosed
			return
		default:
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *LongConnMgr) writePump(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())

			log.ZWarn(ctx, "writePump panic", nil, "panic info", err)
		}
	}()

	log.ZDebug(ctx, "writePump start", "goroutine ID:", getGoroutineID())

	defer func() {
		c.close()
		close(c.send)
	}()
	for {
		select {
		case <-ctx.Done():
			c.closedErr = ctx.Err()
			log.ZInfo(c.ctx, "writePump done, sdk logout.....")
			return
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				_ = c.conn.SetWriteDeadline(writeWait)
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.ZError(c.ctx, "send close message error", err)
				}
				c.closedErr = ErrChanClosed
				return
			}
			log.ZDebug(c.ctx, "writePump recv message", "reqIdentifier", message.Message.ReqIdentifier,
				"operationID", message.Message.OperationID, "sendID", message.Message.SendID)
			resp, err := c.sendAndWaitResp(&message.Message)
			if err != nil {
				resp = &GeneralWsResp{
					ReqIdentifier: message.Message.ReqIdentifier,
					OperationID:   message.Message.OperationID,
					Data:          nil,
				}
				if code, ok := errs.Unwrap(err).(errs.CodeError); ok {
					resp.ErrCode = code.Code()
					resp.ErrMsg = code.Msg()
				} else {
					log.ZError(c.ctx, "writeBinaryMsgAndRetry failed", err, "wsReq", message.Message)
				}

			}
			nErr := c.Syncer.notifyCh(message.Resp, resp, 1)
			if nErr != nil {
				log.ZError(c.ctx, "TriggerCmdNewMsgCome failed", nErr, "wsResp", resp)
			}
		}
	}
}

func (c *LongConnMgr) heartbeat(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Sprintf("panic: %+v\n%s", r, debug.Stack())

			log.ZWarn(ctx, "heartbeat panic", nil, "panic info", err)
		}
	}()

	log.ZDebug(ctx, "heartbeat start", "goroutine ID:", getGoroutineID())
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		log.ZWarn(c.ctx, "heartbeat closed", nil, "heartbeat", "heartbeat done sdk logout.....")
	}()
	for {
		select {
		case <-ctx.Done():
			log.ZInfo(ctx, "heartbeat done sdk logout.....")
			return
		case <-ticker.C:
			log.ZInfo(ctx, "sendPingMessage", "goroutine ID:", getGoroutineID())
			c.sendPingMessage(ctx)
		}
	}

}

func (c *LongConnMgr) sendPingMessage(ctx context.Context) {
	c.connWrite.Lock()
	defer c.connWrite.Unlock()
	opid := utils.OperationIDGenerator()
	log.ZDebug(ctx, "ping Message Started", "goroutine ID:", getGoroutineID(), "opid", opid)
	if c.IsConnected() {
		log.ZDebug(ctx, "ping Message Started isConnected", "goroutine ID:", getGoroutineID(), "opid", opid)
		c.conn.SetWriteDeadline(writeWait)
		if err := c.conn.WriteMessage(PingMessage, []byte(opid)); err != nil {
			log.ZWarn(ctx, "ping Message failed", err, "goroutine ID:", getGoroutineID(), "opid", opid)
			return
		}
	} else {
		log.ZDebug(ctx, "ping Message failed, connection", "connStatus", c.GetConnectionStatus(), "goroutine ID:", getGoroutineID(), "opid", opid)
	}
}

func getGoroutineID() int64 {
	buf := make([]byte, 64)
	buf = buf[:runtime.Stack(buf, false)]
	idField := strings.Fields(strings.TrimPrefix(string(buf), "goroutine "))[0]
	id, err := strconv.ParseInt(idField, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func (c *LongConnMgr) sendAndWaitResp(msg *GeneralWsReq) (*GeneralWsResp, error) {
	tempChan, err := c.writeBinaryMsgAndRetry(msg)
	defer c.Syncer.DelCh(msg.MsgIncr)
	if err != nil {
		return nil, err
	} else {
		select {
		case resp := <-tempChan:
			return resp, nil
		case <-time.After(sendAndWaitTime):
			return nil, sdkerrs.ErrNetworkTimeOut
		}

	}
}

func (c *LongConnMgr) writeBinaryMsgAndRetry(msg *GeneralWsReq) (chan *GeneralWsResp, error) {
	msgIncr, tempChan := c.Syncer.AddCh(msg.SendID)
	msg.MsgIncr = msgIncr
	if c.GetConnectionStatus() != Connected && msg.ReqIdentifier == constant.GetNewestSeq {
		return tempChan, sdkerrs.ErrNetwork.WrapMsg("connection closed,conning...")
	}
	for i := 0; i < maxReconnectAttempts; i++ {
		err := c.writeBinaryMsg(*msg)
		if err != nil {
			log.ZError(c.ctx, "send binary message error", err, "message", msg)
			c.closedErr = err
			_ = c.close()
			time.Sleep(time.Second * 1)
			continue
		} else {
			return tempChan, nil
		}
	}
	return nil, sdkerrs.ErrNetwork.WrapMsg("send binary message error")
}

func (c *LongConnMgr) writeBinaryMsgAndNotRetry(msg *GeneralWsReq) (chan *GeneralWsResp, error) {
	msgIncr, tempChan := c.Syncer.AddCh(msg.SendID)
	msg.MsgIncr = msgIncr
	if err := c.writeBinaryMsg(*msg); err != nil {
		c.Syncer.DelCh(msgIncr)
		return nil, err
	}
	return tempChan, nil
}

func (c *LongConnMgr) writeBinaryMsg(req GeneralWsReq) error {
	c.connWrite.Lock()
	defer c.connWrite.Unlock()
	return c.writeBinaryMsgNoLock(req)
}

func (c *LongConnMgr) writeSubInfo(subscribeUserID, unsubscribeUserID []string, lock bool) error {
	opID := utils.OperationIDGenerator()
	sCtx := ccontext.WithOperationID(c.ctx, opID)
	log.ZInfo(sCtx, "writeSubInfo start", "goroutine ID:", getGoroutineID())
	subReq := sdkws.SubUserOnlineStatus{
		SubscribeUserID:   subscribeUserID,
		UnsubscribeUserID: unsubscribeUserID,
	}
	data, err := proto.Marshal(&subReq)
	if err != nil {
		log.ZError(sCtx, "proto.Marshal", err)
		return err
	}
	req := GeneralWsReq{
		ReqIdentifier: constant.WsSubUserOnlineStatus,
		SendID:        ccontext.Info(sCtx).UserID(),
		OperationID:   opID,
		MsgIncr:       utils.OperationIDGenerator(),
		Data:          data,
	}
	if lock {
		return c.writeBinaryMsg(req)
	} else {
		return c.writeBinaryMsgNoLock(req)
	}
}

func (c *LongConnMgr) writeBinaryMsgNoLock(req GeneralWsReq) error {
	encodeBuf, err := c.encoder.Encode(req)
	if err != nil {
		return err
	}
	if c.GetConnectionStatus() != Connected {
		return sdkerrs.ErrNetwork.WrapMsg("connection closed,re conning...")
	}
	_ = c.conn.SetWriteDeadline(writeWait)
	if c.IsCompression {
		resultBuf, compressErr := c.compressor.CompressWithPool(encodeBuf)
		if compressErr != nil {
			return compressErr
		}
		return c.conn.WriteMessage(MessageBinary, resultBuf)
	} else {
		return c.conn.WriteMessage(MessageBinary, encodeBuf)
	}
}

func (c *LongConnMgr) close() error {
	c.w.Lock()
	defer c.w.Unlock()
	if c.connStatus == Closed || c.connStatus == Connecting || c.connStatus == DefaultNotConnect {
		return nil
	}
	c.connStatus = Closed
	log.ZWarn(c.ctx, "conn closed", c.closedErr)
	return c.conn.Close()
}

func (c *LongConnMgr) handleMessage(message []byte) error {
	if c.IsCompression {
		var decompressErr error
		message, decompressErr = c.compressor.DecompressWithPool(message)
		if decompressErr != nil {
			log.ZError(c.ctx, "DeCompress failed", decompressErr, message)
			return sdkerrs.ErrMsgDeCompression
		}
	}
	var wsResp GeneralWsResp
	err := c.encoder.Decode(message, &wsResp)
	if err != nil {
		log.ZError(c.ctx, "decodeBinaryWs err", err, "message", message)
		return sdkerrs.ErrMsgDecodeBinaryWs
	}
	ctx := context.WithValue(c.ctx, "operationID", wsResp.OperationID)
	log.ZInfo(ctx, "recv msg", "errCode", wsResp.ErrCode, "errMsg", wsResp.ErrMsg,
		"reqIdentifier", wsResp.ReqIdentifier)
	switch wsResp.ReqIdentifier {
	case constant.PushMsg:
		if err = c.doPushMsg(ctx, wsResp); err != nil {
			log.ZError(ctx, "doWSPushMsg failed", err, "wsResp", wsResp)
		}
	case constant.LogoutMsg:
		if err := c.Syncer.NotifyResp(ctx, wsResp); err != nil {
			log.ZError(ctx, "notifyResp failed", err, "wsResp", wsResp)
		}
		return sdkerrs.ErrLoginOut
	case constant.KickOnlineMsg:
		log.ZDebug(ctx, "socket receive client kicked offline")

		err = errs.ErrTokenKicked.WrapMsg("socket receive client kicked offline")
		ccontext.GetApiErrCodeCallback(ctx).OnError(ctx, err)
		return err
	case constant.GetNewestSeq:
		fallthrough
	case constant.PullMsgByRange:
		fallthrough
	case constant.PullMsgBySeqList:
		fallthrough
	case constant.GetConvMaxReadSeq:
		fallthrough
	case constant.PullConvLastMessage:
		fallthrough
	case constant.SendMsg:
		fallthrough
	case constant.SendSignalMsg:
		fallthrough
	case constant.SetBackgroundStatus:
		if err := c.Syncer.NotifyResp(ctx, wsResp); err != nil {
			log.ZError(ctx, "notifyResp failed", err, "reqIdentifier", wsResp.ReqIdentifier, "errCode",
				wsResp.ErrCode, "errMsg", wsResp.ErrMsg, "msgIncr", wsResp.MsgIncr, "operationID", wsResp.OperationID)
		}
	case constant.WsSubUserOnlineStatus:
		if err := c.handlerUserOnlineChange(ctx, wsResp); err != nil {
			log.ZError(ctx, "handlerUserOnlineChange failed", err, "wsResp", wsResp)
		}
	default:
		return sdkerrs.ErrMsgBinaryTypeNotSupport
	}
	return nil
}

func (c *LongConnMgr) handlerUserOnlineChange(ctx context.Context, wsResp GeneralWsResp) error {
	if wsResp.ErrCode != 0 {
		return errs.New("handlerUserOnlineChange failed")
	}
	var tips sdkws.SubUserOnlineStatusTips
	if err := proto.Unmarshal(wsResp.Data, &tips); err != nil {
		return err
	}
	log.ZDebug(ctx, "handlerUserOnlineChange", "tips", &tips)
	c.callbackUserOnlineChange(c.sub.setUserState(tips.Subscribers))
	return nil
}

func (c *LongConnMgr) GetUserOnlinePlatformIDs(ctx context.Context, userIDs []string) (map[string][]int32, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	exist, wait, subUserIDs, unsubUserIDs := c.sub.getUserOnline(userIDs)
	if len(subUserIDs)+len(unsubUserIDs) > 0 {
		if err := c.writeSubInfo(subUserIDs, unsubUserIDs, true); err != nil {
			c.sub.writeFailed(wait, err)
			return nil, err
		}
	}
	for userID, statues := range wait {
		select {
		case <-ctx.Done():
			return nil, context.Cause(ctx)
		case <-statues.Done():
			online, err := statues.Result()
			if err != nil {
				return nil, err
			}
			exist[userID] = online
		}
	}
	return exist, nil
}

func (c *LongConnMgr) UnsubscribeUserOnlinePlatformIDs(ctx context.Context, userIDs []string) error {
	if len(userIDs) > 0 {
		c.sub.unsubscribe(userIDs)
	}
	return nil
}

func (c *LongConnMgr) writeConnFirstSubMsg(ctx context.Context) error {
	userIDs := c.sub.getNewConnSubUserIDs()
	log.ZDebug(ctx, "writeConnFirstSubMsg getNewConnSubUserIDs", "userIDs", userIDs)
	if len(userIDs) == 0 {
		return nil
	}
	if err := c.writeSubInfo(userIDs, nil, false); err != nil {
		c.sub.onConnClosed(err)
		return err
	}
	return nil
}

func (c *LongConnMgr) callbackUserOnlineChange(users map[string][]int32) {
	log.ZDebug(c.ctx, "#### ===> callbackUserOnlineChange", "users", users)
	if len(users) == 0 {
		return
	}
	c.userOnline(users)
	//for userID, onlinePlatformIDs := range users {
	//	status := userPb.OnlineStatus{
	//		UserID:      userID,
	//		PlatformIDs: onlinePlatformIDs,
	//	}
	//	if len(status.PlatformIDs) == 0 {
	//		status.Status = constant.Offline
	//	} else {
	//		status.Status = constant.Online
	//	}
	//	c.userOnline.OnUserStatusChanged(utils.StructToJsonString(users))
	//}
}

func (c *LongConnMgr) IsConnected() bool {
	c.w.Lock()
	defer c.w.Unlock()
	if c.connStatus == Connected {
		return true
	}
	return false

}
func (c *LongConnMgr) GetConnectionStatus() int {
	c.w.Lock()
	defer c.w.Unlock()
	return c.connStatus
}

func (c *LongConnMgr) SetConnectionStatus(status int) {
	c.w.Lock()
	defer c.w.Unlock()
	c.connStatus = status
}

func (c *LongConnMgr) reConn(ctx context.Context, num *int) (needRecon bool, err error) {
	if c.IsConnected() {
		return true, nil
	}
	c.connWrite.Lock()
	defer c.connWrite.Unlock()
	c.listener.OnConnecting()
	c.SetConnectionStatus(Connecting)
	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d&operationID=%s&isBackground=%t",
		ccontext.Info(ctx).WsAddr(), ccontext.Info(ctx).UserID(), ccontext.Info(ctx).Token(),
		ccontext.Info(ctx).PlatformID(), ccontext.Info(ctx).OperationID(), c.GetBackground())
	if c.IsCompression {
		url += fmt.Sprintf("&compression=%s", "gzip")
	}
	log.ZDebug(ctx, "conn start", "url", url)
	resp, err := c.conn.Dial(url, nil)
	if err != nil {
		c.SetConnectionStatus(Closed)
		if resp != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return true, err
			}
			log.ZInfo(ctx, "reConn resp", "body", string(body))
			var apiResp struct {
				ErrCode int    `json:"errCode"`
				ErrMsg  string `json:"errMsg"`
				ErrDlt  string `json:"errDlt"`
			}
			if err := json.Unmarshal(body, &apiResp); err != nil {
				return true, err
			}
			err = errs.NewCodeError(apiResp.ErrCode, apiResp.ErrMsg).WithDetail(apiResp.ErrDlt).Wrap()
			ccontext.GetApiErrCodeCallback(ctx).OnError(ctx, err)
			switch apiResp.ErrCode {
			case
				errs.TokenExpiredError,
				errs.TokenInvalidError,
				errs.TokenMalformedError,
				errs.TokenNotValidYetError,
				errs.TokenUnknownError,
				errs.TokenNotExistError,
				errs.TokenKickedError:
				return false, err
			default:
				return true, err
			}
		}
		c.listener.OnConnectFailed(sdkerrs.NetworkError, err.Error())
		return true, err
	}
	if err := c.writeConnFirstSubMsg(ctx); err != nil {
		log.ZError(ctx, "first write user online sub info error", err)
		ccontext.GetApiErrCodeCallback(ctx).OnError(ctx, err)
		c.listener.OnConnectFailed(sdkerrs.NetworkError, err.Error())
		c.conn.Close()
		return true, err
	}
	c.listener.OnConnectSuccess()
	c.sub.onConnSuccess()
	c.ctx = newContext(c.conn.LocalAddr())
	c.ctx = context.WithValue(ctx, "ConnContext", c.ctx)
	c.SetConnectionStatus(Connected)
	c.conn.SetPongHandler(c.pongHandler)
	c.conn.SetPingHandler(c.pingHandler)
	*num++
	log.ZInfo(c.ctx, "long conn establish success", "localAddr", c.conn.LocalAddr(), "connNum", *num)
	c.reconnectStrategy.Reset()
	_ = common.TriggerCmdConnected(ctx, c.pushMsgAndMaxSeqCh)
	return true, nil
}

func (c *LongConnMgr) doPushMsg(ctx context.Context, wsResp GeneralWsResp) error {
	var msg sdkws.PushMessages
	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		return err
	}
	return common.TriggerCmdPushMsg(ctx, &msg, c.pushMsgAndMaxSeqCh)
}
func (c *LongConnMgr) Close(ctx context.Context) {
	if c.GetConnectionStatus() == Connected {
		log.ZInfo(ctx, "network change conn close")
		c.closedErr = errors.New("closed by client network change")
		_ = c.close()
	} else {
		log.ZInfo(ctx, "conn already closed")
	}

}
func (c *LongConnMgr) GetBackground() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.IsBackground
}
func (c *LongConnMgr) SetBackground(isBackground bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.IsBackground = isBackground
}

// receive ping and send pong.
func (c *LongConnMgr) pingHandler(_ string) error {
	if err := c.conn.SetReadDeadline(pongWait); err != nil {
		return err
	}

	return c.writePongMsg()
}

// when client send pong.
func (c *LongConnMgr) pongHandler(appData string) error {
	log.ZDebug(c.ctx, "server Pong Message Received", "appData", appData)
	if err := c.conn.SetReadDeadline(pongWait); err != nil {
		return err
	}
	return nil
}

func (c *LongConnMgr) writePongMsg() error {
	c.connWrite.Lock()
	defer c.connWrite.Unlock()

	err := c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(PongMessage, nil)
}
