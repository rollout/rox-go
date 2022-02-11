package roxx

type EvaluationResult struct {
	value interface{}
}

func NewEvaluationResult(value interface{}) EvaluationResult {
	return EvaluationResult{value: value}
}

func (ev EvaluationResult) Value() interface{} {
	return ev.value
}

func (ev EvaluationResult) BoolValue() bool {
	if ev.value == nil {
		return false
	} else if value, ok := ev.value.(bool); ok {
		return value
	}

	return false
}

func (ev EvaluationResult) StringValue() string {
	if value, ok := ev.value.(string); ok {
		return value
	} else if value, ok := ev.value.(bool); ok {
		if value {
			return FlagTrueValue
		}
		return FlagFalseValue
	}

	return ""
}
