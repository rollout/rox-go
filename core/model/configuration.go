package model

import "time"

type FetcherError int

const (
	FetcherErrorNoError FetcherError = iota
	FetcherErrorCorruptedJSON
	FetcherErrorEmptyJSON
	FetcherErrorSignatureVerification
	FetcherErrorNetwork
	FetcherErrorMismatchAppKey
	FetcherErrorUnknown
)

type FetcherStatus int

const (
	FetcherStatusAppliedFromEmbedded FetcherStatus = iota
	FetcherStatusAppliedFromLocalStorage
	FetcherStatusAppliedFromNetwork
	FetcherStatusErrorFetchedFailed
)

type ConfigurationFetchedHandler = func(args *ConfigurationFetchedArgs)

type ConfigurationFetchedArgs struct {
	FetcherStatus FetcherStatus
	CreationDate  time.Time
	HasChanges    bool
	ErrorDetails  FetcherError
}

func NewConfigurationFetchedArgs(fetcherStatus FetcherStatus, creationDate time.Time, hasChanges bool) ConfigurationFetchedArgs {
	return ConfigurationFetchedArgs{
		FetcherStatus: fetcherStatus,
		CreationDate:  creationDate,
		HasChanges:    hasChanges,
		ErrorDetails:  FetcherErrorNoError,
	}
}

func NewErrorConfigurationFetchedArgs(errorDetails FetcherError) ConfigurationFetchedArgs {
	return ConfigurationFetchedArgs{
		FetcherStatus: FetcherStatusErrorFetchedFailed,
		CreationDate:  time.Time{},
		HasChanges:    false,
		ErrorDetails:  errorDetails,
	}
}

type Experiment struct {
	Name       string
	Identifier string
	IsArchived bool
	Labels     []string
}

func NewExperiment(experiment *ExperimentModel) *Experiment {
	return &Experiment{
		Name:       experiment.Name,
		Identifier: experiment.ID,
		IsArchived: experiment.IsArchived,
		Labels:     experiment.Labels,
	}
}

type TargetGroupModel struct {
	ID        string
	Condition string
}

func NewTargetGroupModel(id, condition string) *TargetGroupModel {
	return &TargetGroupModel{
		ID:        id,
		Condition: condition,
	}
}
