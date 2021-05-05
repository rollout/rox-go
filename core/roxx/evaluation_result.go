package roxx

import "fmt"

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

func (ev EvaluationResult) IntValue() (int, error) {
	if value, ok := ev.value.(int); ok {
		return value, nil
	} else {
		return 0, fmt.Errorf("evaluation result is not an int")
	}
}

func (ev EvaluationResult) DoubleValue() (float64, error) {
	if value, ok := ev.value.(float64); ok {
		return value, nil
	} else if value, ok := ev.value.(int); ok {
		return float64(value), nil
	} else {
		return 0, fmt.Errorf("evaluation result is not a number")
	}
}
