package util

import "time"

func AddJob(networkId string, fn func(networkId string) error) {

	timer := time.NewTimer(time.Second * 30)
	for {
		timer.Reset(time.Second * 30)
		select {
		case <-timer.C:
			fn(networkId)
		}
	}
}
