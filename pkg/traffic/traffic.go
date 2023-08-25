package traffic

import (
	"sync"
	"time"
)

const (
	normalSize = 1024 * 1024 * 1024 * 5 //2GB

)

type traffic struct {
	lock  sync.Mutex
	Size  int64
	start time.Time
	end   time.Time
}

type Interface interface {
	Add(userId string, size int64) error
	Clear(userId string) error
}

func NewTraffic() Interface {
	return &traffic{}
}

func (t *traffic) Add(userId string, size int64) error {
	return nil
}

func (t *traffic) Clear(userId string) error {
	return nil
}
