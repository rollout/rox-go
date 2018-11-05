package notifications

import (
	"fmt"
	"github.com/rollout/rox-go/core/logging"
	"github.com/rollout/sse"
	"io"
	"strings"
	"sync"
	"time"
)

const (
	connectionRetryInterval = time.Second * 3
)

type Event struct {
	EventName string
	Data      string
}

type EventHandler = func(event Event)

type NotificationListener struct {
	listenURL string
	appKey    string

	handlers      map[string][]EventHandler
	handlersMutex sync.RWMutex
	stop          chan struct{}
}

func NewNotificationListener(listenURL, appKey string) *NotificationListener {
	return &NotificationListener{
		listenURL: listenURL,
		appKey:    appKey,
		handlers:  make(map[string][]EventHandler),
	}
}

func (nl *NotificationListener) Start() {
	sseURL := fmt.Sprintf("%s/%s", strings.TrimSuffix(nl.listenURL, "/"), nl.appKey)
	nl.stop = make(chan struct{})
	go nl.run(sseURL)
}

func (nl *NotificationListener) Stop() {
	close(nl.stop)
}

func (nl *NotificationListener) On(eventName string, handler EventHandler) {
	nl.handlersMutex.Lock()
	nl.handlers[eventName] = append(nl.handlers[eventName], handler)
	nl.handlersMutex.Unlock()
}

func (nl *NotificationListener) run(sseURL string) {
	for {
		sseClient := sse.NewClientWithoutRetry(sseURL)
		events := make(chan *sse.Event)
		sseCloser, err := sseClient.SubscribeChan("", events)
		if err != nil {
			logging.GetLogger().Warn("Can't subscribe to SSE events", err)
			time.Sleep(connectionRetryInterval)

			select {
			case <-nl.stop:
				return
			default:
			}
		} else {
			nl.readEvents(events, sseCloser)
		}
	}
}

func (nl *NotificationListener) readEvents(events <-chan *sse.Event, sseCloser io.Closer) {
	for {
		select {
		case <-nl.stop:
			err := sseCloser.Close()
			if err != nil {
				logging.GetLogger().Warn("Can't close SSE closer", err)
			}
			return
		case event, ok := <-events:
			if !ok {
				return
			}
			nl.invokeHandlers(event)
		}
	}
}

func (nl *NotificationListener) invokeHandlers(rawEvent *sse.Event) {
	event := Event{EventName: string(rawEvent.Event), Data: string(rawEvent.Data)}

	nl.handlersMutex.RLock()
	handlers := make([]EventHandler, len(nl.handlers[event.EventName]))
	copy(handlers, nl.handlers[event.EventName])
	nl.handlersMutex.RUnlock()

	for _, handler := range handlers {
		nl.invokeHandler(handler, event)
	}
}

func (nl *NotificationListener) invokeHandler(handler EventHandler, event Event) {
	defer func() {
		if r := recover(); r != nil {
			logging.GetLogger().Error(fmt.Sprintf("SSE handler panics: %s", r), nil)
		}
	}()

	handler(event)
}
