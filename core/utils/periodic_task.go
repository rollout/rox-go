package utils

import "time"

func RunPeriodicTask(action func(), period time.Duration, quit <-chan struct{}) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	select {
		case <-ticker.C:
			action()
		case <-quit:
			return
	}
}
