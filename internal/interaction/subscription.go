package interaction

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"unsafe"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/utils/datautil"
)

type subscriptionStatues struct {
	done   chan struct{}
	err    error
	online []int32
}

func (s *subscriptionStatues) finish(online []int32, err error) {
	select {
	case <-s.done:
		s.online = online
		s.err = err
	default:
		s.online = online
		s.err = err
		close(s.done)
	}
}

func (s *subscriptionStatues) Done() <-chan struct{} {
	return s.done
}

func (s *subscriptionStatues) Result() ([]int32, error) {
	return s.online, s.err
}

func newSubscription() *subscription {
	return &subscription{
		load:  make(map[string]*subscriptionStatues),
		unsub: make(map[string]struct{}),
		sub:   make(map[string]struct{}),
	}
}

type subscription struct {
	lock  sync.Mutex
	load  map[string]*subscriptionStatues
	unsub map[string]struct{}
	sub   map[string]struct{}
}

func (s *subscription) getNewConnSubUserIDs() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return datautil.Keys(s.sub)
}

func (s *subscription) onConnClosed(err error) {
	if err == nil {
		err = fmt.Errorf("connection closed")
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	clear(s.unsub)
	for userID, statues := range s.load {
		statues.finish(nil, err)
		delete(s.load, userID)
	}
}

func (s *subscription) onConnSuccess() {
	s.lock.Lock()
	defer s.lock.Unlock()
	clear(s.unsub)
}

func (s *subscription) setUserState(changes []*sdkws.SubUserOnlineStatusElem) map[string][]int32 {
	if len(changes) == 0 {
		return nil
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	change := make(map[string][]int32)
	for _, v := range changes {
		if v.OnlinePlatformIDs == nil {
			v.OnlinePlatformIDs = []int32{}
		}
		if status, ok := s.load[v.UserID]; ok {
			delete(s.unsub, v.UserID)
			if !slices.Equal(status.online, v.OnlinePlatformIDs) {
				change[v.UserID] = v.OnlinePlatformIDs
			}
			status.finish(v.OnlinePlatformIDs, nil)
		} else {
			if _, ok := s.sub[v.UserID]; ok {
				done := make(chan struct{})
				s.load[v.UserID] = &subscriptionStatues{
					done:   done,
					err:    nil,
					online: v.OnlinePlatformIDs,
				}
				change[v.UserID] = v.OnlinePlatformIDs
			} else {
				s.unsub[v.UserID] = struct{}{}
			}
		}
	}
	return change
}

func (s *subscription) unsubscribe(userIDs []string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, userID := range userIDs {
		if _, ok := s.sub[userID]; ok {
			delete(s.sub, userID)
			s.unsub[userID] = struct{}{}
		}
		if status, ok := s.load[userID]; ok {
			status.finish(nil, errors.New("unsubscribe"))
			delete(s.load, userID)
		}
	}
}

func (s *subscription) getUserOnline(userIDs []string) (map[string][]int32, map[string]*subscriptionWait, []string, []string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(userIDs) == 0 {
		userIDs = datautil.Keys(s.sub)
	}
	tmp := make(map[string]struct{})
	exist := make(map[string][]int32)
	wait := make(map[string]*subscriptionWait)
	subUserIDs := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		if _, ok := tmp[userID]; ok {
			continue
		}
		tmp[userID] = struct{}{}
		delete(s.unsub, userID)
		if status, ok := s.load[userID]; ok {
			select {
			case <-status.Done():
				exist[userID] = status.online
			default:
				wait[userID] = &subscriptionWait{status: status, first: false}
			}
		} else {
			delete(s.unsub, userID)
			status := &subscriptionStatues{
				done:   make(chan struct{}),
				online: nil,
			}
			s.load[userID] = status
			wait[userID] = &subscriptionWait{status: status, first: true}
			subUserIDs = append(subUserIDs, userID)
		}
	}
	if len(subUserIDs) == 0 {
		return exist, wait, nil, nil
	}
	defer clear(s.unsub)
	return exist, wait, subUserIDs, datautil.Keys(s.unsub)
}

func (s *subscription) writeFailed(waits map[string]*subscriptionWait, err error) {
	if len(waits) == 0 {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	for userID, wait := range waits {
		if !wait.first {
			continue
		}
		if status, ok := s.load[userID]; ok && uintptr(unsafe.Pointer(status)) == uintptr(unsafe.Pointer(wait.status)) {
			delete(s.load, userID)
		}
		wait.status.finish(nil, err)
	}
}

type subscriptionWait struct {
	status *subscriptionStatues
	first  bool
}

func (s *subscriptionWait) Done() <-chan struct{} {
	return s.status.done
}

func (s *subscriptionWait) Result() ([]int32, error) {
	return s.status.online, s.status.err
}
