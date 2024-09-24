package roxx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rollout/rox-go/v6/core/context"
	"github.com/rollout/rox-go/v6/core/extensions"
	"github.com/rollout/rox-go/v6/core/properties"
	"github.com/rollout/rox-go/v6/core/repositories"
	"github.com/rollout/rox-go/v6/core/roxx"
	"github.com/stretchr/testify/assert"
)

func TestParserSimpleTokenization(t *testing.T) {
	operators := []string{"eq", "lt"}
	tokens := roxx.NewTokenizedExpression(`eq(false, lt(-123, "123"))`, operators).GetTokens()

	assert.Equal(t, 5, len(tokens))
	assert.Equal(t, roxx.NodeTypeRator, tokens[0].Type)
	assert.Equal(t, false, tokens[1].Value)
	assert.Equal(t, -123, tokens[3].Value)
	assert.Equal(t, "123", tokens[4].Value)
}

func TestParserTokenType(t *testing.T) {
	assert.True(t, roxx.TokenTypeFromToken("123").IsNumber())
	assert.True(t, roxx.TokenTypeFromToken("-123").IsNumber())
	assert.True(t, roxx.TokenTypeFromToken("-123.23").IsNumber())
	assert.True(t, roxx.TokenTypeFromToken("123.23").IsNumber())

	assert.False(t, roxx.TokenTypeFromToken("-123").IsString())
	assert.True(t, roxx.TokenTypeFromToken(`"-123"`).IsString())
	assert.True(t, roxx.TokenTypeFromToken(`"undefined"`).IsString())
	assert.False(t, roxx.TokenTypeFromToken("undefined").IsString())

	assert.True(t, roxx.TokenTypeFromToken("false").IsBoolean())
	assert.True(t, roxx.TokenTypeFromToken("true").IsBoolean())
	assert.False(t, roxx.TokenTypeFromToken("undefined").IsBoolean())

	assert.True(t, roxx.TokenTypeFromToken("undefined").IsUndefined())
	assert.False(t, roxx.TokenTypeFromToken("false").IsUndefined())
}

func TestParserSimpleExpressionEvaluation(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, "", parser.EvaluateExpression(`""`, nil).Value())
	assert.Equal(t, "", parser.EvaluateExpression(`\"\"`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`true`, nil).Value())
	assert.Equal(t, "red", parser.EvaluateExpression(`"red"`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`and(true, or(true, true))`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`and(true, or(false, true))`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`not(and(false, or(false, true)))`, nil).Value())
}

func TestParserEqExpressionsEvaluation(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, true, parser.EvaluateExpression(`eq("la la", "la la")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`eq("la la", "la,la")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`eq("lala", "lala")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`ne(100.123, 100.321)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`not(eq(undefined, undefined))`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`not(eq(not(undefined), undefined))`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`not(undefined)`, nil).Value())

	roxxString := `la \"la\" la`
	assert.Equal(t, true, parser.EvaluateExpression(fmt.Sprintf(`eq("%s", "la \"la\" la")`, roxxString), nil).Value())
}

func TestParserComparisonExpressionsEvaluation(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, false, parser.EvaluateExpression(`lt(500, 100)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`lt(500, 500)`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`lt(500, 500.54)`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`lte(500, 500)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`lt("500", "100")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`lt("500", "500")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`lt("500", "500.54")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`lte("500", "500")`, nil).Value())

	assert.Equal(t, true, parser.EvaluateExpression(`gt(500, 100)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`gt(500, 500)`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`gt(500.54, 500)`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`gte(500, 500)`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`gt("500", "100")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`gt("500", "500")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`gt("500.54", "500")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`gte("500", "500")`, nil).Value())

	assert.Equal(t, true, parser.EvaluateExpression(`gte("500", 500)`, nil).Value())

}

func TestParserNumEqualityExpressionsEvaluation(t *testing.T) {

	parser := roxx.NewParser()

	assert.Equal(t, true, parser.EvaluateExpression(`numeq(500, 500)`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`numeq("500", "500")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`numeq(500, "500")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`numeq("500", 500)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`numeq(500, 501)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`numeq("500", "501")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`numeq(500, "501")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`numeq("500", 501)`, nil).Value())

	assert.Equal(t, false, parser.EvaluateExpression(`numneq(500, 500)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`numneq("500", "500")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`numneq(500, "500")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`numneq("500", 500)`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`numneq(500, 501)`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`numneq("500", "501")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`numneq(500, "501")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`numneq("500", 501)`, nil).Value())

}

