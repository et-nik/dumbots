package main

import (
	"fmt"
	metamod "github.com/et-nik/metamod-go"
	"github.com/et-nik/metamod-go/engine"
	"strconv"
	"strings"
	"time"
)

// Global development variable for debugging and testing purposes.
var globalDev = struct {
	laser               int              // laser model index
	graphEntities       []*metamod.Edict // entites which showed room graph
	lasersUpdateTaskIDs []uint16         // task id in scheduler for updating lasers
}{}

func botsFn(p *Plugin) func(_ int, _ ...string) {
	return func(argc int, argv ...string) {
		p.Funcs.Engine.ServerPrint("+----------------------------------------------------------+\n")
		p.Funcs.Engine.ServerPrint("|  ID  |    Name    |               Location               |\n")
		p.Funcs.Engine.ServerPrint("+----------------------------------------------------------+\n")

		// print table
		for id, bot := range p.BotManager.Bots {
			p.Funcs.Engine.ServerPrint(fmt.Sprintf("| %4d | %10s | %2f |\n", id, bot.Name, bot.Ent.EntVars().Origin()))
		}

		p.Funcs.Engine.ServerPrint("+----------------------------------------------------------+\n")
	}
}

func addBotFn(p *Plugin) func(argc int, argv ...string) {
	return func(argc int, argv ...string) {
		_, err := p.BotManager.CreateBot()
		if err != nil {
			p.Funcs.Engine.ServerPrint("Failed to create bot: " + err.Error() + "\n")
			return
		}

		p.Funcs.Engine.ServerPrint("Bot added\n")
	}
}

func addMockFn(p *Plugin) func(argc int, argv ...string) {
	return func(argc int, argv ...string) {
		bot, err := p.BotManager.CreateBot()
		if err != nil {
			p.Funcs.Engine.ServerPrint("Failed to create bot: " + err.Error() + "\n")
			return
		}

		bot.IsMock = true

		p.Logger.Debug("Mock added")
	}
}

func setBotAnglesFn(p *Plugin) func(argc int, argv ...string) {
	return func(argc int, argv ...string) {
		if argc < 3 {
			p.Funcs.Engine.ServerPrintf("Usage: %s <id> <pitch> <yaw>\n", argv[0])
			return
		}

		id, err := strconv.Atoi(argv[1])
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid ID\n")

			return
		}

		bot, ok := p.BotManager.Bots[id]
		if !ok {
			p.Funcs.Engine.ServerPrint("Bot not found\n")

			return
		}

		pitch, err := strconv.ParseFloat(argv[2], 32)
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid pitch\n")

			return
		}

		yaw, err := strconv.ParseFloat(argv[3], 32)
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid yaw\n")

			return
		}

		bot.Ent.EntVars().SetAngles([3]float32{float32(pitch), float32(yaw), bot.Ent.EntVars().VAngle()[2]})
	}
}

func moveBotFn(p *Plugin) func(argc int, argv ...string) {
	return func(argc int, argv ...string) {
		if argc < 5 {
			p.Funcs.Engine.ServerPrintf("Usage: %s <id> <x> <y> <z>\n", argv[0])

			return
		}

		id, err := strconv.Atoi(argv[1])
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid ID\n")

			return
		}

		bot, ok := p.BotManager.Bots[id]
		if !ok {
			p.Funcs.Engine.ServerPrint("Bot not found\n")

			return
		}

		x, err := strconv.ParseFloat(argv[2], 32)
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid X\n")

			return
		}

		y, err := strconv.ParseFloat(argv[3], 32)
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid Y\n")

			return
		}

		z, err := strconv.ParseFloat(argv[4], 32)
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid Z\n")

			return
		}

		bot.Ent.EntVars().SetOrigin([3]float32{float32(x), float32(y), float32(z)})

		p.Logger.Debug("Bot moved")
	}
}

