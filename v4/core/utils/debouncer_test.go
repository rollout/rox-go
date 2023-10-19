package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWillTestDebouncerCalledAfterInterval(t *testing.T) {
	counter := 0
	debouncer := NewDebouncer(2000, func() {
		counter++
	})

	assert.Equal(t, 0, counter)
	debouncer.Invoke()
	assert.Equal(t, 0, counter)
	timer := time.NewTimer(1000 * time.Millisecond)
	<-timer.C
	assert.Equal(t, 0, counter)
	timer = time.NewTimer(1500 * time.Millisecond)
	<-timer.C
	assert.Equal(t, 1, counter)
}

func TestWillTestDebouncerSkipDoubleInvoke(t *testing.T) {
	counter := 0
	debouncer := NewDebouncer(1000, func() {
		counter++
	})

	assert.Equal(t, 0, counter)
	debouncer.Invoke()
	assert.Equal(t, 0, counter)
	timer := time.NewTimer(500 * time.Millisecond)
	<-timer.C
	assert.Equal(t, 0, counter)
	debouncer.Invoke()
	assert.Equal(t, 0, counter)
	timer = time.NewTimer(600 * time.Millisecond)
	<-timer.C
	assert.Equal(t, 1, counter)
	timer = time.NewTimer(600 * time.Millisecond)
	<-timer.C
	assert.Equal(t, 1, counter)
}

func TestWillTestDebouncerInvokeAfterInvoke(t *testing.T) {
	counter := 0
	debouncer := NewDebouncer(1000, func() {
		counter++
	})

	assert.Equal(t, 0, counter)
	debouncer.Invoke()
	assert.Equal(t, 0, counter)
	timer := time.NewTimer(1100 * time.Millisecond)
	<-timer.C
	assert.Equal(t, 1, counter)
	debouncer.Invoke()
	assert.Equal(t, 1, counter)
	timer = time.NewTimer(800 * time.Millisecond)
	<-timer.C
	assert.Equal(t, 1, counter)
	timer = time.NewTimer(300 * time.Millisecond)
	<-timer.C
	assert.Equal(t, 2, counter)
}
