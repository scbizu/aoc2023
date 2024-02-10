package main

import (
	"container/list"
	"context"
	"fmt"
	"sort"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/matrix"
	"github.com/magejiCoder/magejiAoc/queue"
	"github.com/magejiCoder/magejiAoc/stack"
	"github.com/magejiCoder/set"
)

const (
	line = 141
	col  = 141
)

type trails struct {
	start   grid.Point
	m       matrix.Matrix[byte]
	end     grid.Point
	visited map[grid.Point]int
	graph   map[grid.Point][]graphNode
}

type graphNode struct {
	head, end grid.Point
	weight    int
}

func (g *graphNode) String() string {
	return fmt.Sprintf("head: %v, end: %v, weight: %d", g.head, g.end, g.weight)
}

type state struct {
	current grid.Point
	from    grid.Point
}

type walkState struct {
	start grid.Point
}

func direction(from, to grid.Point) grid.Vec {
	return grid.Vec{
		X: to.X - from.X,
		Y: to.Y - from.Y,
	}
}

type stateV2 struct {
	current grid.Point
	node    graphNode
	count   int
	visited *set.Set[grid.Point]
}

func (t *trails) walkV2() []int {
	// DFS
	var pos []int
	s := stack.Stack[stateV2]{
		List: list.New(),
	}
	s.Push(stateV2{
		current: t.start,
		node: graphNode{
			head: t.start,
		},
		count: 0,
		visited: set.New[grid.Point](
			t.start,
		),
	})
	for {
		if s.Len() == 0 {
			break
		}
		start := s.Pop()
		visited := start.visited
		if start.node.end == t.end {
			pos = append(pos, start.count)
			continue
		}
		c := start.count
		for _, r := range t.graph[start.current] {
			if visited.Has(r.end) {
				continue
			}
			vc := visited.Copy()
			vc.Add(r.head)
			s.Push(stateV2{
				current: r.end,
				node:    r,
				count:   c + r.weight,
				visited: vc,
			})
		}
	}
	return pos
}

type weightPath struct {
	current grid.Point
	weight  int
	visited *set.Set[grid.Point]
}

// minimize 用 dfs 最小化整个迷宫，忽略直线的路径
func (t *trails) minimize() {
	q := stack.Stack[walkState]{
		List: list.New(),
	}
	q.Push(walkState{
		start: t.start,
	})
	t.visited[t.start] = 0
	cross := set.New[grid.Point](
		t.start,
	)
	for {
		if q.Len() == 0 {
			break
		}
		start := q.Pop()
		sp := start.start
		var nbs []grid.Point
		for _, p := range t.m.GetNeighbor(sp.X, sp.Y) {
			if t.m.Get(p.X, p.Y) != '#' {
				// 因为只是遍历一次迷宫，找到所有岔路口，所以一旦走过的点不再走
				if _, ok := t.visited[grid.Point(p)]; ok {
					continue
				}
				t.visited[grid.Point(p)] = 0
				nbs = append(nbs, grid.Point(p))
			}
		}
		if len(nbs) > 1 {
			cross.Add(sp)
		}
		for _, p := range nbs {
			q.Push(walkState{
				start: grid.Point(p),
			})
		}
	}
	fmt.Printf("cross: %v\n", cross.List())
	cross.Each(func(item grid.Point) bool {
		if item == t.end {
			return true
		}
		cs := stack.Stack[weightPath]{
			List: list.New(),
		}
		cs.Push(weightPath{
			current: item,
			weight:  0,
			visited: set.New[grid.Point](
				item,
			),
		})
		for {
			if cs.Len() == 0 {
				break
			}
			start := cs.Pop()
			cur := start.current
			weight := start.weight
			visited := start.visited.Copy()
			var nbs []grid.Point
			for _, p := range t.m.GetNeighbor(cur.X, cur.Y) {
				if t.m.Get(p.X, p.Y) == '#' {
					continue
				}
				if visited.Has(grid.Point(p)) {
					continue
				}
				if cross.Has(grid.Point(p)) || grid.Point(p) == t.end || grid.Point(p) == t.start {
					node := graphNode{
						head:   item,
						weight: weight + 1,
						end:    grid.Point(p),
					}
					t.graph[item] = append(t.graph[item], node)
					continue
				}

				nbs = append(nbs, grid.Point(p))
			}
			for _, p := range nbs {
				visited.Add(grid.Point(p))
				cs.Push(weightPath{
					current: grid.Point(p),
					weight:  weight + 1,
					visited: visited,
				})
			}
		}
		return true
	})
}

