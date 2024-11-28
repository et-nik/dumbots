package main

import (
	"container/heap"
	"fmt"
	"github.com/chewxy/math32"
	"github.com/et-nik/metamod-go/vector"
)

type Graph struct {
	list map[Vertex][]Vertex
}

func NewGraph() *Graph {
	return &Graph{
		list: make(map[Vertex][]Vertex),
	}
}

func (g *Graph) AddVertex(v Vertex) {
	g.list[v] = make([]Vertex, 0, 4)
}

func (g *Graph) AddEdge(v1, v2 Vertex) {
	g.list[v1] = append(g.list[v1], v2)
	g.list[v2] = append(g.list[v2], v1)
}

func (g *Graph) Exists(v Vertex) bool {
	_, exists := g.list[v]

	return exists
}

func (g *Graph) Iterate(yield func(Vertex) bool) {
	for v := range g.list {
		if !yield(v) {
			return
		}
	}
}

func (g *Graph) Len() int {
	return len(g.list)
}

func (g *Graph) Merge(other *Graph) {
	for v, slice := range other.list {
		if _, exists := g.list[v]; !exists {
			g.list[v] = make([]Vertex, 0, 4)
		}

		for _, vertex := range slice {
			if _, exists := g.list[vertex]; !exists {
				g.list[vertex] = make([]Vertex, 0, 4)
			}

			g.list[v] = append(g.list[v], vertex)
			g.list[vertex] = append(g.list[vertex], v)
		}
	}
}

func (g *Graph) PrintGraph() {
	fmt.Println("Graph connections:")

	for v, slice := range g.list {
		fmt.Printf("%s %v: ", v.Name, v.Coord)

		for _, vertex := range slice {
			fmt.Printf("%s ", vertex.Name)
		}

		fmt.Println()
	}
}

type Vertex struct {
	Name  string
	Coord vector.Vector
	Tags  [8]string
}

// PriorityQueue - очередь с приоритетом для A*
type PriorityQueue []*Node

type Node struct {
	vertex   Vertex
	priority float32
	index    int
}

// Вспомогательные функции для работы с PriorityQueue
func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*Node)
	node.index = len(*pq)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

func heuristic(v1, v2 Vertex) float32 {
	dx := v1.Coord[0] - v2.Coord[0]
	dy := v1.Coord[1] - v2.Coord[1]
	dz := v1.Coord[2] - v2.Coord[2]

	return math32.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (g *Graph) AStar(start, goal Vertex) ([]Vertex, bool) {
	openSet := &PriorityQueue{}
	heap.Init(openSet)

	cameFrom := make(map[Vertex]*Vertex)
	gScore := make(map[Vertex]float32)
	fScore := make(map[Vertex]float32)

	for v := range g.list {
		gScore[v] = math32.Inf(1)
		fScore[v] = math32.Inf(1)
	}
	gScore[start] = 0
	fScore[start] = heuristic(start, goal)

	heap.Push(openSet, &Node{vertex: start, priority: fScore[start]})

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*Node).vertex

		if current == goal {
			return reconstructPath(cameFrom, current), true
		}

		for _, neighbor := range g.list[current] {
			tentativeGScore := gScore[current] + heuristic(current, neighbor)
			if tentativeGScore < gScore[neighbor] {
				cameFrom[neighbor] = &current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = gScore[neighbor] + heuristic(neighbor, goal)

				heap.Push(openSet, &Node{vertex: neighbor, priority: fScore[neighbor]})
			}
		}
	}

	return nil, false
}

// Восстановление пути
func reconstructPath(cameFrom map[Vertex]*Vertex, current Vertex) []Vertex {
	var path []Vertex
	for {
		path = append([]Vertex{current}, path...)
		if parent, exists := cameFrom[current]; exists {
			current = *parent
		} else {
			break
		}
	}
	return path
}
