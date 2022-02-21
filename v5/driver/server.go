package main

import (
	"context"
	"encoding/json"
	roxContext "github.com/rollout/rox-go/v5/core/context"
	"github.com/rollout/rox-go/v5/server"
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
	Env           string           `json:"env"`
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

type setCustomPropertyToThrow struct {
	Key string `json:"key"`
}

type staticFlagIsEnabled struct {
	Flag    string           `json:"flag"`
	Context *json.RawMessage `json"context, omitempty"`
}

var srv *http.Server

func main() {

	var rox *server.Rox
	rox = server.NewRox()

	var port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "1234"
	}

	mux := http.NewServeMux()
	srv = &http.Server{
		Addr:    ":" + port,
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
		case "setupAndAwait":
			var setup setupAndAwait
			if err := json.Unmarshal(payload, &setup); err != nil {
				log.Fatal(err)
			}

			var configMap map[string]string

			json.Unmarshal(*setup.Options.Configuration, &configMap)

			switch setup.Options.Env {
			case "localhost":
				os.Setenv("ROLLOUT_MODE", "LOCAL")
			case "qa":
				os.Setenv("ROLLOUT_MODE", "QA")
			default:
				os.Setenv("ROLLOUT_MODE", "")
			}
			options := server.NewRoxOptions(server.RoxOptionsBuilder{})
			<-rox.Setup(setup.Key, options)

			sendDone(w)
			return
		case "setCustomStringProperty":
			var setCustom setCustomString
			if err = json.Unmarshal(payload, &setCustom); err != nil {
				log.Fatal(err)
			}
			rox.SetCustomStringProperty(setCustom.Key, setCustom.Value)
			sendDone(w)
			return
		case "dynamicFlagIsEnabled":
			var dynamicFlag dynamicFlagIsEnabled
			if err := json.Unmarshal(payload, &dynamicFlag); err != nil {
				log.Fatal(err)
			}

			contextMap := make(map[string]interface{})
			json.Unmarshal(*dynamicFlag.Context, &contextMap)
			rCtx := roxContext.NewContext(contextMap)
			result := rox.DynamicAPI().IsEnabled(dynamicFlag.Flag, dynamicFlag.DefaultValue, rCtx)
			sendResult(w, struct {
				Result bool `json:"result"`
			}{result})
			return
		case "stop":
			sendResult(w, "done")
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
			result := rox.DynamicAPI().Value(dynamicFlag.Flag, dynamicFlag.DefaultValue, []string{}, rCtx)
			sendResult(w, struct {
				Result string `json:"result"`
			}{result})
			return
		case "registerStaticContainers":
			rox.Register("namespace", container)
			log.Println("Registered static container in namespace \"namespace\"")
			sendDone(w)
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
			if strings.Contains(staticFlag.Flag, "BoolDefaultFalse") {
				result = container.BoolDefaultFalse.IsEnabled(rCtx)
			} else {
				result = container.BoolDefaultTrue.IsEnabled(rCtx)
			}

			sendResult(w, struct {
				Result bool `json:"result"`
			}{result})
			return

		case "setCustomPropertyToThrow":
			var prop setCustomPropertyToThrow
			if err := json.Unmarshal(payload, &prop); err != nil {
				log.Fatal(err)
			}
			rox.SetCustomComputedStringProperty(prop.Key, func(context roxContext.Context) string {
				panic("")
				return ""
			})
			sendDone(w)
			return

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

func sendResult(w http.ResponseWriter, doneStruct interface{}) {
	doneBody, err := json.Marshal(doneStruct)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(doneBody)
}

func sendDone(w http.ResponseWriter) {
	doneStruct := struct {
		Result string `json:"result"`
	}{"done"}
	sendResult(w, doneStruct)
}

func statusCheck(w http.ResponseWriter, req *http.Request) {
	return
}
