package util

import (
	metamod "github.com/et-nik/metamod-go"
	"github.com/et-nik/metamod-go/engine"
	"github.com/et-nik/metamod-go/vector"
)

func CheckWallInFront(engFuncs *metamod.EngineFuncs, e *metamod.Edict) bool {
	origin := e.EntVars().Origin()

	forward := origin.Add(
		AnglesToForward(e.EntVars().VAngle()).Mul(40.0),
	)

	tr := engFuncs.TraceHull(
		origin,
		forward,
		engine.TraceIgnoreMonsters,
		engine.TraceHullPoint,
		e.EntVars().PContainingEntity(),
	)

	return tr.Fraction < 1.0
}

func CheckWallBehind(engFuncs *metamod.EngineFuncs, e *metamod.Edict) bool {
	origin := e.EntVars().Origin()

	behind := origin.Add(
		AnglesToForward(e.EntVars().VAngle()).Mul(-40.0),
	)

	tr := engFuncs.TraceHull(
		origin,
		behind,
		engine.TraceDontIgnoreMonsters,
		engine.TraceHullPoint,
		e.EntVars().PContainingEntity(),
	)

	return tr.Fraction < 1.0
}

func CheckWallOnLeft(engFuncs *metamod.EngineFuncs, e *metamod.Edict) bool {
	origin := e.EntVars().Origin()

	left := origin.Add(
		AnglesToRight(e.EntVars().VAngle()).Mul(-40.0),
	)

	tr := engFuncs.TraceHull(
		origin,
		left,
		engine.TraceDontIgnoreMonsters,
		engine.TraceHullPoint,
		e.EntVars().PContainingEntity(),
	)

	return tr.Fraction < 1.0
}

func CheckWallOnRight(engFuncs *metamod.EngineFuncs, e *metamod.Edict) bool {
	origin := e.EntVars().Origin()

	right := origin.Add(
		AnglesToRight(e.EntVars().VAngle()).Mul(40.0),
	)

	tr := engFuncs.TraceHull(
		origin,
		right,
		engine.TraceDontIgnoreMonsters,
		engine.TraceHullPoint,
		e.EntVars().PContainingEntity(),
	)

	return tr.Fraction < 1.0
}

const (
	minHeightToJump = 18.0
	maxHeightToJump = 62.0
)

// CheckObstacleInFront checks if there is an obstacle in front of the player
func CheckObstacleInFront(engFuncs *metamod.EngineFuncs, e *metamod.Edict) bool {
	pos := e.EntVars().Origin().Add(e.EntVars().ViewOfs())

	// find ground
	tr := engFuncs.TraceLine(
		pos,
		pos.Add(vector.Vector{0, 0, -100}),
		engine.TraceIgnoreMonsters,
		e.EntVars().PContainingEntity(),
	)

	if tr.Fraction == 1.0 {
		return false
	}

	groundPos := tr.EndPos

	startPos := pos.Add(
		AnglesToForward(e.EntVars().VAngle()).Mul(40.0),
	)

	downPos := vector.Vector{startPos[0], startPos[1], groundPos[2]}

	tr = engFuncs.TraceLine(
		startPos,
		downPos,
		engine.TraceIgnoreMonsters,
		e.EntVars().PContainingEntity(),
	)

	if tr.Fraction == 1.0 {
		return false
	}

	height := tr.EndPos[2] - groundPos[2]

	return height > minHeightToJump && height < maxHeightToJump
}

// CheckObstacleInFront2 checks if there is an obstacle in front of the player
// Second way to check
func CheckObstacleInFront2(engFuncs *metamod.EngineFuncs, e *metamod.Edict) bool {
	eyesPos := e.EntVars().Origin().Add(e.EntVars().ViewOfs())

	// find ground
	tr := engFuncs.TraceLine(
		eyesPos,
		eyesPos.Add(vector.Vector{0, 0, -100}),
		engine.TraceIgnoreMonsters,
		e.EntVars().PContainingEntity(),
	)

	if tr.Fraction == 1.0 {
		return false
	}

	groundPos := tr.EndPos

	startPos := vector.Vector{groundPos[0], groundPos[1], groundPos[2] + minHeightToJump - 1}
	endPos := startPos.Add(
		AnglesToForward(e.EntVars().VAngle()).Mul(80),
	)

	tr1 := engFuncs.TraceLine(
		startPos,
		endPos,
		engine.TraceDontIgnoreMonsters,
		e.EntVars().PContainingEntity(),
	)

	startUpPos := vector.Vector{groundPos[0], groundPos[1], groundPos[2] + maxHeightToJump}
	endEyePos := startUpPos.Add(
		AnglesToForward(e.EntVars().VAngle()).Mul(80),
	)

	tr2 := engFuncs.TraceLine(
		startUpPos,
		endEyePos,
		engine.TraceDontIgnoreMonsters,
		e.EntVars().PContainingEntity(),
	)

	return tr1.Fraction != tr2.Fraction
}

const (
	minHeightToDuck = 37.0
	maxHeightToDuck = 73.0
)

func CheckDuckInFront(engFuncs *metamod.EngineFuncs, e *metamod.Edict) bool {
	pos := e.EntVars().Origin().Add(e.EntVars().ViewOfs())

	// find ground
	tr := engFuncs.TraceLine(
		pos,
		pos.Add(vector.Vector{0, 0, -100}),
		engine.TraceIgnoreMonsters,
		e.EntVars().PContainingEntity(),
	)

	if tr.Fraction == 1.0 {
		return false
	}

	groundPos := tr.EndPos

	startPos := groundPos.Add(
		AnglesToForward(e.EntVars().VAngle()).Mul(40.0),
	)

	upPos := startPos.Add(vector.Vector{0, 0, 100})

	tr = engFuncs.TraceLine(
		startPos,
		upPos,
		engine.TraceIgnoreMonsters,
		e.EntVars().PContainingEntity(),
	)

	if tr.Fraction == 1.0 {
		return false
	}

	height := tr.EndPos[2] - groundPos[2]

	return height >= minHeightToDuck && height < maxHeightToDuck
}
