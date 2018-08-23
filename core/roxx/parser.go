package roxx

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/rollout/rox-go/core/context"
	"time"
)

type Operation = func(p Parser, stack *CoreStack, context context.Context)

type Parser interface {
	EvaluateExpression(expression string, context context.Context) EvaluationResult
	AddOperator(name string, operation Operation)
}

type roxxParser struct {
	operatorsMap map[string]Operation
}

func NewParser() Parser {
	p := &roxxParser{
		operatorsMap: make(map[string]Operation),
	}
	p.setBasicOperators()
	NewValueCompareExtensions(p).Extend()
	NewRegularExpressionExtensions(p).Extend()
	return p
}

func (p *roxxParser) AddOperator(name string, operation Operation) {
	p.operatorsMap[name] = operation
}

func (p *roxxParser) EvaluateExpression(expression string, context context.Context) EvaluationResult {
	operators := make([]string, 0, len(p.operatorsMap))
	for operator := range p.operatorsMap {
		operators = append(operators, operator)
	}

	tokens := NewTokenizedExpression(expression, operators).GetTokens()
	p.reverseTokens(tokens)

	defer func() {
		if r := recover(); r != nil {
			// TODO logger
			fmt.Printf("Roxx Exception: Failed evaluate expression %s\n", r)
		}
	}()

	stack := NewCoreStack()
	var result interface{}

	for _, token := range tokens {
		if token.Type == NodeTypeRand {
			stack.Push(token.Value)
		} else if token.Type == NodeTypeRator {
			if handler, ok := p.operatorsMap[token.Value.(string)]; ok {
				handler(p, stack, context)
			}
		} else {
			return EvaluationResult{nil}
		}
	}

	result = stack.Pop()
	if result == TokenTypeUndefined {
		result = nil
	}
	return EvaluationResult{result}
}

func (p *roxxParser) setBasicOperators() {
	p.AddOperator("isUndefined", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		if tokenType, ok := op1.(*TokenType); !ok {
			stack.Push(false)
		} else {
			stack.Push(tokenType == TokenTypeUndefined)
		}
	})

	p.AddOperator("now", func(p Parser, stack *CoreStack, context context.Context) {
		stack.Push(int(time.Now().UnixNano() / 1e6))
	})

	p.AddOperator("and", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		op2 := stack.Pop()

		if op1 == TokenTypeUndefined {
			op1 = false
		}

		if op2 == TokenTypeUndefined {
			op2 = false
		}

		stack.Push(op1.(bool) && op2.(bool))
	})

	p.AddOperator("or", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		op2 := stack.Pop()

		if op1 == TokenTypeUndefined {
			op1 = false
		}

		if op2 == TokenTypeUndefined {
			op2 = false
		}

		stack.Push(op1.(bool) || op2.(bool))
	})

	p.AddOperator("ne", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		op2 := stack.Pop()

		if op1 == TokenTypeUndefined {
			op1 = false
		}

		if op2 == TokenTypeUndefined {
			op2 = false
		}

		stack.Push(op1 != op2)
	})

	p.AddOperator("eq", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		op2 := stack.Pop()

		if op1 == TokenTypeUndefined {
			op1 = false
		}

		if op2 == TokenTypeUndefined {
			op2 = false
		}

		stack.Push(op1 == op2)
	})

	p.AddOperator("not", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()

		if op1 == TokenTypeUndefined {
			op1 = false
		}

		stack.Push(!op1.(bool))
	})

	p.AddOperator("ifThen", func(p Parser, stack *CoreStack, context context.Context) {
		conditionExpression := stack.Pop().(bool)
		trueExpression := stack.Pop()
		falseExpression := stack.Pop()

		if conditionExpression {
			stack.Push(trueExpression)
		} else {
			stack.Push(falseExpression)
		}
	})

	p.AddOperator("inArray", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		op2 := stack.Pop()

		if op2, ok := op2.([]interface{}); !ok {
			stack.Push(false)
		} else {
			for _, item := range op2 {
				if item == op1 {
					stack.Push(true)
					return
				}
			}
			stack.Push(false)
		}
	})

	p.AddOperator("md5", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop().(string)
		hasher := md5.New()
		hasher.Write([]byte(op1))
		stack.Push(hex.EncodeToString(hasher.Sum(nil)))
	})

	p.AddOperator("concat", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop().(string)
		op2 := stack.Pop().(string)
		stack.Push(fmt.Sprintf("%s%s", op1, op2))
	})
}

func (p *roxxParser) reverseTokens(tokens []*Node) {
	for i, j := 0, len(tokens)-1; i < j; i, j = i+1, j-1 {
		tokens[i], tokens[j] = tokens[j], tokens[i]
	}
}
