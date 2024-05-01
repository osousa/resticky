package main

import "github.com/osousa/resticky/internal"

// ╔════════════════════════════════╗
// ║             Config             ║
// ╚════════════════════════════════╝

type AppConfig struct {
	Server  internal.ServerConfig
	ConnMan internal.ConnManagerConfig
}

func ProvideAppConfig() *AppConfig {
	return &AppConfig{
		Server:  internal.NewServerConfig(),
		ConnMan: internal.NewConnManagerConfig(),
	}
}
