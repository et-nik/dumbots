package main

import (
	"github.com/chewxy/math32"
	"github.com/et-nik/dumbots/util"
	metamod "github.com/et-nik/metamod-go"
	"github.com/et-nik/metamod-go/engine"
	"github.com/et-nik/metamod-go/vector"
)

const (
	kValue = 1.0 / 3.0

	roomGraphRegenerateTime = 5
	roomMaxDistance         = 350.0
	roomRandomizeAngle      = 30.0

	curvatureFactor = 0.8
)

type BotAI struct {
	Globals *metamod.GlobalVars

	Funcs  *PluginFuncs
	Logger *Logger
}

func NewBotAI(globals *metamod.GlobalVars, funcs *PluginFuncs, logger *Logger) *BotAI {
	return &BotAI{
		Globals: globals,
		Funcs:   funcs,
		Logger:  logger,
	}
}

func (ai *BotAI) Think(b *Bot) error {
	if ai.Globals.Time()-b.PreviousTickTime < 0.04 {
		return nil
	}

	if metamod.IsNullEntity(b.Ent) {
		return nil
	}

	botEntVars := b.Ent.EntVars()

	if !botEntVars.IsValid() {
		return nil
	}

	botEntVars.ButtonClear()

	msec := b.ComputeMsec(ai.Globals.Time())
	b.PreviousTickTime = ai.Globals.Time()

	if ai.Globals.Time()-b.Conditions.ConditionsCheckTime > 0.1 {
		b.Conditions.ConditionsCheckTime = ai.Globals.Time()

		ai.updateConditions(b)
	}

	if b.IsMock {
		vangle := botEntVars.Angles()

		ai.Funcs.Engine.RunPlayerMove(
			b.Ent,
			vangle,
			b.ForwardMove,
			b.SideMove,
			0,
			uint16(botEntVars.Button()),
			uint16(botEntVars.Impulse()),
			msec,
		)

		return nil
	}

	if botEntVars.IdealPitch() == 0 {
		botEntVars.SetIdealPitch(botEntVars.VAngle()[0])
	}

	if botEntVars.IdealYaw() == 0 {
		botEntVars.SetIdealYaw(botEntVars.VAngle()[1])
	}

	if botEntVars.PitchSpeed() == 0 {
		botEntVars.SetPitchSpeed(1)
	}

	if botEntVars.YawSpeed() == 0 {
		botEntVars.SetYawSpeed(1)
	}

	if b.PrimaryTask != nil {
		ai.RunPrimaryTask(b)
	}

	if b.MovementTask != nil {
		ai.RunMovementTask(b)
	} else {
		ai.makeUpMovementTask(b)
	}

	if b.Enemy == nil || metamod.IsNullEntity(b.Enemy) || !metamod.IsAlive(b.Enemy) {
		b.Enemy = ai.findEnemy(b)
		b.EnemyFoundTime = ai.Globals.Time()
	}

	if b.Enemy != nil {
		b.EnemyAngles = ai.AnglesToEntity(b, b.Enemy)
	}

	vangle := botEntVars.Angles()

	if b.Conditions.WallInFront && b.ForwardMove > 0 {
		b.ForwardMove = -b.ForwardMove

		botEntVars.SetYawSpeed(botEntVars.YawSpeed() * 0.25)

		util.ExecuteRandomFunc([]func(){
			func() {
				botEntVars.SetIdealYaw(wrapAngle180(-botEntVars.IdealYaw()))
			},
			func() {
				botEntVars.SetIdealYaw(wrapAngle180(botEntVars.IdealYaw() + 90))
			},
			func() {
				botEntVars.SetIdealYaw(wrapAngle180(botEntVars.IdealYaw() - 90))
			},
		})
	}

	if b.Conditions.WallBehind && b.ForwardMove < 0 {
		b.ForwardMove = -b.ForwardMove
	}

	if b.Conditions.WallOnLeft && b.SideMove < 0 || b.Conditions.WallOnRight && b.SideMove > 0 {
		b.SideMove = -b.SideMove
	}

	if b.Conditions.ObstacleInFront {
		if ai.Funcs.Engine.RandomLong(0, 3) == 0 {
			botEntVars.SetButtonBit(engine.InButtonJump)
			botEntVars.SetButtonBit(engine.InButtonDuck)

			b.ForwardMove = math32.Abs(b.ForwardMove)
			b.SideMove = 0

			b.PreviousMoveChangeTime = ai.Globals.Time()
		}
	}

	if b.Conditions.DuckInFront {
		botEntVars.SetButtonBit(engine.InButtonDuck)

		b.ForwardMove = math32.Abs(b.ForwardMove)
		b.SideMove = 0

		b.PreviousMoveChangeTime = ai.Globals.Time()
	}

	vangle[0] = util.WrapPitch(
		util.CurveValue(
			util.SmoothAngle(vangle[0], botEntVars.IdealPitch(), botEntVars.PitchSpeed()),
			curvatureFactor*2,
		),
	)
	vangle[1] = util.WrapAngle180(
		util.CurveValue(
			util.SmoothAngle(vangle[1], botEntVars.IdealYaw(), botEntVars.YawSpeed()),
			curvatureFactor,
		),
	)

	ai.Funcs.Engine.RunPlayerMove(
		b.Ent,
		vangle,
		b.ForwardMove,
		b.SideMove,
		0,
		uint16(botEntVars.Button()),
		uint16(botEntVars.Impulse()),
		msec,
	)

	botEntVars.SetAngles(
		[3]float32{vangle[0] * kValue, vangle[1], botEntVars.Angles()[2]},
	)

	botEntVars.SetVAngle([3]float32{vangle[0], vangle[1], vangle[2]})

	b.OldButtons = botEntVars.Button()
	b.OldVangle = vangle

	return nil
}

