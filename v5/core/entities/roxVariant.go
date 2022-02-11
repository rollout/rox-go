package entities

type roxVariant struct {
	name     string
	flagType int
}

func (v *roxVariant) Name() string {
	return v.name
}

func (v *roxVariant) FlagType() int {
	return v.flagType
}
