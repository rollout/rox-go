package repositories_test

import (
	"github.com/rollout/rox-go/core/properties"
	"github.com/rollout/rox-go/core/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomPropertyRepositoryWillReturnNullWhenPropNotFound(t *testing.T) {
	repo := repositories.NewCustomPropertyRepository()

	assert.Nil(t, repo.GetCustomProperty("harti"))
}

func TestCustomPropertyRepositoryWillAddProp(t *testing.T) {
	repo := repositories.NewCustomPropertyRepository()
	cp := properties.NewStringProperty("prop1", "123")
	repo.AddCustomProperty(cp)

	assert.Equal(t, "prop1", repo.GetCustomProperty("prop1").Name)
}

func TestCustomPropertyRepositoryWillNotOverrideProp(t *testing.T) {
	repo := repositories.NewCustomPropertyRepository()
	cp := properties.NewStringProperty("prop1", "123")
	cp2 := properties.NewStringProperty("prop1", "234")

	repo.AddCustomPropertyIfNotExists(cp)
	repo.AddCustomPropertyIfNotExists(cp2)

	assert.Equal(t, "123", repo.GetCustomProperty("prop1").Value(nil))
}

func TestCustomPropertyRepositoryWillOverrideProp(t *testing.T) {
	repo := repositories.NewCustomPropertyRepository()
	cp := properties.NewStringProperty("prop1", "123")
	cp2 := properties.NewStringProperty("prop1", "234")

	repo.AddCustomPropertyIfNotExists(cp)
	repo.AddCustomProperty(cp2)

	assert.Equal(t, "234", repo.GetCustomProperty("prop1").Value(nil))
}

func TestCustomPropertyRepositoryWillRaisePropAddedEvent(t *testing.T) {
	repo := repositories.NewCustomPropertyRepository()
	cp := properties.NewStringProperty("prop1", "123")

	var propFromEvent *properties.CustomProperty
	repo.RegisterPropertyAddedHandler(func(property *properties.CustomProperty) {
		propFromEvent = property
	})

	repo.AddCustomProperty(cp)

	assert.Equal(t, "prop1", propFromEvent.Name)
}

func TestCustomPropertyRepositoryRecoversFromPanic(t *testing.T) {
	repo := repositories.NewCustomPropertyRepository()
	cp := properties.NewStringProperty("prop1", "123")

	var propFromEvent *properties.CustomProperty
	repo.RegisterPropertyAddedHandler(func(property *properties.CustomProperty) {
		propFromEvent = property
		panic("mwahahahahaEvilUser")
	})

	repo.AddCustomProperty(cp)

	assert.Equal(t, "prop1", propFromEvent.Name)
}
