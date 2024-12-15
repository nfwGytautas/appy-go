package appy_runtime

import (
	appy_config "github.com/nfwGytautas/appy-go/config"
	appy_logger "github.com/nfwGytautas/appy-go/logger"
)

type HookFn func(ctx *AppyContext)

type AppyContext struct {
	config appy_config.AppyConfig
	hook   HookFn

	HttpEngine *HttpEngine
}

func Initialize(config appy_config.AppyConfig, hook HookFn) *AppyContext {
	appy_logger.Logger().Debug("Appy initializing")
	return &AppyContext{
		config: config,
		hook:   hook,
	}
}

func (ac *AppyContext) Takeover() {
	appy_logger.Logger().Debug("Appy taking over")

	// HTTP
	if ac.config.Http != nil {
		ac.HttpEngine = createHttpEngine(ac.config.Http)
	}

	// Hook
	if ac.hook != nil {
		ac.hook(ac)
	}

	// Run
	if ac.HttpEngine != nil {
		appy_logger.Logger().Debug("Running HTTP server")
		err := ac.HttpEngine.run()
		if err != nil {
			appy_logger.Logger().Error("Failed to run HTTP server")
			panic(err)
		}
	}
}
