//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	internal "github.com/osousa/resticky/internal"
)

var appSet = wire.NewSet(
	ProvideApp,
)

var commonSet = wire.NewSet(
	internal.ProvideConnManager,
	internal.ProvideServer,
	internal.ProvideZap,
)

var megaSet = wire.NewSet(
	commonSet,
	appSet,
)

func InitializeApp(
	ctx context.Context,
	cfg internal.ServerConfig,
	cnn internal.ConnManagerConfig,
) (*App, error) {

	// Set up the injection runtime. App top-level injector.
	wire.Build(megaSet)
	return &App{}, nil
}
