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

	VariantWithContext server.RoxString

	Variant            server.RoxString
	VariantOverwritten server.RoxString

	FlagForDependency             server.RoxFlag
	FlagColorsForDependency       server.RoxString
	FlagDependent                 server.RoxFlag
	FlagColorDependentWithContext server.RoxString
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

	VariantWithContext: server.NewRoxString("red", []string{"red", "blue", "green"}),

	Variant:            server.NewRoxString("red", []string{"red", "blue", "green"}),
	VariantOverwritten: server.NewRoxString("red", []string{"red", "blue", "green"}),

	FlagForDependency:             server.NewRoxFlag(false),
	FlagColorsForDependency:       server.NewRoxString("White", []string{"White", "Blue", "Green", "Yellow"}),
	FlagDependent:                 server.NewRoxFlag(false),
	FlagColorDependentWithContext: server.NewRoxString("White", []string{"White", "Blue", "Green", "Yellow"}),
}
