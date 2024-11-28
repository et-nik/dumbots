package main

import (
	"github.com/et-nik/dumbots/scheduler"
	metamod "github.com/et-nik/metamod-go"
	"github.com/pkg/errors"
)

type PluginFuncs struct {
	Engine   *metamod.EngineFuncs
	MetaUtil *metamod.MUtilFuncs
	Game     *metamod.GameDLLFuncs
}

type Plugin struct {
	Funcs   *PluginFuncs
	Logger  *Logger
	Globals *metamod.GlobalVars
	CVars   struct {
		MaxBots *metamod.CVar
	}

	BotManager *BotManager

	HookedCommand struct {
		Argv []string
	}

	Scheduler *scheduler.Scheduler
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

// Init initializes the plugin.
// It initializes the state, logger, bot manager, and scheduler.
// This method should be called only once during the plugin's lifecycle, in GameDLLInit callback.
func (p *Plugin) Init() error {
	p.Logger = NewLogger(p.Funcs.MetaUtil)

	p.BotManager = NewBotManager(p.Globals, p.Funcs, p.Logger)

	p.Scheduler = scheduler.NewScheduler()

	p.Logger.Debug("Plugin initialized")

	return nil
}

func (p *Plugin) Tick() error {
	err := p.Scheduler.Run()
	if err != nil {
		return errors.WithMessage(err, "failed to run scheduler")
	}

	if len(p.BotManager.Bots) < p.CVars.MaxBots.Int() {
		_, err := p.BotManager.CreateBot()
		if err != nil {
			return errors.WithMessage(err, "failed to create bot")
		}

		p.Logger.Debug("Bot created")
	}

	err = p.BotManager.Tick()
	if err != nil {
		return errors.WithMessage(err, "failed to tick bot manager")
	}

	return nil
}

func (p *Plugin) OnClientDisconnect(e *metamod.Edict) {
	p.BotManager.OnClientDisconnect(e)
}

func (p *Plugin) Shutdown() {
	p.BotManager.RemoveAllBots()

	p.Logger.Debug("Plugin shutdown")
}
