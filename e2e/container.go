package e2e

import "github.com/rollout/rox-go/server"

type Container struct {
	SimpleFlag            server.RoxFlag
	SimpleFlagOverwritten server.RoxFlag

	FlagForImpression                         server.RoxFlag
	FlagForImpressionWithExperimentAndContext server.RoxFlag

	FlagCustomProperties server.RoxFlag

	FlagTargetGroupsAll  server.RoxFlag
	FlagTargetGroupsAny  server.RoxFlag
	FlagTargetGroupsNone server.RoxFlag

	VariantWithContext server.RoxVariant

	Variant            server.RoxVariant
	VariantOverwritten server.RoxVariant

	FlagForDependency             server.RoxFlag
	FlagColorsForDependency       server.RoxVariant
	FlagDependent                 server.RoxFlag
	FlagColorDependentWithContext server.RoxVariant
}

var container = &Container{
	SimpleFlag:            server.NewRoxFlag(true),
	SimpleFlagOverwritten: server.NewRoxFlag(true),

	FlagForImpression:                         server.NewRoxFlag(false),
	FlagForImpressionWithExperimentAndContext: server.NewRoxFlag(false),

	FlagCustomProperties: server.NewRoxFlag(false),

	FlagTargetGroupsAll:  server.NewRoxFlag(false),
	FlagTargetGroupsAny:  server.NewRoxFlag(false),
	FlagTargetGroupsNone: server.NewRoxFlag(false),

	VariantWithContext: server.NewRoxVariant("red", []string{"red", "blue", "green"}),

	Variant:            server.NewRoxVariant("red", []string{"red", "blue", "green"}),
	VariantOverwritten: server.NewRoxVariant("red", []string{"red", "blue", "green"}),

	FlagForDependency:             server.NewRoxFlag(false),
	FlagColorsForDependency:       server.NewRoxVariant("White", []string{"White", "Blue", "Green", "Yellow"}),
	FlagDependent:                 server.NewRoxFlag(false),
	FlagColorDependentWithContext: server.NewRoxVariant("White", []string{"White", "Blue", "Green", "Yellow"}),
}
