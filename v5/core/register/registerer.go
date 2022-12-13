package register

import (
	"fmt"
	"github.com/rollout/rox-go/v5/core/model"
	"reflect"
	"sync"
)

//The name of the tag we use to set the actual flag name in flag structs
const flagStructTagName = "fflag"

type Registerer struct {
	flagRepository model.FlagRepository
	namespaces     map[string]bool
	mutex          sync.Mutex
}

func NewRegisterer(flagRepository model.FlagRepository) *Registerer {
	return &Registerer{
		flagRepository: flagRepository,
		namespaces:     make(map[string]bool),
	}
}

func (r *Registerer) RegisterInstance(container interface{}, ns string) {
	r.mutex.Lock()
	if r.namespaces[ns] {
		panic(fmt.Sprintf("A container with the given namesapce (%s) has already been registered", ns))
	} else {
		r.namespaces[ns] = true
	}
	r.mutex.Unlock()

	v := reflect.Indirect(reflect.ValueOf(container))
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}

		variant, ok := v.Field(i).Interface().(model.Variant)
		if !ok {
			continue
		}

		name := v.Type().Field(i).Name
		//check for our tag in struct definition
		tag := v.Type().Field(i).Tag.Get(flagStructTagName)
		if tag == "" {
			//always set the tag
			tag = name
		}

		if ns != "" {
			name = fmt.Sprintf("%s.%s", ns, name)
			tag = fmt.Sprintf("%s.%s", ns, tag)
		}

		r.flagRepository.AddFlag(variant, name, tag)
	}
}