func TestParserSemVerComparisonEvaluation(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, false, parser.EvaluateExpression(`semverLt("1.1.0", "1.1")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`semverLte("1.1.0", "1.1")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`semverGte("1.1.0", "1.1")`, nil).Value())

	assert.Equal(t, false, parser.EvaluateExpression(`semverEq("1.0.0", "1")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`semverNe("1.0.1", "1.0.0.1")`, nil).Value())

	assert.Equal(t, true, parser.EvaluateExpression(`semverLt("1.1", "1.2")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`semverLte("1.1", "1.2")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`semverGt("1.1.1", "1.2")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`semverGt("1.2.1", "1.2")`, nil).Value())
}

func TestParserComparisonWithUndefinedEvaluation(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, false, parser.EvaluateExpression(`gte(500, undefined)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`gt(500, undefined)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`lte(500, undefined)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`lt(500, undefined)`, nil).Value())

	assert.Equal(t, false, parser.EvaluateExpression(`semverGte("1.1", undefined)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`semverGt("1.1", undefined)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`semverLte("1.1", undefined)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`semverLt("1.1", undefined)`, nil).Value())
}

func TestParserUnknownOperatorEvaluation(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, nil, parser.EvaluateExpression(`NOT_AN_OPERATOR(500, 500)`, nil).Value())
	assert.Equal(t, nil, parser.EvaluateExpression(`JUSTAWORD(500, 500)`, nil).Value())
}

func TestParserUndefinedEvaluation(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, true, parser.EvaluateExpression(`isUndefined(undefined)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`isUndefined(123123)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`isUndefined("undefined")`, nil).Value())
}

func TestParserNowEvaluation(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, true, parser.EvaluateExpression(`gte(now(), now())`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`gte(now(), 2458.123)`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`gte(now(), 1534759307565)`, nil).Value())
}

