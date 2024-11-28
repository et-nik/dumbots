package main

import (
	"github.com/et-nik/dumbots/util"
	"github.com/et-nik/metamod-go/engine"
	"github.com/et-nik/metamod-go/vector"
	"math/rand/v2"
)

func (ai *BotAI) makeUpMovementTask(b *Bot) {
	moveTask := ai.createMovementFromRoomGraph(b)
	if moveTask != nil {
		b.MovementTask = moveTask

		return
	}

	moveTask = ai.generateRandomMovement(b)
	if moveTask != nil {
		b.MovementTask = moveTask

		return
	}

	ai.dumbMovement(b)
}

func (ai *BotAI) createMovementFromRoomGraph(b *Bot) BotTask {
	if b.RoomGraph == nil || ai.Globals.Time()-b.RoomGraphGeneratedTime > roomGraphRegenerateTime {
		b.RoomGraphGeneratedTime = ai.Globals.Time()
		ai.makeRoomGraph(b)
	}

	if b.RoomGraph.Len() <= 1 {
		return nil
	}

	if b.Path == nil || ai.Globals.Time()-b.PathGeneratedTime > roomGraphRegenerateTime {
		b.PathGeneratedTime = ai.Globals.Time()
		b.Path = ai.generatePath(b)
		b.CurrentPathIndex = 0
	}

	if len(b.Path) == 0 {
		return nil
	}

	if b.CurrentPathIndex >= len(b.Path) {
		b.RoomGraph = nil
		return ai.createMovementFromRoomGraph(b)
	}

	target := b.Path[b.CurrentPathIndex]

	if !ai.seesPoint(b, target.Coord) {
		return nil
	}

	return &BotTaskMoveToGraphVertex{
		Target:     target,
		TimeToFail: 2,
	}
}

func (ai *BotAI) generateRandomMovement(b *Bot) BotTask {
	generator := b.directionGenerator(b.MaxSpeed * 2)

	if int(ai.Globals.Time()*10.0)%5 == 0 {
		rand.Shuffle(len(generator), func(i, j int) {
			generator[i], generator[j] = generator[j], generator[i]
		})
	}

	for _, f := range generator {
		endTraceCheck := f()

		tr := ai.Funcs.Engine.TraceHull(
			b.Ent.EntVars().Origin(),
			endTraceCheck,
			engine.TraceIgnoreMonsters,
			engine.TraceHullPoint,
			b.Ent.EntVars().PContainingEntity(),
		)

		if tr.EndPos.Distance(b.Ent.EntVars().Origin()) < 100 {
			continue
		}

		target := util.MiddleOfDistance(b.Ent.EntVars().Origin(), tr.EndPos)

		return &BotTaskMoveToPoint{
			Target:     target,
			TimeToFail: 1,
		}
	}

	return nil
}

// probeGraphDirection probes the direction of the graph
func (ai *BotAI) probeGraphDirection(b *Bot, g *Graph, src, dst vector.Vector) *Vertex {
	traceResult := ai.Funcs.Engine.TraceLine(
		src,
		dst,
		engine.TraceIgnoreMonsters,
		b.Ent.EntVars().PContainingEntity(),
	)

	middleDistance := util.MiddleOfDistance(src, traceResult.EndPos)
	newVertex := Vertex{Coord: middleDistance}

	nearExists := false

	for vertex := range g.Iterate {
		if vertex.Coord.Distance(newVertex.Coord) < 100 {
			nearExists = true

			break
		}
	}

	//g.Iterate(func(vertex Vertex) (cont bool) {
	//	if vertex.Coord.Distance(newVertex.Coord) < 100 {
	//		nearExists = true
	//
	//		return false
	//	}
	//
	//	return true
	//})

	if nearExists {
		return nil
	}

	return &newVertex
}

func (ai *BotAI) makeRoomGraph(b *Bot) {
	g := NewGraph()

	vangle := b.Ent.EntVars().VAngle()

	origin := b.Ent.EntVars().Origin()

	playerVertex := Vertex{Name: "pl", Coord: b.Ent.EntVars().Origin()}

	g.AddVertex(playerVertex)

	baseDirections := []vector.Vector{
		origin.Add(util.AnglesToRight(vangle).Mul(roomMaxDistance)),
		origin.Add(util.AnglesToRight(
			vector.Vector{
				vangle[0],
				util.WrapAngle180(vangle[1] + 45),
			},
		).Mul(roomMaxDistance)),
		origin.Add(util.AnglesToRight(vangle).Mul(-roomMaxDistance)),
		origin.Add(util.AnglesToRight(
			vector.Vector{
				vangle[0],
				util.WrapAngle180(vangle[1] + 45),
			},
		).Mul(-roomMaxDistance)),
		origin.Add(util.AnglesToForward(vangle).Mul(roomMaxDistance)),
		origin.Add(util.AnglesToForward(
			vector.Vector{
				vangle[0],
				util.WrapAngle180(vangle[1] + 45),
			},
		).Mul(roomMaxDistance)),
		origin.Add(util.AnglesToForward(vangle).Mul(-roomMaxDistance)),
		origin.Add(util.AnglesToForward(
			vector.Vector{
				vangle[0],
				util.WrapAngle180(vangle[1] + 45),
			},
		).Mul(-roomMaxDistance)),
	}

	for _, baseDirection := range baseDirections {
		v := ai.probeGraphDirection(
			b, g, origin, baseDirection,
		)

		if v == nil {
			continue
		}

		g.AddVertex(*v)
		g.AddEdge(playerVertex, *v)

		nextDirections := []vector.Vector{
			v.Coord.Add(util.AnglesToRight(
				vector.Vector{
					vangle[0],
					util.WrapAngle180(vangle[1] + ai.Funcs.Engine.RandomFloat(-roomRandomizeAngle, roomRandomizeAngle)),
				},
			).Mul(roomMaxDistance)),
			v.Coord.Add(util.AnglesToRight(
				vector.Vector{
					vangle[0],
					util.WrapAngle180(vangle[1] + ai.Funcs.Engine.RandomFloat(-roomRandomizeAngle, roomRandomizeAngle)),
				},
			).Mul(-roomMaxDistance)),
			v.Coord.Add(util.AnglesToForward(
				vector.Vector{
					vangle[0],
					util.WrapAngle180(vangle[1] + ai.Funcs.Engine.RandomFloat(-roomRandomizeAngle, roomRandomizeAngle)),
				},
			).Mul(roomMaxDistance)),
			v.Coord.Add(util.AnglesToForward(
				vector.Vector{
					vangle[0],
					util.WrapAngle180(vangle[1] + ai.Funcs.Engine.RandomFloat(-roomRandomizeAngle, roomRandomizeAngle)),
				},
			).Mul(-roomMaxDistance)),
		}

		for _, next := range nextDirections {

			randomized := ai.Funcs.Engine.RandomFloat(0.8, 1.2)

			nextV := ai.probeGraphDirection(
				b, g, v.Coord, vector.Vector{
					next[0] * randomized,
					next[1] * randomized,
					next[2],
				},
			)

			if nextV == nil {
				continue
			}

			g.AddVertex(*nextV)
			g.AddEdge(*v, *nextV)
		}
	}

	b.RoomGraph = g
}

