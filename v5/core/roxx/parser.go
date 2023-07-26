package roxx

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rollout/rox-go/v5/core/context"
	"github.com/rollout/rox-go/v5/core/logging"
)

type Parser interface {
	EvaluateExpression(expression string, context context.Context) EvaluationResult
	AddOperator(name string, operation Operation)
}

type Operation = func(p Parser, stack *CoreStack, context context.Context)

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
			logging.GetLogger().Warn(fmt.Sprintf("Roxx Exception: Failed evaluate expression %s\n", r), nil)
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
			return NewEvaluationResult(nil)
		}
	}

	result = stack.Pop()
	if result == TokenTypeUndefined {
		result = nil
	}
	return NewEvaluationResult(result)
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
		op1, ok := stack.Pop().(string)
		if ok {
			hasher := md5.New()
			hasher.Write([]byte(op1))
			stack.Push(hex.EncodeToString(hasher.Sum(nil)))
		} else {
			stack.Push(TokenTypeUndefined)
		}
	})

	p.AddOperator("concat", func(p Parser, stack *CoreStack, context context.Context) {
		op1, ok1 := stack.Pop().(string)
		op2, ok2 := stack.Pop().(string)
		if ok1 && ok2 {
			stack.Push(fmt.Sprintf("%s%s", op1, op2))
		} else {
			stack.Push(TokenTypeUndefined)
		}
	})

	p.AddOperator("b64d", func(p Parser, stack *CoreStack, context context.Context) {
		op1, ok1 := stack.Pop().(string)
		if ok1 {
			sDec, _ := base64.StdEncoding.DecodeString(op1)
			stack.Push(string(sDec))
		} else {
			stack.Push(TokenTypeUndefined)
		}
	})
	
	p.AddOperator("tsToNum", func(p Parser, stack *CoreStack, context context.Context) {
		op1, ok1 := stack.Pop().(time.Time)
		if ok1 {
			// for better precision using milli
			stack.Push(float64(op1.UnixMilli()) / 1000)
		} else {
			stack.Push(TokenTypeUndefined)
		}
	})
}

func (p *roxxParser) reverseTokens(tokens []*Node) {
	for i, j := 0, len(tokens)-1; i < j; i, j = i+1, j-1 {
		tokens[i], tokens[j] = tokens[j], tokens[i]
	}
}