func (ai *BotAI) updateConditions(b *Bot) {
	if b.Conditions.WallConditionsCheckTime == 0 || ai.Globals.Time()-b.Conditions.WallConditionsCheckTime > 1 {
		b.Conditions.WallInFront = util.CheckWallInFront(ai.Funcs.Engine, b.Ent)
		b.Conditions.WallOnLeft = util.CheckWallOnLeft(ai.Funcs.Engine, b.Ent)
		b.Conditions.WallOnRight = util.CheckWallOnRight(ai.Funcs.Engine, b.Ent)
		b.Conditions.WallBehind = util.CheckWallBehind(ai.Funcs.Engine, b.Ent)
	}

	b.Conditions.ObstacleInFront = util.CheckObstacleInFront(ai.Funcs.Engine, b.Ent)
	if !b.Conditions.ObstacleInFront {
		b.Conditions.ObstacleInFront = util.CheckObstacleInFront2(ai.Funcs.Engine, b.Ent)
	}

	b.Conditions.DuckInFront = util.CheckDuckInFront(ai.Funcs.Engine, b.Ent)
}

func (ai *BotAI) makeUpPrimaryTask(b *Bot) {
	ai.makeRoomGraph(b)
}

func (ai *BotAI) RunPrimaryTask(b *Bot) {
	switch v := b.PrimaryTask.(type) {
	case *BotTaskMoveToPoint:
		ai.executeBotTaskMoveToPoint(b, v)
	}

	if b.PrimaryTask != nil && b.PrimaryTask.Finished() {
		b.PrimaryTask = nil
	}
}

func (ai *BotAI) findEnemy(b *Bot) *metamod.Edict {
	e := ai.FindClosestPlayer(b)

	if metamod.IsNullEntity(e) {
		return nil
	}

	return e
}

