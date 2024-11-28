package main

import (
	"github.com/et-nik/dumbots/scheduler"
	"testing"
)

//func Test_Experiments(t *testing.T) {
//	v1 := vector.Vector{1, 1, 1}
//	v2 := vector.Vector{2, 2, 1}
//
//	t.Logf("Middle of distance: %v", FindMiddleOfDistance(v1, v2))
//}

//func FindMiddleOfDistance(v1, v2 vector.Vector) vector.Vector {
//	return vector.Vector{
//		(v1[0] + v2[0]) / 2,
//		(v1[1] + v2[1]) / 2,
//		(v1[2] + v2[2]) / 2,
//	}
//}
//
//func Test_FindPath(t *testing.T) {
//	graph := NewGraph()
//
//	vA := Vertex{Name: "A", Coord: vector.Vector{0, 0, 0}}
//	vB := Vertex{Name: "B", Coord: vector.Vector{1, 1, 0}}
//	vC := Vertex{Name: "C", Coord: vector.Vector{2, 0, 0}}
//	vD := Vertex{Name: "D", Coord: vector.Vector{3, 1, 0}}
//	vE := Vertex{Name: "E", Coord: vector.Vector{4, 2, 0}}
//	vF := Vertex{Name: "F", Coord: vector.Vector{5, 2, 0}}
//	vH := Vertex{Name: "H", Coord: vector.Vector{5, 3, 0}}
//
//	graph.AddVertex(vA)
//	graph.AddVertex(vB)
//	graph.AddVertex(vC)
//	graph.AddVertex(vD)
//	graph.AddVertex(vE)
//	graph.AddVertex(vF)
//	graph.AddVertex(vH)
//
//	graph.AddEdge(vA, vB)
//	graph.AddEdge(vA, vC)
//	graph.AddEdge(vB, vD)
//	graph.AddEdge(vC, vD)
//	graph.AddEdge(vD, vE)
//	graph.AddEdge(vD, vH)
//	graph.AddEdge(vF, vH)
//	graph.AddEdge(vF, vE)
//
//	path, found := graph.AStar(vA, vF)
//	if found {
//		for _, v := range path {
//			println(v.Name)
//		}
//	} else {
//		println("Путь не найден")
//	}
//}

func Test_Scheduler(t *testing.T) {
	s := scheduler.NewScheduler()

	s.ScheduleOnce(func() (completed bool, err error) {
		println("Once")
		return true, nil
	}, 1)

	s.ScheduleRepeating(func() (completed bool, err error) {
		println("Repeating")
		return false, nil
	}, 1)

	s.Run()
}
