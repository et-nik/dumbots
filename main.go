package main

import (
	"fmt"
	metamod "github.com/et-nik/metamod-go"
	"strings"
)

func init() {
	err := metamod.SetPluginInfo(&metamod.PluginInfo{
		InterfaceVersion: metamod.MetaInterfaceVersion,
		Name:             "DumBots",
		Version:          "0.9.0",
		Date:             "2024-11-28",
		Author:           "KNiK",
		Url:              "https://github.com/et-nik/dumbots",
		LogTag:           "DumBots",
		Loadable:         metamod.PluginLoadTimeStartup,
		Unloadable:       metamod.PluginLoadTimeAnyTime,
	})
	if err != nil {
		panic(err)
	}

	plugin := NewPlugin()

	err = metamod.SetMetaCallbacks(&metamod.MetaCallbacks{
		MetaQuery:  metaQueryFn(plugin),
		MetaAttach: metaAttachFn(plugin),
	})
	if err != nil {
		panic(err)
	}

	err = metamod.SetApiCallbacks(&metamod.APICallbacks{
		GameDLLInit: func() metamod.APICallbackResult {
			err := plugin.Init()
			if err != nil {
				fmt.Println("Failed to init plugin: ", err.Error())

				return metamod.APICallbackResultHandled
			}

			return metamod.APICallbackResultHandled
		},
		Spawn: func(e *metamod.Edict) (metamod.APICallbackResult, int) {
			plugin.Funcs.Engine.PrecacheModel("sprites/glow01.spr")
			plugin.Funcs.Engine.PrecacheModel("sprites/glow02.spr")
			plugin.Funcs.Engine.PrecacheModel("sprites/glow03.spr")

			globalDev.laser = plugin.Funcs.Engine.PrecacheModel("sprites/lgtning.spr")

			return metamod.APICallbackResultHandled, 0
		},
		StartFrame: func() metamod.APICallbackResult {
			err = plugin.Tick()
			if err != nil {
				plugin.Logger.Errorf("Failed to tick: %s", err.Error())

				return metamod.APICallbackResultHandled
			}

			return metamod.APICallbackResultHandled
		},
		ClientDisconnect: func(e *metamod.Edict) metamod.APICallbackResult {
			plugin.OnClientDisconnect(e)

			return metamod.APICallbackResultHandled
		},
		ServerDeactivate: func() metamod.APICallbackResult {
			plugin.Shutdown()

			return metamod.APICallbackResultHandled
		},
	})
	if err != nil {
		panic(err)
	}

	err = metamod.SetEngineHooks(&metamod.EngineHooks{
		CmdArgv: func(arg int) (metamod.EngineHookResult, string) {
			if plugin.HookedCommand.Argv == nil {
				return metamod.EngineHookResultIgnored, ""
			}

			if arg < 0 || arg >= len(plugin.HookedCommand.Argv) {
				return metamod.EngineHookResultSupercede, ""
			}

			return metamod.EngineHookResultSupercede, plugin.HookedCommand.Argv[arg]
		},
		CmdArgs: func() (metamod.EngineHookResult, string) {
			if plugin.HookedCommand.Argv == nil {
				return metamod.EngineHookResultIgnored, ""
			}

			if len(plugin.HookedCommand.Argv) == 1 {
				return metamod.EngineHookResultSupercede, ""
			}

			b := strings.Builder{}
			b.Grow(32)

			for i, arg := range plugin.HookedCommand.Argv[1:] {
				if i > 0 {
					b.WriteByte(' ')
				}

				if strings.Contains(arg, " ") {
					b.WriteByte('"')
					b.WriteString(arg)
					b.WriteByte('"')

					continue
				} else {
					b.WriteString(arg)
				}
			}

			return metamod.EngineHookResultSupercede, b.String()
		},
		CmdArgc: func() (metamod.EngineHookResult, int) {
			if plugin.HookedCommand.Argv == nil {
				return metamod.EngineHookResultIgnored, 0
			}

			return metamod.EngineHookResultSupercede, len(plugin.HookedCommand.Argv)
		},
	})
	if err != nil {
		panic(err)
	}
}

func metaQueryFn(p *Plugin) func() int {
	return func() int {
		var err error

		if p.Funcs == nil {
			p.Funcs = &PluginFuncs{}
		}

		p.Funcs.Engine, err = metamod.GetEngineFuncs()
		if err != nil {
			fmt.Println("Failed to get engine funcs: ", err.Error())

			return 0
		}

		p.Funcs.MetaUtil, err = metamod.GetMetaUtilFuncs()
		if err != nil {
			fmt.Println("Failed to get meta util funcs: ", err.Error())

			return 0
		}

		p.Globals = metamod.GetGlobalVars()

		return 1
	}
}

func metaAttachFn(p *Plugin) func(now int) int {
	return func(now int) int {
		var err error

		p.Funcs.Game, err = metamod.GetGameDLLFuncs()
		if err != nil {
			fmt.Println("Failed to get game dll funcs: ", err.Error())

			return 0
		}

		p.Funcs.Engine.AddServerCommand("dbt_list", botsFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_add", addBotFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_add_mock", addMockFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_move", moveBotFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_goto", gotoBotFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_angles", setBotAnglesFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_rm", deleteBotFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_rmall", deleteAllBotsFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_press", pressButtonFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_room_graph", showBotRoomGraphFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_hide_room_graph", hideBotRoomGraphFn(p))
		p.Funcs.Engine.AddServerCommand("dbt_say", botSayFn(p))

		p.CVars.MaxBots = p.Funcs.Engine.CVarRegister(
			metamod.NewCVar("dbt_max_bots", "1", metamod.CVarServer|metamod.CVarExtdll),
		)

		return 1
	}
}

func main() {}
