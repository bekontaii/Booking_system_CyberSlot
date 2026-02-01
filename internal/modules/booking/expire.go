package booking

import "time"

func AutoExpire() {
	go func() {
		time.Sleep(10 * time.Second)
	}()
}
