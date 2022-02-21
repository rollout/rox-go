package register

import (
	"fmt"
	"github.com/rollout/rox-go/v4/core/model"
	"reflect"
	"sync"
)

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
		if ns != "" {
			name = fmt.Sprintf("%s.%s", ns, name)
		}

		r.flagRepository.AddFlag(variant, name)
	}
}
