package main

import "github.com/rollout/rox-go/v6/server"

type Container struct {
	BoolDefaultFalse server.RoxFlag
	BoolDefaultTrue  server.RoxFlag
}

var container = &Container{
	BoolDefaultFalse: server.NewRoxFlag(false),
	BoolDefaultTrue:  server.NewRoxFlag(true),
}
