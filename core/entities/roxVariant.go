package entities

type roxVariant struct {
	name string
	flagType string
}

func (v *roxVariant) Name() string {
	return v.name
}

func (v *roxVariant) FlagType() string {
	return v.flagType
}
