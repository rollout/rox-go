package roxx

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	dictStartDelimiter      = "{"
	dictEndDelimiter        = "}"
	arrayStartDelimiter     = "["
	arrayEndDelimiter       = "]"
	tokenDelimiters         = "{}[]():, \t\r\n\""
	prePostStringChar       = ""
	stringDelimiter         = `"`
	escapedQuote            = `\"`
	escapedQuotePlaceholder = `\RO_Q`
)

type TokenizedExpression struct {
	expression       string
	operators        []string
	resultList       []*Node
	arrayAccumulator []interface{}
	dictAccumulator  map[string]interface{}
	dictKey          string
}

func NewTokenizedExpression(expression string, operators []string) *TokenizedExpression {
	return &TokenizedExpression{
		expression: expression,
		operators:  operators,
	}
}

func (te *TokenizedExpression) GetTokens() []*Node {
	return te.tokenize(te.expression)
}

func (te *TokenizedExpression) pushNode(node *Node) {
	if te.dictAccumulator != nil && te.dictKey == "" {
		te.dictKey = node.Value.(string)
	} else if te.dictAccumulator != nil && te.dictKey != "" {
		te.dictAccumulator[te.dictKey] = node.Value
		te.dictKey = ""
	} else if te.arrayAccumulator != nil {
		te.arrayAccumulator = append(te.arrayAccumulator, node.Value)
	} else {
		te.resultList = append(te.resultList, node)
	}
}

func (te *TokenizedExpression) tokenize(expression string) []*Node {
	te.resultList = nil
	te.dictAccumulator = nil
	te.arrayAccumulator = nil
	te.dictKey = ""

	delimitersToUse := tokenDelimiters
	normalizedExpression := strings.Replace(expression, escapedQuote, escapedQuotePlaceholder, -1)
	tokenizer := NewStringTokenizer(normalizedExpression, delimitersToUse, true)

	var prevToken, token string
	for tokenizer.hasMoreTokens() {
		prevToken = token
		token = tokenizer.nextToken(delimitersToUse)
		inString := delimitersToUse == stringDelimiter

		if !inString && token == dictStartDelimiter {
			te.dictAccumulator = make(map[string]interface{})
		} else if !inString && token == dictEndDelimiter {
			dictResult := te.dictAccumulator
			te.dictAccumulator = nil
			te.pushNode(te.nodeFromDict(dictResult))
		} else if !inString && token == arrayStartDelimiter {
			te.arrayAccumulator = make([]interface{}, 0)
		} else if !inString && token == arrayEndDelimiter {
			arrayResult := te.arrayAccumulator
			te.arrayAccumulator = nil
			te.pushNode(te.nodeFromArray(arrayResult))
		} else if token == stringDelimiter {
			if prevToken == stringDelimiter {
				te.pushNode(te.nodeFromToken(`""`))
			}

			if inString {
				delimitersToUse = tokenDelimiters
			} else {
				delimitersToUse = stringDelimiter
			}
		} else {
			if delimitersToUse == stringDelimiter {
				te.pushNode(NewNode(NodeTypeRand, strings.Replace(token, escapedQuotePlaceholder, escapedQuote, -1)))
			} else if !strings.Contains(tokenDelimiters, token) && token != prePostStringChar {
				te.pushNode(te.nodeFromToken(token))
			}
		}
	}

	return te.resultList
}

func (te *TokenizedExpression) nodeFromArray(items []interface{}) *Node {
	return NewNode(NodeTypeRand, items)
}

func (te *TokenizedExpression) nodeFromDict(items map[string]interface{}) *Node {
	return NewNode(NodeTypeRand, items)
}

func (te *TokenizedExpression) nodeFromToken(token string) *Node {
	if te.isOperator(token) {
		return NewNode(NodeTypeRator, token)
	} else {
		if token == roxxTrue {
			return NewNode(NodeTypeRand, true)
		}
		if token == roxxFalse {
			return NewNode(NodeTypeRand, false)
		}
		if token == roxxUndefined {
			return NewNode(NodeTypeRand, TokenTypeUndefined)
		}

		tokenType := TokenTypeFromToken(token)
		switch tokenType {
		case TokenTypeString:
			return NewNode(NodeTypeRand, token[1:len(token)-1])
		case TokenTypeNumber:
			intNumber, err := strconv.Atoi(token)
			if err == nil {
				return NewNode(NodeTypeRand, intNumber)
			}

			number, err := strconv.ParseFloat(token, 64)
			if err == nil {
				return NewNode(NodeTypeRand, number)
			}

			panic(fmt.Sprintf("Excepted Number, got '%s' (%s)", token, tokenType.text))
		}
	}

	return NewNode(NodeTypeUnknown, nil)
}

func (te *TokenizedExpression) isOperator(token string) bool {
	for _, operator := range te.operators {
		if operator == token {
			return true
		}
	}
	return false
}