func gotoBotFn(p *Plugin) func(argc int, argv ...string) {
	return func(argc int, argv ...string) {
		if argc < 5 {
			p.Funcs.Engine.ServerPrintf("Usage: %s <id> <x> <y> <z>\n", argv[0])

			return
		}

		id, err := strconv.Atoi(argv[1])
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid ID\n")

			return
		}

		bot, ok := p.BotManager.Bots[id]
		if !ok {
			p.Funcs.Engine.ServerPrint("Bot not found\n")

			return
		}

		x, err := strconv.ParseFloat(argv[2], 32)
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid X\n")

			return
		}

		y, err := strconv.ParseFloat(argv[3], 32)
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid Y\n")

			return
		}

		z, err := strconv.ParseFloat(argv[4], 32)
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid Z\n")

			return
		}

		bot.MovementTask = &BotTaskMoveToPoint{
			TimeToFail: 20,
			Target:     [3]float32{float32(x), float32(y), float32(z)},
		}

		p.Logger.Debug("Bot goto task added")
	}
}

func deleteBotFn(p *Plugin) func(argc int, argv ...string) {
	return func(argc int, argv ...string) {
		if argc < 2 {
			p.Funcs.Engine.ServerPrintf("Usage: %s <id>\n", argv[0])
			return
		}

		id, err := strconv.Atoi(argv[1])
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid ID\n")

			return
		}

		_, ok := p.BotManager.Bots[id]
		if !ok {
			p.Funcs.Engine.ServerPrint("Bot not found\n")
			return
		}

		p.Funcs.Game.ClientDisconnect(p.BotManager.Bots[id].Ent)

		p.Logger.Debug("Bot deleted")
	}
}

func deleteAllBotsFn(p *Plugin) func(argc int, argv ...string) {
	return func(argc int, argv ...string) {
		p.BotManager.RemoveAllBots()

		p.Logger.Debug("All bots deleted")
	}
}

func pressButtonFn(p *Plugin) func(argc int, argv ...string) {
	return func(argc int, argv ...string) {
		if argc < 3 {
			p.Funcs.Engine.ServerPrintf("Usage: %s <id> <button>\n", argv[0])
			return
		}

		id, err := strconv.Atoi(argv[1])
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid ID\n")
			return
		}

		bot, ok := p.BotManager.Bots[id]
		if !ok {
			p.Funcs.Engine.ServerPrint("Bot not found\n")
			return
		}

		button := strings.TrimPrefix(strings.ToLower(argv[2]), "in_")

		switch {
		case strings.HasPrefix(button, "at") && strings.HasSuffix(button, "2"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonAttack2)
		case strings.HasPrefix(button, "at"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonAttack)
		case strings.HasPrefix(button, "ju"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonJump)
		case strings.HasPrefix(button, "du"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonDuck)
		case strings.HasPrefix(button, "fo"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonForward)
		case strings.HasPrefix(button, "ba"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonBack)
		case strings.HasPrefix(button, "us"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonUse)
		case strings.HasPrefix(button, "ca"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonCancel)
		case strings.HasPrefix(button, "le"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonLeft)
		case strings.HasPrefix(button, "ri"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonRight)
		case strings.HasPrefix(button, "re"):
			bot.Ent.EntVars().ButtonToggle(engine.InButtonReload)
		default:
			p.Funcs.Engine.ServerPrint("Invalid button\n")

			return
		}

		p.Funcs.Engine.RunPlayerMove(
			bot.Ent,
			bot.Ent.EntVars().Angles(),
			0,
			0,
			0,
			uint16(bot.Ent.EntVars().Button()),
			uint16(bot.Ent.EntVars().Impulse()),
			100,
		)

		p.Logger.Debug("Button pressed")
	}
}

