package relay

import "time"

// alive upload node status every 2 minute
func alive() {
	timer := time.NewTimer(time.Minute * 2)

	go func() {
		upload()
		timer.Reset(time.Minute * 2)
	}()
}

func upload() {

}
