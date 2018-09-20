package roxx

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	TokenTypeNotAType  = newTokenType("NOT_A_TYPE", "")
	TokenTypeString    = newTokenType("StringType", `"((\\.)|[^\\\\"])*"`)
	TokenTypeNumber    = newTokenType("NumberType", `[\-]{0,1}\d+[\.]\d+|[\-]{0,1}\d+`)
	TokenTypeBoolean   = newTokenType("BooleanType", fmt.Sprintf("%s|%s", roxxTrue, roxxFalse))
	TokenTypeUndefined = newTokenType("UndefinedType", roxxUndefined)
)

func TokenTypeFromToken(token string) *TokenType {
	if token != "" {
		testedToken := strings.ToLower(token)
		for _, tokenType := range []*TokenType{TokenTypeString, TokenTypeNumber, TokenTypeBoolean, TokenTypeUndefined} {
			if tokenType.pattern.MatchString(testedToken) {
				return tokenType
			}
		}
	}
	return TokenTypeNotAType
}

type TokenType struct {
	text    string
	pattern *regexp.Regexp
}

func newTokenType(text, pattern string) *TokenType {
	return &TokenType{
		text:    text,
		pattern: regexp.MustCompile(pattern),
	}
}

func (tt *TokenType) IsNumber() bool {
	return tt == TokenTypeNumber
}

func (tt *TokenType) IsString() bool {
	return tt == TokenTypeString
}

func (tt *TokenType) IsBoolean() bool {
	return tt == TokenTypeBoolean
}

func (tt *TokenType) IsUndefined() bool {
	return tt == TokenTypeUndefined
}
