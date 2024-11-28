package main

import (
	"fmt"
	metamod "github.com/et-nik/metamod-go"
	"github.com/et-nik/metamod-go/engine"
	"github.com/pkg/errors"
)

type BotManager struct {
	Globals *metamod.GlobalVars

	Funcs  *PluginFuncs
	Logger *Logger

	BotAI *BotAI

	// state
	Bots map[int]*Bot
}

func NewBotManager(globals *metamod.GlobalVars, funcs *PluginFuncs, logger *Logger) *BotManager {
	botAI := NewBotAI(globals, funcs, logger)

	return &BotManager{
		Globals: globals,
		Funcs:   funcs,
		Logger:  logger,

		BotAI: botAI,

		Bots: make(map[int]*Bot),
	}
}

func (bm *BotManager) CreateBot() (*Bot, error) {
	bm.Logger.Debug("Creating bot")

	name := BotNames[bm.Funcs.Engine.RandomLong(0, int32(len(BotNames)-1))]

	ent := bm.Funcs.Engine.CreateFakeClient(name)
	if ent == nil {
		return nil, errors.New("failed to create bot")
	}

	bot := NewBot(name, ent, bm.Globals.Time())

	bot.MaxSpeed = bm.Funcs.Engine.CVarGetFloat("sv_maxspeed")

	infoBuffer := bm.Funcs.Engine.GetInfoKeyBuffer(bot.Ent)

	if infoBuffer == nil {
		return nil, errors.New("failed to get info buffer")
	}

	botIndex := bm.Funcs.Engine.IndexOfEdict(bot.Ent)

	bm.Funcs.Engine.SetClientKeyValue(botIndex, *infoBuffer, "_vgui_menus", "0")
	bm.Funcs.Engine.SetClientKeyValue(botIndex, *infoBuffer, "_ah", "0")

	model := BotModels[bm.Funcs.Engine.RandomLong(0, int32(len(BotModels)-1))]
	bm.Funcs.Engine.SetClientKeyValue(botIndex, *infoBuffer, "model", model)

	botAddress := fmt.Sprintf("127.0.0.%d", len(bm.Bots)+101)

	result, reason := bm.Funcs.Game.ClientConnect(bot.Ent, bot.Name, botAddress)
	if result != metamod.Success {
		return nil, errors.New("failed to connect bot: " + reason)
	}
	bm.Funcs.Game.ClientPutInServer(bot.Ent)

	bot.Ent.EntVars().SetFlagsBit(engine.EdictFlagFakeClient)

	bm.Bots[metamod.EntityIndex(bot.Ent)] = bot

	return bot, nil
}

func (bm *BotManager) RemoveAllBots() {
	for _, bot := range bm.Bots {
		bm.Funcs.Game.ClientDisconnect(bot.Ent)
	}

	bm.Bots = make(map[int]*Bot)
}

func (bm *BotManager) Tick() error {
	for _, bot := range bm.Bots {
		err := bm.BotTick(bot)
		if err != nil {
			return errors.WithMessage(err, "failed to tick bot")
		}
	}

	return nil
}

func (bm *BotManager) BotTick(b *Bot) error {
	err := bm.BotAI.Think(b)
	if err != nil {
		return errors.WithMessage(err, "failed to think bot")
	}

	if bm.Globals.Time()-b.PreviousUpdated > 0.1 {
		b.PreviousUpdated = bm.Globals.Time()

		b.MaxSpeed = bm.Funcs.Engine.CVarGetFloat("sv_maxspeed")
	}

	return nil
}

func (bm *BotManager) OnClientDisconnect(e *metamod.Edict) {
	bm.Logger.Debugf("Client disconnect: %d", metamod.EntityIndex(e))

	delete(bm.Bots, metamod.EntityIndex(e))
}

func (bm *BotManager) OnServerDeactivate() {

}
