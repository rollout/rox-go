package utils

import (
	"testing"
	"time"
)

// TestRunPeriodicTask tests that RunPeriodicTask will run...periodically.
func TestRunPeriodicTask(t *testing.T) {

	// If this test is not done in 10 seconds, the test will fail
	passed := false
	go func() {
		timeLimit := time.Now().Add(10 * time.Second)
		for time.Now().Before(timeLimit) {
			if !passed {
				// Not passed yet. Sleep and check again.
				time.Sleep(1 * time.Second)
			} else {
				// Passed. Return.
				return
			}
		}
		// If we made it to here, fail the test.
		t.Fail()
		return
	}()

	done := make(chan struct{})

	// Start off a task that will send an increasing counter back on a channel
	task1Chan := make(chan int)
	taks1Counter := 1
	RunPeriodicTask(func() {
		taks1Counter += 1
		task1Chan <- taks1Counter
	}, time.Second, done)

	// Start off a second task that will send an increasing counter back on a channel
	task2Chan := make(chan int)
	task2Counter := 10
	RunPeriodicTask(func() {
		task2Counter += 1
		task2Chan <- task2Counter
	}, time.Second, done)

	// Wait for the first task to do at least 2 ticks.
	task1Done := make(chan bool)
	go func() {
		for {
			// Read the next value from the channel
			value := <-task1Chan
			// Check if it's at least 2 before signalling that the task is done
			if value >= 2 {
				task1Done <- true
				return
			}
		}
	}()

	// Wait for the second task to do at least 2 ticks.
	task2Done := make(chan bool)
	go func() {
		for {
			// Read the next value from the channel
			value := <-task2Chan
			// Check if it's at least 12 before signalling that the task is done
			if value >= 12 {
				task2Done <- true
				return
			}
		}
	}()

	// Wait for both tasks to say they're done
	<-task1Done
	<-task2Done

	// Close the done channel to stop the tasks
	close(done)

	// Call off the fail check goroutine
	passed = true
}
