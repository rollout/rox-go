package reporting

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/rollout/rox-go/core/logging"
	"github.com/rollout/rox-go/core/model"
)

const (
	SentryDsn = "https://1d1f883fe6af42b0bb276f1156760949@sentry.io/1546802"
)

type errorReporter struct {
	deviceProperties model.DeviceProperties
	buid             model.BUID
	sentryHub        *sentry.Hub
}

func disableIntegrations([]sentry.Integration) []sentry.Integration {
	return []sentry.Integration{}
}

func NewErrorReporter(deviceProperties model.DeviceProperties, buid model.BUID) model.ErrorReporter {
	// Do not use the Sentry.init() method but create instead our own client, scope and hub.
	// This way we are allowing our users to use their own Sentry in their own app.
	sentryScope := sentry.NewScope()
	sentryClient, _ := sentry.NewClient(sentry.ClientOptions{
		Dsn:            SentryDsn,
		MaxBreadcrumbs: 0,
		Integrations:   disableIntegrations,
		Environment:    deviceProperties.RolloutEnvironment(),
		Release:        "Rollout Go SDK " + deviceProperties.LibVersion(),
	})
	sentryHub := sentry.NewHub(sentryClient, sentryScope)

	return &errorReporter{
		deviceProperties: deviceProperties,
		buid:             buid,
		sentryHub:        sentryHub,
	}
}

func (er *errorReporter) Report(message string, err error) {
	if er.deviceProperties.RolloutEnvironment() == "LOCAL" {
		return
	}

	logging.GetLogger().Error(fmt.Sprintf("Error report: %s", message), err)

	stackTrace := sentry.NewStacktrace()
	event := er.createEvent(message, err, stackTrace)
	er.sendError(event)
}

func (er *errorReporter) sendError(event *sentry.Event) {
	logging.GetLogger().Debug("Sending sentry error report...", nil)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logging.GetLogger().Error("Failed to send sentry error", r)
			}
		}()

		eventID := er.sentryHub.CaptureEvent(event)
		if eventID == nil {
			logging.GetLogger().Error("Failed to send sentry error report", nil)
		} else {
			logging.GetLogger().Debug("Sentry error report was sent", nil)
		}
	}()
}

func (er *errorReporter) createEvent(message string, err error, stackTrace *sentry.Stacktrace) *sentry.Event {
	user := sentry.User{
		ID: er.deviceProperties.RolloutKey(),
	}
	exception := sentry.Exception{
		Stacktrace: stackTrace,
		Value:      err.Error(),
	}
	extra := map[string]interface{}{
		"deviceId": er.deviceProperties.DistinctID(),
		"buid":     er.buid.String(),
	}
	event := sentry.Event{
		Message:   message,
		User:      user,
		Exception: []sentry.Exception{exception},
		Extra:     extra,
	}
	return &event
}
