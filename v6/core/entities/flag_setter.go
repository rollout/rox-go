package entities

import (
	"github.com/rollout/rox-go/v6/core/model"
	"github.com/rollout/rox-go/v6/core/roxx"
	"github.com/rollout/rox-go/v6/core/utils"
)

type FlagSetter struct {
	flagRepository       model.FlagRepository
	parser               roxx.Parser
	experimentRepository model.ExperimentRepository
	impressionInvoker    model.ImpressionInvoker
}

func NewFlagSetter(flagRepository model.FlagRepository, parser roxx.Parser, experimentRepository model.ExperimentRepository, impressionInvoker model.ImpressionInvoker) *FlagSetter {
	fs := &FlagSetter{
		flagRepository:       flagRepository,
		parser:               parser,
		experimentRepository: experimentRepository,
		impressionInvoker:    impressionInvoker,
	}

	fs.flagRepository.RegisterFlagAddedHandler(func(variant model.Variant) {
		exp := fs.experimentRepository.GetExperimentByFlag(variant.Name())
		fs.setFlagData(variant, exp)
	})

	return fs
}

func (fs *FlagSetter) SetExperiments() {
	var flagsWithCondition []string
	for _, exp := range fs.experimentRepository.GetAllExperiments() {
		for _, flagName := range exp.Flags {
			flag := fs.flagRepository.GetFlag(flagName)
			if flag != nil {
				fs.setFlagData(flag, exp)
				flagsWithCondition = append(flagsWithCondition, flagName)
			}
		}
	}

	for _, flag := range fs.flagRepository.GetAllFlags() {
		if !utils.ContainsString(flagsWithCondition, flag.Name()) {
			fs.setFlagData(flag, nil)
		}
	}
}

func (fs *FlagSetter) setFlagData(variant model.Variant, experiment *model.ExperimentModel) {
	variant.(model.InternalVariant).SetForEvaluation(fs.parser, experiment, fs.impressionInvoker)
}