func TestParserRegularExpressionEvaluation(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, false, parser.EvaluateExpression(`match("111", "222", "")`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`match(".*", "222", "")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("22222", ".*", "")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("22222", "^2*$", "")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("22222", "^2*$", \"\")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("22222", "^2*$", b64d(\"\"))`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("test@shimi.com", ".*(com|ca)", "")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("test@jet.com", ".*jet\.com$", "")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("US", ".*IL|US", "")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("US", "IL|US"), ""`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("US", "(IL|US)", "")`, nil).Value())

	// Test flags
	assert.Equal(t, false, parser.EvaluateExpression(`match("Us", "(IL|US)", "")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("uS", "(IL|US)", "i")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`match("\n", ".", "s")`, nil).Value())

	// Unsupported flags (x)
	// assert.Equal(t, true, parser.EvaluateExpression(`match("uS", "IL|US#Comment", "xi")`, nil).Value())
	// assert.Equal(t, true, parser.EvaluateExpression(`match("HELLO\nTeST\n#This is a comment", "^TEST$", "ixm")`, nil).Value())
}

func TestParserIfThenExpressionEvaluationString(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, `AB`, parser.EvaluateExpression(`ifThen(and(true, or(true, true)), "AB", "CD")`, nil).Value())
	assert.Equal(t, `CD`, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), "AB", "CD")"`, nil).Value())

	assert.Equal(t, `AB`, parser.EvaluateExpression(`ifThen(and(true, or(true, true)), "AB", ifThen(and(true, or(true, true)), "EF", "CD"))`, nil).Value())
	assert.Equal(t, `EF`, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), "AB", ifThen(and(true, or(true, true)), "EF", "CD"))`, nil).Value())
	assert.Equal(t, `CD`, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), "AB", ifThen(and(true, or(false, false)), "EF", "CD"))`, nil).Value())

	assert.Equal(t, nil, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), "AB", ifThen(and(true, or(false, false)), "EF", undefined))`, nil).Value())
}

func TestParserIfThenExpressionEvaluationIntNumber(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, 1, parser.EvaluateExpression(`ifThen(and(true, or(true, true)), 1, 2)`, nil).Value())
	assert.Equal(t, 2, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), 1, 2)`, nil).Value())

	assert.Equal(t, 1, parser.EvaluateExpression(`ifThen(and(true, or(true, true)), 1, ifThen(and(true, or(true, true)), 3, 2))`, nil).Value())
	assert.Equal(t, 3, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), 1, ifThen(and(true, or(true, true)), 3, 2))`, nil).Value())
	assert.Equal(t, 2, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), 1, ifThen(and(true, or(false, false)), 3, 2))`, nil).Value())

	assert.Equal(t, nil, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), 1, ifThen(and(true, or(false, false)), 3, undefined))`, nil).Value())
}

func TestParserIfThenExpressionEvaluationFloatNumber(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, 1.1, parser.EvaluateExpression(`ifThen(and(true, or(true, true)), 1.1, 2.2)`, nil).Value())
	assert.Equal(t, 2.2, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), 1.1, 2.2)`, nil).Value())

	assert.Equal(t, 1.1, parser.EvaluateExpression(`ifThen(and(true, or(true, true)), 1.1, ifThen(and(true, or(true, true)), 3.3, 2.2))`, nil).Value())
	assert.Equal(t, 3.3, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), 1.1, ifThen(and(true, or(true, true)), 3.3, 2.2))`, nil).Value())
	assert.Equal(t, 2.2, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), 1.1, ifThen(and(true, or(false, false)), 3.3, 2.2))`, nil).Value())

	assert.Equal(t, nil, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), 1.1, ifThen(and(true, or(false, false)), 3.3, undefined))`, nil).Value())
}

func TestParserIfThenExpressionEvaluationBoolean(t *testing.T) {
	parser := roxx.NewParser()

	assert.Equal(t, true, parser.EvaluateExpression(`ifThen(and(true, or(true, true)), true, false)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), true, false)`, nil).Value())

	assert.Equal(t, false, parser.EvaluateExpression(`ifThen(and(true, or(true, true)), false, ifThen(and(true, or(true, true)), true, true))`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), false, ifThen(and(true, or(true, true)), true, false))`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), true, ifThen(and(true, or(false, false)), true, false))`, nil).Value())

	assert.Equal(t, true, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), false, ifThen(and(true, or(false, false)), false, (and(true,true))))`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), true, ifThen(and(true, or(false, false)), true, (and(true,false))))`, nil).Value())

	assert.Equal(t, nil, parser.EvaluateExpression(`ifThen(and(false, or(true, true)), true, ifThen(and(true, or(false, false)), true, undefined))`, nil).Value())
}

func TestParserInArray(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()

	parser.AddOperator(`mergeSeed`, func(p roxx.Parser, stack *roxx.CoreStack, context context.Context) {
		seed1 := stack.Pop()
		seed2 := stack.Pop()
		stack.Push(fmt.Sprintf("%s.%s", seed1, seed2))
	})

	assert.Equal(t, false, parser.EvaluateExpression(`inArray("123", ["222", "233"])`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`inArray("123", ["123", "233"])`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`inArray("123", [123, "233"])`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`inArray("123", [123, "123", "233"])`, nil).Value())

	assert.Equal(t, true, parser.EvaluateExpression(`inArray(123, [123, "233"])`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`inArray(123, ["123", "233"])`, nil).Value())

	assert.Equal(t, false, parser.EvaluateExpression(`inArray("123", [])`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`inArray("1 [23", ["1 [23", "]"])`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`inArray("123", undefined)`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`inArray(undefined, [])`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`inArray(undefined, [undefined, 123])`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`inArray(undefined, undefined)`, nil).Value())

	assert.Equal(t, true, parser.EvaluateExpression(`inArray(mergeSeed("123", "456"), ["123.456", "233"])`, nil).Value())
	assert.Equal(t, false, parser.EvaluateExpression(`inArray("123.456", [mergeSeed("123", "456"), "233"])`, nil).Value()) // THIS CASE IS NOT SUPPORTED

	assert.Equal(t, `07915255d64730d06d2349d11ac3bfd8`, parser.EvaluateExpression(`md5("stam")`, nil).Value())
	assert.Equal(t, `stamstam2`, parser.EvaluateExpression(`concat("stam","stam2")`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`inArray(md5(concat("st","am")), ["07915255d64730d06d2349d11ac3bfd8"]`, nil).Value())
	assert.Equal(t, true, parser.EvaluateExpression(`eq(md5(concat("st",property("am"))), undefined)`, nil).Value())
}

func TestB64d(t *testing.T) {
	parser := roxx.NewParser()
	assert.Equal(t, `stam`, parser.EvaluateExpression(`b64d("c3RhbQ==")`, nil).Value())
	assert.Equal(t, `𩸽`, parser.EvaluateExpression(`b64d("8Km4vQ==")`, nil).Value())
	assert.Equal(t, "", parser.EvaluateExpression(`b64d(\"\")`, nil).Value())
}

func TestTsToNum(t *testing.T) {
	customPropertiesRepository := repositories.NewCustomPropertyRepository()
	parser := roxx.NewParser()
	now := time.Now()
	extensions.NewPropertiesExtensions(parser, customPropertiesRepository, nil).Extend()
	customPropertiesRepository.AddCustomProperty(properties.NewTimeProperty("cp1", now))

	assert.Equal(t, float64(now.UnixMilli())/1000, parser.EvaluateExpression(`tsToNum(property("cp1"))`, nil).Value())

	// wrong type property
	customPropertiesRepository.AddCustomProperty(properties.NewStringProperty("cp2", "notADateTime"))
	assert.Equal(t, nil, parser.EvaluateExpression(`tsToNum(property("cp2"))`, nil).Value())

	// non existent custom property
	assert.Equal(t, nil, parser.EvaluateExpression(`tsToNum(property("cp3"))`, nil).Value())
}
