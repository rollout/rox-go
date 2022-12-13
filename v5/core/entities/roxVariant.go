package entities

type roxVariant struct {
	name     string
	flagType int
	tag      string
}

func (v *roxVariant) Name() string {
	return v.name
}

func (v *roxVariant) Tag() string {
	return v.tag
}

func (v *roxVariant) FlagType() int {
	return v.flagType
}
