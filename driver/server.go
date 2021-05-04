package main

import (
	"context"
	"encoding/json"
	roxContext "github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/server"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type testRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

type setupAndAwait struct {
	Key     string        `json:"key"`
	Options serverOptions `json:"options, omitempty"`
}

type serverOptions struct {
	Configuration *json.RawMessage `json:"configuration, omitempty"`
}

type dynamicFlagIsEnabled struct {
	Flag         string           `json:"flag"`
	DefaultValue bool             `json:"defaultValue"`
	Context      *json.RawMessage `json:"context, omitempty"`
}

type dynamicFlagValue struct {
	Flag         string           `json:"flag"`
	DefaultValue string           `json:"defaultValue"`
	Context      *json.RawMessage `json:"context, omitempty"`
}

type setCustomString struct {
	Key   string
	Value string
}

type staticFlagIsEnabled struct {
	Flag    string           `json:"flag"`
	Context *json.RawMessage `json"context, omitempty"`
}

var srv *http.Server

func main() {

	var rox *server.Rox
	rox = server.NewRox()

	mux := http.NewServeMux()
	srv = &http.Server{
		Addr:    ":1234",
		Handler: mux,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		var payload json.RawMessage
		tr := testRequest{Payload: &payload}
		byteBody, err := ioutil.ReadAll(req.Body)
		if err := json.Unmarshal(byteBody, &tr); err != nil {
			log.Fatal(err)
		}

		switch tr.Action {
		case "shutdown":
			rox.Shutdown()
			return
		case "setupAndAwait":
			var setup setupAndAwait
			if err := json.Unmarshal(payload, &setup); err != nil {
				log.Fatal(err)
			}

			var configMap map[string]string

			json.Unmarshal(*setup.Options.Configuration, &configMap)

			if s, ok := configMap["env"]; ok {
				switch s {
				case "localhost":
					os.Setenv("ROLLOUT_MODE", "LOCAL")
				case "qa":
					os.Setenv("ROLLOUT_MODE", "QA")
				default:
					os.Setenv("ROLLOUT_MODE", "")

				}
			}
			options := server.NewRoxOptions(server.RoxOptionsBuilder{})
			err := <-rox.Setup(setup.Key, options)

			var result string

			if err != nil {
				result = err.Error()
			} else {
				result = "done"
			}

			doneStruct := struct {
				Result string `json:"result"`
			}{result}
			doneBody, err := json.Marshal(doneStruct)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(doneBody)
			return
		case "setCustomStringProperty":
			var setCustom setCustomString
			if err = json.Unmarshal(payload, &setCustom); err != nil {
				log.Fatal(err)
			}
			rox.SetCustomStringProperty(setCustom.Key, setCustom.Value)
			doneStruct := struct {
				Result string `json:"result"`
			}{"done"}
			doneBody, err := json.Marshal(doneStruct)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(doneBody)
			return
		case "dynamicFlagIsEnabled":
			var dynamicFlag dynamicFlagIsEnabled
			if err := json.Unmarshal(payload, &dynamicFlag); err != nil {
				log.Fatal(err)
			}

			contextMap := make(map[string]interface{})
			json.Unmarshal(*dynamicFlag.Context, &contextMap)
			rCtx := roxContext.NewContext(contextMap)
			for key, value := range contextMap {

				switch value.(type) {
				case int:
					rox.SetCustomIntegerProperty(key, value.(int))
					continue
				case float64:
					rox.SetCustomFloatProperty(key, value.(float64))
					continue
				case string:
					rox.SetCustomStringProperty(key, value.(string))
					continue
				case bool:
					rox.SetCustomBooleanProperty(key, value.(bool))
					continue
				}
			}
			result := rox.DynamicAPI().IsEnabled(dynamicFlag.Flag, dynamicFlag.DefaultValue, rCtx)
			doneStruct := struct {
				Result bool `json:"result"`
			}{result}
			doneBody, err := json.Marshal(doneStruct)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(doneBody)
			return
		case "stop":
			doneStruct := struct {
				Result string `json:"result"`
			}{"done"}
			doneBody, err := json.Marshal(doneStruct)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(doneBody)
			cancel()
			return
		case "dynamicFlagValue":
			var dynamicFlag dynamicFlagValue
			if err := json.Unmarshal(payload, &dynamicFlag); err != nil {
				log.Fatal(err)
			}
			contextMap := make(map[string]interface{})
			json.Unmarshal(*dynamicFlag.Context, &contextMap)
			rCtx := roxContext.NewContext(contextMap)

			for key, value := range contextMap {
				switch value.(type) {
				case int:
					rox.SetCustomIntegerProperty(key, value.(int))
					continue
				case float64:
					rox.SetCustomFloatProperty(key, value.(float64))
					continue
				case string:
					rox.SetCustomStringProperty(key, value.(string))
					continue
				case bool:
					rox.SetCustomBooleanProperty(key, value.(bool))
					continue
				}
			}
			result := rox.DynamicAPI().StringValue(dynamicFlag.Flag, dynamicFlag.DefaultValue, []string{}, rCtx)
			doneStruct := struct {
				Result string `json:"result"`
			}{result}
			doneBody, err := json.Marshal(doneStruct)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(doneBody)
			return
		case "registerStaticContainers":
			rox.Register("namespace", container)
			doneStruct := struct {
				Result string `json:"result"`
			}{"done"}
			doneBody, err := json.Marshal(doneStruct)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(doneBody)
			return
		case "staticFlagIsEnabled":
			var staticFlag staticFlagIsEnabled
			if err := json.Unmarshal(payload, &staticFlag); err != nil {
				log.Fatal(err)
			}

			contextMap := make(map[string]interface{})
			json.Unmarshal(*staticFlag.Context, &contextMap)
			rCtx := roxContext.NewContext(contextMap)

			var result bool
			if strings.Contains(staticFlag.Flag, "boolDefaultFalse") {
				result = container.BoolDefaultFalse.IsEnabled(rCtx)
			} else {
				result = container.BoolDefaultTrue.IsEnabled(rCtx)
			}

			doneStruct := struct {
				Result bool `json:"result"`
			}{result}
			doneBody, err := json.Marshal(doneStruct)
			if err != nil {
				log.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(doneBody)

		default:
			return
		}
	})
	mux.HandleFunc("/status-check", statusCheck)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Println("server started")
	select {
	case <-ctx.Done():
		// Shutdown the server when the context is canceled
		srv.Shutdown(ctx)
	}

	log.Println("server stopped")

}

func statusCheck(w http.ResponseWriter, req *http.Request) {
	return
}
