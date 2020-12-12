package implementation

import (
	"sync"
	"time"
)

type activity struct {
	currentTime time.Time
	count       int
}

type senderActivity struct {
	users           map[string]*activity
	maxMsgFrequency int
	sync.Mutex
}

func newSenderActivity(maxMsgFrequency int) *senderActivity {
	return &senderActivity{
		users:           make(map[string]*activity),
		maxMsgFrequency: maxMsgFrequency,
	}
}

func (s *senderActivity) DoesUserLadyGaga(userID string, messageCount int) bool {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UTC()
	if val, exist := s.users[userID]; exist {

		diff := now.Sub(val.currentTime)

		if int(diff.Minutes()) == 0 {
			val.count += messageCount
		} else {
			val.currentTime = now
			val.count = messageCount
		}
	} else {
		s.users[userID] = &activity{
			currentTime: now,
			count:       messageCount,
		}
	}

	if s.users[userID].count >= s.maxMsgFrequency {
		delete(s.users, userID)

		return true
	}

	return false
}
