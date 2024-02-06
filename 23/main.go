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
	"github.com/magejiCoder/magejiAoc/set"
)

const (
	line = 23
	col  = 23
)

type trails struct {
	start   grid.Vec
	m       matrix.Matrix[byte]
	end     grid.Vec
	visited map[grid.Vec]int
}

type state struct {
	current grid.Vec
	from    grid.Vec
}

type stateV2 struct {
	start  grid.Vec
	traces *set.Set[grid.Vec]
}

func direction(from, to grid.Vec) grid.Vec {
	return grid.Vec{
		X: to.X - from.X,
		Y: to.Y - from.Y,
	}
}

func (t *trails) walkV2() []int {
	q := queue.Queue[stateV2]{
		List: list.New(),
	}
	q.Push(stateV2{
		start:  t.start,
		traces: set.New[grid.Vec](),
	})
	var poss []int
	for {
		if q.Len() == 0 {
			break
		}
		pp := q.Pop()
		start := pp.start
		trace := pp.traces
		if start == t.end {
			fmt.Printf("trace: %d\n", trace.Size())
			poss = append(poss, trace.Size())
			continue
		}
		if p, ok := t.visited[start]; ok && p > trace.Size() {
			fmt.Printf("p: %d, trace: %d\n", p, trace.Size())
			continue
		} else {
			if ok {
				fmt.Printf("p: %d, trace: %d\n", p, trace.Size())
			}
			t.visited[start] = trace.Size()
		}
		for _, point := range t.m.GetNeighbor(start.X, start.Y) {
			// 不允许走回头路
			if trace.Has(grid.Vec(point)) {
				continue
			}
			v := t.m.Get(point.X, point.Y)
			if v == '#' {
				continue
			}
			// fmt.Printf("next point: %v\n", point)
			switch v {
			// 滑行不会导致撞墙
			case '>', '<', '^', 'v':
				// 在行进方向继续滑行,不需要考虑上下滑梯
				next := grid.Vec{
					X: point.X + direction(start, grid.Vec(point)).X,
					Y: point.Y + direction(start, grid.Vec(point)).Y,
				}
				if trace.Has(next) {
					continue
				}
				newTrace := trace.Copy()
				newTrace.Add(start)
				newTrace.Add(grid.Vec(point))
				q.Push(stateV2{
					start:  next,
					traces: newTrace,
				})
			default:
				nt := trace.Copy()
				nt.Add(start)
				q.Push(stateV2{
					start:  grid.Vec(point),
					traces: nt,
				})
			}
		}
	}
	return poss
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
			if p, ok := t.visited[grid.Vec(point)]; ok {
				if t.visited[start] <= p-1 {
					continue
				}
			}
			// fmt.Printf("next point: %v\n", point)
			switch v {
			// 滑行不会导致撞墙
			case '>':
				next := grid.Vec{
					X: point.X,
					Y: point.Y + 1,
				}
				// 不能上滑梯
				if next.X == start.X && next.Y == start.Y {
					continue
				}
				t.visited[grid.Vec(point)] = t.visited[start] + 1
				q.Push(state{
					current: next,
					from:    start,
				})
				t.visited[next] = t.visited[start] + 2
			case '<':
				next := grid.Vec{
					X: point.X,
					Y: point.Y - 1,
				}
				if next.X == start.X && next.Y == start.Y {
					continue
				}
				t.visited[grid.Vec(point)] = t.visited[start] + 1
				q.Push(state{
					current: next,
					from:    start,
				})
				t.visited[next] = t.visited[start] + 2
			case '^':
				t.visited[grid.Vec(point)] = t.visited[start] + 1
				next := grid.Vec{
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
				t.visited[grid.Vec(point)] = t.visited[start] + 1
				next := grid.Vec{
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
					current: grid.Vec(point),
					from:    start,
				})
				t.visited[grid.Vec(point)] = t.visited[start] + 1
			}
		}
	}
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[byte](line, col)
	start := grid.Vec{}
	end := grid.Vec{}
	txt.ReadByLineEx(context.Background(), func(i int, l string) error {
		for j, v := range l {
			if i == 0 && v == '.' {
				start = grid.Vec{X: i, Y: j}
			}
			if i == line-1 && v == '.' {
				end = grid.Vec{X: i, Y: j}
			}
			m.Add(i, j, byte(v))
		}
		return nil
	})
	ts := &trails{start: start, m: m, end: end, visited: make(map[grid.Vec]int)}
	// ts.m.PrintEx("%c")
	ts.walk()
	// fmt.Printf("visited: %v\n", ts.visited)
	distance := ts.visited[end]
	fmt.Printf("p1: %d\n", distance)
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[byte](line, col)
	start := grid.Vec{}
	end := grid.Vec{}
	txt.ReadByLineEx(context.Background(), func(i int, l string) error {
		for j, v := range l {
			if i == 0 && v == '.' {
				start = grid.Vec{X: i, Y: j}
			}
			if i == line-1 && v == '.' {
				end = grid.Vec{X: i, Y: j}
			}
			m.Add(i, j, byte(v))
		}
		return nil
	})
	ts := &trails{start: start, m: m, end: end, visited: make(map[grid.Vec]int)}
	// ts.m.PrintEx("%c")
	poss := ts.walkV2()
	sort.Slice(poss, func(i, j int) bool {
		return poss[i] < poss[j]
	})
	fmt.Printf("p2: %d\n", poss[len(poss)-1])
	// fmt.Printf("visited: %v\n", ts.visited)
	// fmt.Printf("p2: %d\n", ts.visited[end])
}

func main() {
	// p1()
	p2()
}
