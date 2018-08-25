package reporting

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/rollout/rox-go/core/logging"
	"github.com/rollout/rox-go/core/model"
	"runtime"
)

const (
	BugsnagNotifyUrl = "https://notify.bugsnag.com"
)

type errorReporter struct {
	request          model.Request
	deviceProperties model.DeviceProperties
	buid             model.BUID
}

func NewErrorReporter(request model.Request, deviceProperties model.DeviceProperties, buid model.BUID) model.ErrorReporter {
	return &errorReporter{
		request:          request,
		deviceProperties: deviceProperties,
		buid:             buid,
	}
}

func (er *errorReporter) Report(message string, err error) {
	if er.deviceProperties.RolloutEnvironment() == "LOCAL" {
		return
	}

	logging.GetLogger().Error(fmt.Sprintf("Error report: %s", message), err)

	stackTrace := er.getStackTrace()

	payload := er.createPayload(message, err, stackTrace)
	er.sendError(payload)
}

func (er *errorReporter) sendError(payload interface{}) {
	logging.GetLogger().Debug("Sending bugsnag error report...", nil)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logging.GetLogger().Error("Failed to send bugsnag error", r)
			}
		}()

		_, err := er.request.SendPost(BugsnagNotifyUrl, payload)
		if err != nil {
			logging.GetLogger().Error("Failed to send bugsnag error", err)
		} else {
			logging.GetLogger().Debug("Bugsnag error report was sent", nil)
		}
	}()
}

func (er *errorReporter) createPayload(message string, err error, stackTrace []*errors.StackFrame) map[string]interface{} {
	payload := make(map[string]interface{})
	er.addApiKey(payload)
	er.addNotifier(payload)
	er.addEvents(message, err, stackTrace, payload)
	return payload
}

func (er *errorReporter) addMetadata(message string, ev map[string]interface{}) {
	innerData := map[string]string{
		"message":  message,
		"deviceId": er.deviceProperties.DistinctID(),
		"buid":     er.buid.String(),
	}
	metaData := map[string]interface{}{
		"data": innerData,
	}
	ev["metaData"] = metaData
}

func (er *errorReporter) addApiKey(payload map[string]interface{}) {
	payload["apiKey"] = "9569ec14f61546c6aa2a97856492bf4d"
}

func (er *errorReporter) addEvents(message string, err error, stackTrace []*errors.StackFrame, payload map[string]interface{}) {
	evs := make([]map[string]interface{}, 0)
	er.addEvent(message, err, stackTrace, &evs)
	payload["events"] = evs
}

func (er *errorReporter) addEvent(message string, err error, stackTrace []*errors.StackFrame, events *[]map[string]interface{}) {
	ev := er.addPayloadVersion()
	er.addExceptions(message, err, stackTrace, ev)
	er.addUser("id", er.deviceProperties.RolloutKey(), ev)
	er.addMetadata(message, ev)
	er.addApp(ev)

	*events = append(*events, ev)
}

func (er *errorReporter) addPayloadVersion() map[string]interface{} {
	return map[string]interface{}{
		"payloadVersion": "2",
	}
}

func (er *errorReporter) addNotifier(payload map[string]interface{}) {
	notifier := map[string]string{
		"name":    "Rollout Go SDK",
		"version": er.deviceProperties.LibVersion(),
	}
	payload["notifier"] = notifier
}

func (er *errorReporter) addUser(id, rolloutKey string, ev map[string]interface{}) {
	user := map[string]string{
		"id":         id,
		"rolloutKey": rolloutKey,
	}
	ev["user"] = user
}

func (er *errorReporter) addExceptions(message string, err error, stackTrace []*errors.StackFrame, ev map[string]interface{}) {
	exception := make(map[string]interface{})

	if err == nil {
		exception["errorClass"] = message
		exception["message"] = message
		exception["stacktrace"] = []string{}
	} else {
		exception["errorClass"] = err.Error()
		exception["message"] = err.Error()

		var stacktrace []map[string]interface{}
		for _, frame := range stackTrace {
			fr := map[string]interface{}{
				"file":         frame.File,
				"method":       frame.Name,
				"lineNumber":   frame.LineNumber,
				"columnNumber": 0,
			}
			stacktrace = append(stacktrace, fr)
		}
		exception["stacktrace"] = stacktrace
	}

	exceptions := make([]map[string]interface{}, 1)
	exceptions[0] = exception
	ev["exceptions"] = exceptions
}

func (er *errorReporter) addApp(ev map[string]interface{}) {
	app := map[string]string{
		"releaseStage": er.deviceProperties.RolloutEnvironment(),
		"version":      er.deviceProperties.LibVersion(),
	}
	ev["app"] = app
}

func (er *errorReporter) getStackTrace() []*errors.StackFrame {
	maxStackDepth := 100
	stack := make([]uintptr, maxStackDepth)
	length := runtime.Callers(2, stack)
	frames := make([]*errors.StackFrame, length)
	for i, s := range stack[:length] {
		frame := errors.NewStackFrame(s)
		frames[i] = &frame
	}
	return frames
}
