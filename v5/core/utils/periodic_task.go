package utils

import (
	"time"
)

func RunPeriodicTask(action func(), period time.Duration, quit <-chan struct{}) {
	go func() {
		// Start a new ticker to execute the task
		ticker := time.NewTicker(period)

		// Loop forever, executing the task on each tick. Only stop when the "quit"
		// channel tells us to, or the channel is closed.
		for {
			select {
			case <-ticker.C:
				action()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
