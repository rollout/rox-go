package utils

import "time"

func RunPeriodicTask(action func(), period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			action()
		}
	}
}
