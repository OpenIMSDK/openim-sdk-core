package conversation_msg

import "sync"

type MaxSeqRecorder struct {
	seqs map[string]int64
	lock sync.RWMutex
}

func NewMaxSeqRecorder() MaxSeqRecorder {
	m := make(map[string]int64)
	return MaxSeqRecorder{seqs: m}
}

func (m *MaxSeqRecorder) Get(conversationID string) int64 {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.seqs[conversationID]
}

func (m *MaxSeqRecorder) Set(conversationID string, seq int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.seqs[conversationID] = seq
}

func (m *MaxSeqRecorder) Incr(conversationID string, num int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.seqs[conversationID] += num
}

func (m *MaxSeqRecorder) IsNewMsg(conversationID string, seq int64) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	currentSeq := m.seqs[conversationID]
	return seq > currentSeq
}