func (ai *BotAI) generatePath(b *Bot) []Vertex {
	if b.RoomGraph == nil {
		return nil
	}

	var nearestVertex *Vertex
	nearestDist := float32(99999.0)

	var longestVertex *Vertex
	longestDist := float32(0)

	for vertex := range b.RoomGraph.Iterate {
		dist := float32(vertex.Coord.Distance(b.Ent.EntVars().Origin()))

		if dist < nearestDist {
			nearestDist = dist
			nearestVertex = &vertex
		}

		if dist > longestDist {
			longestDist = dist
			longestVertex = &vertex
		}
	}

	if nearestVertex == nil {
		return nil
	}

	if longestVertex == nil {
		return nil
	}

	path, found := b.RoomGraph.AStar(*nearestVertex, *longestVertex)

	if !found {
		return nil
	}

	if longestVertex.Coord.Distance(b.Ent.EntVars().Origin()) < 250 {
		// Bad path
		return nil
	}

	return path
}

func (ai *BotAI) dumbMovement(b *Bot) {
	entVars := b.Ent.EntVars()

	b.ForwardMove = b.MaxSpeed
	b.SideMove = b.MaxSpeed

	entVars.ButtonClear()

	util.ExecuteRandomFunc([]func(){
		func() {
			b.ForwardMove = -b.ForwardMove
		},
		func() {
			b.SideMove = -b.SideMove
		},
	})

	return
}

func (ai *BotAI) RunMovementTask(b *Bot) {
	switch v := b.MovementTask.(type) {
	case *BotTaskMoveToGraphVertex:
		ai.executeBotTaskMoveToGraphVertex(b, v)
	case *BotTaskMoveToPoint:
		ai.executeBotTaskMoveToPoint(b, v)
	}

	if b.MovementTask != nil && b.MovementTask.Finished() {
		b.MovementTask = nil
	}
}

func (ai *BotAI) executeBotTaskMoveToGraphVertex(b *Bot, task *BotTaskMoveToGraphVertex) {
	if task.startedTime == 0 {
		task.SetStartedTime(ai.Globals.Time())
	}

	if !ai.seesPoint(b, task.Target.Coord) {
		task.SetFinished()

		return
	}

	if ai.Globals.Time()-task.startedTime > task.TimeToFail {
		task.SetFinished()

		return
	}

	if task.Target.Coord.Distance(b.Ent.EntVars().Origin()) < 50 {
		b.CurrentPathIndex++

		if b.CurrentPathIndex >= len(b.Path) {
			b.Path = nil
			b.CurrentPathIndex = 0
		}

		task.SetFinished()

		return
	}

	ai.TurnToPoint(b, task.Target.Coord)

	goalAngles := ai.AnglesToPoint(b, task.Target.Coord)

	b.ForwardMove = b.MaxSpeed

	switch {
	case goalAngles[1] > 0:
		b.SideMove = b.MaxSpeed
	case goalAngles[1] < 0:
		b.SideMove = -b.MaxSpeed
	}
}

func (ai *BotAI) executeBotTaskMoveToPoint(b *Bot, task *BotTaskMoveToPoint) {
	if task.startedTime == 0 {
		task.SetStartedTime(ai.Globals.Time())
	}

	if ai.Globals.Time()-task.startedTime > task.TimeToFail {
		task.SetFinished()

		return
	}

	if !ai.seesPoint(b, task.Target) {
		return
	}

	ai.TurnToPoint(b, task.Target)

	goalAngles := ai.AnglesToPoint(b, task.Target)

	b.ForwardMove = b.MaxSpeed

	switch {
	case goalAngles[1] > 0:
		b.SideMove = b.MaxSpeed
	case goalAngles[1] < 0:
		b.SideMove = -b.MaxSpeed
	}

	if b.Ent.EntVars().Origin().Distance(task.Target) < 100 {
		task.SetFinished()

		return
	}
}