func (t *trails) walk() {
	q := queue.Queue[state]{
		List: list.New(),
	}
	q.Push(state{
		current: t.start,
		from:    t.start,
	})
	t.visited[t.start] = 0
	for {
		if q.Len() == 0 {
			break
		}
		pp := q.Pop()
		// fmt.Printf("point: %v\n", pp)
		start := pp.current
		from := pp.from
		for _, point := range t.m.GetNeighbor(start.X, start.Y) {
			v := t.m.Get(point.X, point.Y)
			if v == '#' {
				continue
			}
			// 不走回头路
			if point.X == from.X && point.Y == from.Y {
				continue
			}
			if p, ok := t.visited[grid.Point(point)]; ok {
				if t.visited[start] <= p-1 {
					continue
				}
			}
			// fmt.Printf("next point: %v\n", point)
			switch v {
			// 滑行不会导致撞墙
			case '>':
				next := grid.Point{
					X: point.X,
					Y: point.Y + 1,
				}
				// 不能上滑梯
				if next.X == start.X && next.Y == start.Y {
					continue
				}
				t.visited[grid.Point(point)] = t.visited[start] + 1
				q.Push(state{
					current: next,
					from:    start,
				})
				t.visited[next] = t.visited[start] + 2
			case '<':
				next := grid.Point{
					X: point.X,
					Y: point.Y - 1,
				}
				if next.X == start.X && next.Y == start.Y {
					continue
				}
				t.visited[grid.Point(point)] = t.visited[start] + 1
				q.Push(state{
					current: next,
					from:    start,
				})
				t.visited[next] = t.visited[start] + 2
			case '^':
				t.visited[grid.Point(point)] = t.visited[start] + 1
				next := grid.Point{
					X: point.X - 1,
					Y: point.Y,
				}
				if next.X == start.X && next.Y == start.Y {
					continue
				}
				q.Push(state{
					current: next,
					from:    start,
				})
				t.visited[next] = t.visited[start] + 2
			case 'v':
				t.visited[grid.Point(point)] = t.visited[start] + 1
				next := grid.Point{
					X: point.X + 1,
					Y: point.Y,
				}
				if next.X == start.X && next.Y == start.Y {
					continue
				}
				q.Push(state{
					current: next,
					from:    start,
				})
				t.visited[next] = t.visited[start] + 2
			default:
				q.Push(state{
					current: grid.Point(point),
					from:    start,
				})
				t.visited[grid.Point(point)] = t.visited[start] + 1
			}
		}
	}
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[byte](line, col)
	start := grid.Point{}
	end := grid.Point{}
	txt.ReadByLineEx(context.Background(), func(i int, l string) error {
		for j, v := range l {
			if i == 0 && v == '.' {
				start = grid.Point{X: i, Y: j}
			}
			if i == line-1 && v == '.' {
				end = grid.Point{X: i, Y: j}
			}
			m.Add(i, j, byte(v))
		}
		return nil
	})
	ts := &trails{start: start, m: m, end: end, visited: make(map[grid.Point]int)}
	// ts.m.PrintEx("%c")
	ts.walk()
	// fmt.Printf("visited: %v\n", ts.visited)
	distance := ts.visited[end]
	fmt.Printf("p1: %d\n", distance)
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[byte](line, col)
	start := grid.Point{}
	end := grid.Point{}
	txt.ReadByLineEx(context.Background(), func(i int, l string) error {
		for j, v := range l {
			if i == 0 && v == '.' {
				start = grid.Point{X: i, Y: j}
			}
			if i == line-1 && v == '.' {
				end = grid.Point{X: i, Y: j}
			}
			m.Add(i, j, byte(v))
		}
		return nil
	})
	ts := &trails{
		start:   start,
		m:       m,
		end:     end,
		visited: make(map[grid.Point]int),
		graph:   make(map[grid.Point][]graphNode),
	}
	ts.minimize()
	// fmt.Printf("graph: %v\n", ts.graph)
	pos := ts.walkV2()
	sort.Slice(pos, func(i, j int) bool {
		return pos[i] < pos[j]
	})
	fmt.Printf("p2: %d\n", pos[len(pos)-1])
}

func main() {
	// p1()
	p2()
}
