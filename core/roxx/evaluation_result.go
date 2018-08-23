package roxx

type EvaluationResult struct {
	value interface{}
}

func (ev EvaluationResult) Value() interface{} {
	return ev.value
}

func (ev EvaluationResult) BoolValue() *bool {
	var result *bool

	if ev.value == nil {
		result = new(bool)
		*result = false
	} else if value, ok := ev.value.(bool); ok {
		result = new(bool)
		*result = value
	}

	return result
}

func (ev EvaluationResult) StringValue() string {
	if value, ok := ev.value.(string); ok {
		return value
	} else if value, ok := ev.value.(bool); ok {
		if value {
			return FlagTrueValue
		} else {
			return FlagFalseValue
		}
	}

	return ""
}
