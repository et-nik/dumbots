package util

import (
	metamod "github.com/et-nik/metamod-go"
	"github.com/et-nik/metamod-go/engine"
)

func IsAlive(e *metamod.Edict) bool {
	entVars := e.EntVars()

	return entVars.DeadFlag() == 0 && entVars.Health() > 0 &&
		!entVars.FlagsHas(engine.EdictFlagNoTarget) && entVars.TakeDamage() != 0 &&
		entVars.Solid() != engine.SolidTypeNot
}
