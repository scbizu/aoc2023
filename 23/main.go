package main

import (
	"container/list"
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/matrix"
	"github.com/magejiCoder/magejiAoc/queue"
)

const (
	line = 141
	col  = 141
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

func (t *trails) walkV2() {
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

func main() {
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
