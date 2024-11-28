package main

import (
	"github.com/chewxy/math32"
	"github.com/et-nik/dumbots/util"
	metamod "github.com/et-nik/metamod-go"
	"github.com/et-nik/metamod-go/engine"
	"github.com/et-nik/metamod-go/vector"
	"math"
)

type BotConditions struct {
	ConditionsCheckTime     float32
	WallConditionsCheckTime float32

	// Wall
	WallInFront bool
	WallOnLeft  bool
	WallOnRight bool
	WallBehind  bool

	// Obstacle
	ObstacleInFront bool
	DuckInFront     bool
}

type Bot struct {
	Name  string
	Model string

	IsMock bool

	Ent *metamod.Edict

	PrimaryTask  BotTask
	MovementTask BotTask

	Conditions BotConditions

	// Temporary private graph.
	// Using to navigate in the current room or area.
	RoomGraph              *Graph
	RoomGraphGeneratedTime float32

	CurrentPathIndex  int
	Path              []Vertex
	PathGeneratedTime float32

	// Enemy
	Enemy          *metamod.Edict
	EnemyAngles    vector.Vector
	EnemyFoundTime float32

	PreviousTickTime float32

	PreviousMoveChangeTime float32
	PreviousUpdated        float32

	MaxSpeed float32

	// Move
	ForwardMove float32
	SideMove    float32

	VAngle [3]float32

	OldButtons engine.InButtonFlag
	OldVangle  [3]float32
}

func NewBot(name string, ent *metamod.Edict, time float32) *Bot {
	return &Bot{
		Name:             name,
		Ent:              ent,
		PreviousTickTime: time,
	}
}

func (b *Bot) ComputeMsec(currentTime float32) uint16 {
	return uint16(min(int32(math.Round(float64((currentTime-b.PreviousTickTime)*1000.0))), 255))
}

func (b *Bot) EyesPos() vector.Vector {
	return b.Ent.EntVars().Origin().Add(b.Ent.EntVars().ViewOfs())
}

func (b *Bot) Forward() vector.Vector {
	return util.AnglesToForward(b.Ent.EntVars().VAngle())
}

func (b *Bot) Right() vector.Vector {
	return util.AnglesToRight(b.VAngle)
}

func (b *Bot) directionGenerator(distance float32) []func() vector.Vector {
	return []func() vector.Vector{
		func() vector.Vector {
			return b.Ent.EntVars().Origin().Add(
				b.Forward().Mul(math32.Abs(distance)),
			)
		},
		func() vector.Vector {
			return b.Ent.EntVars().Origin().Add(
				b.Right().Mul(math32.Abs(distance)),
			)
		},
		func() vector.Vector {
			return b.Ent.EntVars().Origin().Add(
				b.Forward().Mul(-math32.Abs(distance)),
			)
		},
		func() vector.Vector {
			return b.Ent.EntVars().Origin().Add(
				b.Right().Mul(-math32.Abs(distance)),
			)
		},
	}

}