func (ai *BotAI) FindClosestPlayer(b *Bot) *metamod.Edict {
	var pPlayer *metamod.Edict
	closestDist := float32(99999.0)

	for i := 1; i <= ai.Globals.MaxClients(); i++ {
		pEdict := ai.Funcs.Engine.EntityOfEntIndex(i)

		if pEdict == nil || !pEdict.EntVars().IsValid() || metamod.EntityIndex(pEdict) == metamod.EntityIndex(b.Ent) {
			continue
		}

		if !metamod.IsAlive(pEdict) {
			continue
		}

		dist := float32(pEdict.EntVars().Origin().Distance(b.Ent.EntVars().Origin()))

		if dist < closestDist {
			closestDist = dist
			pPlayer = pEdict
		}
	}

	return pPlayer
}

func (ai *BotAI) GetDirectionToPlayer(b *Bot, pPlayer *metamod.Edict) [3]float32 {
	direction := pPlayer.EntVars().Origin().Sub(b.Ent.EntVars().Origin())

	direction[2] = 0

	return direction.Normalize()
}

func (ai *BotAI) AnglesToEntity(b *Bot, ent *metamod.Edict) vector.Vector {
	return ai.AnglesToPoint(b, ent.EntVars().Origin())
}

func (ai *BotAI) AnglesToPoint(b *Bot, point vector.Vector) vector.Vector {
	delta := point.Sub(b.Ent.EntVars().Origin())

	if delta.IsZero() {
		return vector.Vector{}
	}

	angles := ai.Funcs.Engine.VecToAngles(delta)

	return vector.Vector{util.WrapPitch(angles[0]), util.WrapAngle180(angles[1]), 0}
}

func (ai *BotAI) AimAtEntity(b *Bot, ent *metamod.Edict) {
	ai.AimAtPoint(b, ent.EntVars().Origin())
}

func (ai *BotAI) AimAtPoint(b *Bot, point vector.Vector) {
	angles := ai.AnglesToPoint(b, point)

	angles[0] = util.WrapPitch(angles[0])
	angles[1] = util.WrapAngle180(angles[1])

	b.Ent.EntVars().SetIdealPitch(angles[0])
	b.Ent.EntVars().SetIdealYaw(angles[1])

	pitchSpeed := ai.Funcs.Engine.RandomFloat(1, 10)
	pm := math32.Abs(angles[0])
	if pm >= 15 {
		pitchSpeed *= 3
	}

	yawSpeed := ai.Funcs.Engine.RandomFloat(2, 6)
	ym := math32.Abs(angles[1])
	switch {
	case ym > 45:
		yawSpeed *= 4

	case ym > 20:
		yawSpeed *= 3

	case ym > 10:
		yawSpeed *= 2
	}

	b.Ent.EntVars().SetPitchSpeed(pitchSpeed)
	b.Ent.EntVars().SetYawSpeed(yawSpeed)
}

func (ai *BotAI) TurnToPoint(b *Bot, point vector.Vector) {
	angles := ai.AnglesToPoint(b, point)

	angles[1] = util.WrapAngle180(angles[1])

	b.Ent.EntVars().SetIdealYaw(angles[1])

	yawSpeed := ai.Funcs.Engine.RandomFloat(2, 6)
	ym := math32.Abs(angles[1])
	switch {
	case ym > 45:
		yawSpeed *= 4

	case ym > 20:
		yawSpeed *= 3

	case ym > 10:
		yawSpeed *= 2
	}

	b.Ent.EntVars().SetYawSpeed(yawSpeed)
}

func (ai *BotAI) SeesEntity(b *Bot, ent *metamod.Edict, fromBody bool) bool {
	return ai.seesPoint(b, ent.EntVars().Origin())
}

func (ai *BotAI) seesPoint(b *Bot, point vector.Vector) bool {
	start := b.EyesPos()

	tr := ai.Funcs.Engine.TraceLine(
		start,
		point,
		engine.TraceIgnoreMonsters,
		b.Ent.EntVars().PContainingEntity(),
	)

	return tr.Fraction >= 1.0
}
