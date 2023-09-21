package device

import (
	"fmt"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {

	p := NewPool()

	ch := make(chan *Frame, 1000)

	go func() {
		i := 0
		for {
			buffPtr := p.Get()

			buff := *buffPtr

			f := NewFrame()
			f.Buff = buff

			str := fmt.Sprintf("%s%d", "hello, world", i)
			copy(f.Buff[:], str)

			ch <- f

			time.Sleep(5 * time.Second)
			i++
			p.Put(buffPtr)

		}
	}()

	go func() {
		for {
			pkt := <-ch
			fmt.Println("Println: ", pkt.Buff)
			//p.Put(&pkt.Buff)
		}
	}()

	time.Sleep(1000 * time.Minute)
}
