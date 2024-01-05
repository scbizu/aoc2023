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

func main() {
	p1()
	p2()
}

const (
	row = 11
	col = 11
)

const (
	maxStep = 1000
)

type gardens struct {
	gq      map[grid.Vec]*garden
	count   map[gardenPoint]struct{}
	remains int64
}

type gardenPoint struct {
	gardenIndex grid.Vec
	point       matrix.Point[byte]
}

func (gs *gardens) walk2() {
	for {
		if gs.remains == 0 {
			break
		}
		gs.count = make(map[gardenPoint]struct{})
		gs.remains--
		// fmt.Printf("gardens: %d\n", len(gs.gq))
		next := map[grid.Vec][]matrix.Point[byte]{}
		for v, g := range gs.gq {
			// fmt.Printf("at garden: %+v\n", v)
			g.walk()
			// g.print()
			// println()
			// fmt.Printf("count: %d\n", g.count())
			// self
			for p := range g.all {
				gs.count[gardenPoint{
					gardenIndex: v,
					point:       p,
				}] = struct{}{}
			}
			// others
			for _, p := range g.borders {
				if _, ok := gs.gq[v.Add(p.offset)]; ok {
					if _, ok := gs.gq[v.Add(p.offset)].all[p.start]; !ok && gs.gq[v.Add(p.offset)].m.Get(p.start.X, p.start.Y) != '#' {
						// fmt.Printf("offset: %+v,start: %+v\n", v.Add(p.offset), p.start)
						next[v.Add(p.offset)] = append(next[v.Add(p.offset)], p.start)
						gs.count[gardenPoint{
							gardenIndex: v.Add(p.offset),
							point:       p.start,
						}] = struct{}{}
					}
					continue
				}
				// fmt.Printf("new garden: %+v\n", p)
				newG := &garden{
					m: g.m,
					paths: queue.Queue[matrix.Point[byte]]{
						List: list.New(),
					},
					all: make(map[matrix.Point[byte]]struct{}),
				}
				newG.all[p.start] = struct{}{}
				// newG.print()
				// fmt.Printf("count: %d\n", newG.count())
				// gs.gq[v.Add(p.offset)] = newG
				next[v.Add(p.offset)] = append(next[v.Add(p.offset)], p.start)
				gs.count[gardenPoint{
					gardenIndex: v.Add(p.offset),
					point:       p.start,
				}] = struct{}{}
			}
			g.borders = []border{}
		}
		for v, ps := range next {
			for _, p := range ps {
				if _, ok := gs.gq[v]; !ok {
					gs.gq[v] = &garden{
						m:     gs.gq[grid.Vec{X: 0, Y: 0}].m,
						paths: queue.Queue[matrix.Point[byte]]{List: list.New()},
						all:   make(map[matrix.Point[byte]]struct{}),
					}
				}
				gs.gq[v].all[p] = struct{}{}
			}
		}
	}
}

type border struct {
	start  matrix.Point[byte]
	offset grid.Vec
}

type garden struct {
	m       matrix.Matrix[byte]
	paths   queue.Queue[matrix.Point[byte]]
	all     map[matrix.Point[byte]]struct{}
	borders []border
}

func (g *garden) load() {
	for p := range g.all {
		g.paths.Push(p)
	}
	g.all = make(map[matrix.Point[byte]]struct{})
}

func (g *garden) print() {
	pMaze := matrix.NewMatrix[byte](row, col)
	g.m.ForEach(func(i, j int, v byte) {
		pMaze.Add(i, j, v)
	})
	pMaze.ForEach(func(i, j int, _ byte) {
		if _, ok := g.all[matrix.Point[byte]{
			X: i,
			Y: j,
		}]; ok {
			pMaze.Set(i, j, 'O')
		}
	})
	pMaze.PrintEx("%c")
}

func (g *garden) count() int {
	return len(g.all)
}

func (g *garden) walk() {
	g.load()
	for {
		if g.paths.Len() == 0 {
			break
		}
		p := g.paths.Pop()
		nbs := g.m.GetNeighbors(p.X, p.Y, true)
		for _, nb := range nbs {
			if nb.X < 0 {
				b := border{
					start: matrix.Point[byte]{
						X: row - 1,
						Y: nb.Y,
					},
					offset: grid.Vec{
						X: -1,
						Y: 0,
					},
				}
				g.borders = append(g.borders, b)
			}
			if nb.X > row-1 {
				b := border{
					start: matrix.Point[byte]{
						X: 0,
						Y: nb.Y,
					},
					offset: grid.Vec{
						X: 1,
						Y: 0,
					},
				}
				g.borders = append(g.borders, b)
			}
			if nb.Y < 0 {
				b := border{
					start: matrix.Point[byte]{
						X: nb.X,
						Y: col - 1,
					},
					offset: grid.Vec{
						X: 0,
						Y: -1,
					},
				}
				g.borders = append(g.borders, b)
			}
			if nb.Y > col-1 {
				b := border{
					start: matrix.Point[byte]{
						X: nb.X,
						Y: 0,
					},
					offset: grid.Vec{
						X: 0,
						Y: 1,
					},
				}
				g.borders = append(g.borders, b)
			}
			if nb.X < 0 || nb.X > row-1 || nb.Y < 0 || nb.Y > col-1 {
				continue
			}
			if g.m.Get(nb.X, nb.Y) == '#' {
				continue
			}
			if _, ok := g.all[matrix.Point[byte]{
				X: nb.X,
				Y: nb.Y,
			}]; ok {
				continue
			}
			g.all[matrix.Point[byte]{
				X: nb.X,
				Y: nb.Y,
			}] = struct{}{}
		}
	}
}

func (g *garden) TopCalibration() {
	newPoints := make(map[matrix.Point[byte]]struct{})
	for p := range g.all {
		newPoints[matrix.Point[byte]{
			X: p.X + row,
			Y: p.Y,
		}] = struct{}{}
	}
}

func (g *garden) LeftCalibration() {
	newPoints := make(map[matrix.Point[byte]]struct{})
	for p := range g.all {
		newPoints[matrix.Point[byte]{
			X: p.X,
			Y: p.Y + col,
		}] = struct{}{}
	}
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[byte](row, col)
	start := matrix.Point[byte]{}
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		for j, c := range line {
			if c == 'S' {
				start.X = i
				start.Y = j
				start.Value = byte(c)
			}
			m.Add(i, j, byte(c))
		}
		return nil
	})
	g := &garden{
		m: m,
		paths: queue.Queue[matrix.Point[byte]]{
			List: list.New(),
		},
		all: make(map[matrix.Point[byte]]struct{}),
	}
	g.all[start] = struct{}{}
	for i := 0; i < maxStep; i++ {
		g.walk()
		// g.print()
		// println()
	}
	fmt.Printf("p1: %d\n", g.count())
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[byte](row, col)
	start := matrix.Point[byte]{}
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		for j, c := range line {
			if c == 'S' {
				start.X = i
				start.Y = j
				start.Value = byte(c)
			}
			m.Add(i, j, byte(c))
		}
		return nil
	})
	g := &garden{
		m: m,
		paths: queue.Queue[matrix.Point[byte]]{
			List: list.New(),
		},
		all: make(map[matrix.Point[byte]]struct{}),
	}
	g.all[start] = struct{}{}
	gardens := &gardens{
		gq:      map[grid.Vec]*garden{{X: 0, Y: 0}: g},
		remains: maxStep,
	}
	gardens.walk2()
	fmt.Printf("p2: %d\n", len(gardens.count))
}
