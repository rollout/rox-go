package extensions

import (
	"crypto/md5"
	"fmt"
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/roxx"
	"math"
)

type ExperimentsExtensions struct {
	parser                 roxx.Parser
	targetGroupsRepository model.TargetGroupRepository
	flagsRepository        model.FlagRepository
	experimentRepository   model.ExperimentRepository
}

func NewExperimentsExtensions(parser roxx.Parser, targetGroupsRepository model.TargetGroupRepository, flagsRepository model.FlagRepository, experimentRepository model.ExperimentRepository) *ExperimentsExtensions {
	return &ExperimentsExtensions{
		parser:                 parser,
		targetGroupsRepository: targetGroupsRepository,
		flagsRepository:        flagsRepository,
		experimentRepository:   experimentRepository,
	}
}

func (e *ExperimentsExtensions) Extend() {
	e.parser.AddOperator("mergeSeed", func(p roxx.Parser, stack *roxx.CoreStack, context context.Context) {
		seed1 := stack.Pop().(string)
		seed2 := stack.Pop().(string)
		stack.Push(fmt.Sprintf("%s.%s", seed1, seed2))
	})

	e.parser.AddOperator("isInPercentage", func(p roxx.Parser, stack *roxx.CoreStack, context context.Context) {
		percentage := stack.Pop().(float64)
		seed := stack.Pop().(string)

		bucket := e.GetBucket(seed)
		stack.Push(bucket <= percentage)
	})

	e.parser.AddOperator("isInPercentageRange", func(p roxx.Parser, stack *roxx.CoreStack, context context.Context) {
		percentageLow := stack.Pop().(float64)
		percentageHigh := stack.Pop().(float64)
		seed := stack.Pop().(string)

		bucket := e.GetBucket(seed)
		stack.Push(percentageLow <= bucket && bucket < percentageHigh)
	})

	e.parser.AddOperator("flagValue", func(p roxx.Parser, stack *roxx.CoreStack, context context.Context) {
		featureFlagIdentifier := stack.Pop().(string)

		result := roxx.FlagFalseValue
		variant := e.flagsRepository.GetFlag(featureFlagIdentifier)
		if variant != nil {
			result = variant.GetValue(context)
		} else {
			flagsExperiment := e.experimentRepository.GetExperimentByFlag(featureFlagIdentifier)
			if flagsExperiment != nil && flagsExperiment.Condition != "" {
				experimentEvalResult := e.parser.EvaluateExpression(flagsExperiment.Condition, context).StringValue()
				if experimentEvalResult != "" {
					result = experimentEvalResult
				}
			}
		}

		stack.Push(result)
	})

	e.parser.AddOperator("isInTargetGroup", func(p roxx.Parser, stack *roxx.CoreStack, context context.Context) {
		targetGroupIdentifier := stack.Pop().(string)

		targetGroup := e.targetGroupsRepository.GetTargetGroup(targetGroupIdentifier)
		if targetGroup == nil {
			stack.Push(false)
		} else {
			isInTargetGroup := e.parser.EvaluateExpression(targetGroup.Condition, context).BoolValue()
			stack.Push(isInTargetGroup)
		}
	})
}

func (e *ExperimentsExtensions) GetBucket(seed string) float64 {
	hasher := md5.New()
	hasher.Write([]byte(seed))
	bytes := hasher.Sum(nil)
	hash := (uint64(bytes[0]) & 0xFF) | ((uint64(bytes[1]) & 0xFF) << 8) | ((uint64(bytes[2]) & 0xFF) << 16) | ((uint64(bytes[3]) & 0xFF) << 24)
	hash &= 0xFFFFFFFF

	bucket := float64(hash) / (math.Pow(2, 32) - 1)
	if bucket == 1 {
		bucket = 0
	}
	return bucket
}
