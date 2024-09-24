package roxx_test

import (
	"testing"

	"github.com/rollout/rox-go/v6/core/roxx"
	"github.com/stretchr/testify/assert"
)

func TestCoreStackWillPushIntoStackString(t *testing.T) {
	testString := "stringTest"
	stack := roxx.NewCoreStack()
	stack.Push(testString)
	poppedItem := stack.Pop()

	assert.Equal(t, testString, poppedItem)
}

func TestCoreStackWillPushIntoStackInteger(t *testing.T) {
	testInt := 5
	stack := roxx.NewCoreStack()
	stack.Push(testInt)
	poppedItem := stack.Pop()

	assert.Equal(t, testInt, poppedItem)
}

func TestCoreStackWillPushIntoStackIntegerAndString(t *testing.T) {
	testInt := 5
	testString := "stringTest"
	stack := roxx.NewCoreStack()
	stack.Push(testInt)
	stack.Push(testString)
	poppedItemFirst := stack.Pop()
	poppedItemSecond := stack.Pop()

	assert.Equal(t, testString, poppedItemFirst)
	assert.Equal(t, testInt, poppedItemSecond)
}

func TestCoreStackWillPeekFromStack(t *testing.T) {
	testInt := 5
	testString := "stringTest"
	stack := roxx.NewCoreStack()
	stack.Push(testInt)
	stack.Push(testString)
	peekedFirst := stack.Peek()
	poppedItem := stack.Pop()

	assert.Equal(t, peekedFirst, poppedItem)
	assert.Equal(t, testString, poppedItem)

}