func showBotRoomGraphFn(p *Plugin) func(argc int, argv ...string) {
	return func(argc int, argv ...string) {
		if argc < 2 {
			p.Funcs.Engine.ServerPrintf("Usage: %s <id>\n", argv[0])
			return
		}

		id, err := strconv.Atoi(argv[1])
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid ID\n")
			return
		}

		bot, ok := p.BotManager.Bots[id]
		if !ok {
			p.Funcs.Engine.ServerPrint("Bot not found\n")
			return
		}

		if bot.RoomGraph == nil {
			p.Funcs.Engine.ServerPrint("Room graph not found\n")
			return
		}

		hideBotRoomGraphFn(p)(0, "")

		p.Funcs.Engine.ServerPrint("Room graph:\n")
		bot.RoomGraph.PrintGraph()

		for v := range bot.RoomGraph.Iterate {
			e := p.Funcs.Engine.CreateNamedEntity("env_sprite")
			if metamod.IsNullEntity(e) {
				p.Logger.Error("Failed to create env_sprite\n")

				continue
			}

			globalDev.graphEntities = append(globalDev.graphEntities, e)

			e.EntVars().SetOrigin([3]float32{v.Coord[0], v.Coord[1], v.Coord[2]})

			if v.Name == "pl" {
				e.EntVars().SetModel("sprites/glow02.spr")
			} else {
				e.EntVars().SetModel("sprites/glow01.spr")
			}

			e.EntVars().SetRenderMode(5)
			e.EntVars().SetRenderAmt(255)
			e.EntVars().SetScale(1.0)
			e.EntVars().SetFrameRate(0)
			e.EntVars().SetEffects(0)

			result := p.Funcs.Game.Spawn(e)
			if result != 0 {
				p.Logger.Error("Failed to spawn env_sprite\n")

				continue
			}

			continue
		}

		createBeam := func(v1, v2 Vertex) {
			p.Funcs.Engine.MessageBegin(engine.MessageDestBroadcast, engine.SvcTempEntity, nil, nil)
			p.Funcs.Engine.MessageWriteByte(0)            // TE_BEAMPOINTS
			p.Funcs.Engine.MessageWriteCoord(v1.Coord[0]) // start.x
			p.Funcs.Engine.MessageWriteCoord(v1.Coord[1]) // start.y
			p.Funcs.Engine.MessageWriteCoord(v1.Coord[2]) // start.z
			p.Funcs.Engine.MessageWriteCoord(v2.Coord[0]) // end.x
			p.Funcs.Engine.MessageWriteCoord(v2.Coord[1]) // end.y
			p.Funcs.Engine.MessageWriteCoord(v2.Coord[2]) // end.z
			p.Funcs.Engine.MessageWriteShort(globalDev.laser)
			p.Funcs.Engine.MessageWriteByte(1)   // start frame
			p.Funcs.Engine.MessageWriteByte(10)  // frame rate
			p.Funcs.Engine.MessageWriteByte(10)  // life in 0.1's
			p.Funcs.Engine.MessageWriteByte(10)  // width
			p.Funcs.Engine.MessageWriteByte(2)   // noise
			p.Funcs.Engine.MessageWriteByte(255) // red
			p.Funcs.Engine.MessageWriteByte(255) // green
			p.Funcs.Engine.MessageWriteByte(255) // blue
			p.Funcs.Engine.MessageWriteByte(200) // brightness
			p.Funcs.Engine.MessageWriteByte(10)  // speed
			p.Funcs.Engine.MessageEnd()
		}

		for v1, slice := range bot.RoomGraph.list {
			for _, v2 := range slice {
				createBeam(v1, v2)

				taskID := p.Scheduler.ScheduleRepeating(1*time.Second, func() (completed bool, err error) {
					createBeam(v1, v2)

					return false, nil
				})
				globalDev.lasersUpdateTaskIDs = append(globalDev.lasersUpdateTaskIDs, taskID)
			}
		}
	}
}

func hideBotRoomGraphFn(p *Plugin) func(_ int, _ ...string) {
	return func(argc int, argv ...string) {
		if len(globalDev.graphEntities) > 0 {
			for _, e := range globalDev.graphEntities {
				p.Funcs.Engine.RemoveEntity(e)
			}

			globalDev.graphEntities = nil
		}

		if len(globalDev.lasersUpdateTaskIDs) != 0 {
			for _, tid := range globalDev.lasersUpdateTaskIDs {
				p.Scheduler.Cancel(tid)
			}

			globalDev.lasersUpdateTaskIDs = nil
		}
	}
}

func botSayFn(p *Plugin) func(_ int, _ ...string) {
	return func(argc int, argv ...string) {
		if argc < 3 {
			p.Funcs.Engine.ServerPrintf("Usage: %s <id> <message>\n", argv[0])

			return
		}

		id, err := strconv.Atoi(argv[1])
		if err != nil {
			p.Funcs.Engine.ServerPrint("Invalid ID\n")

			return
		}

		bot, ok := p.BotManager.Bots[id]
		if !ok {
			p.Funcs.Engine.ServerPrint("Bot not found\n")

			return
		}

		message := strings.Join(argv[2:], " ")

		p.HookedCommand.Argv = []string{"say", message}
		p.Funcs.Game.ClientCommand(bot.Ent)
		p.HookedCommand.Argv = nil

		p.Logger.Debug("Bot said")
	}
}
