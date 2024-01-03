package main

import (
	"container/list"
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/matrix"
	"github.com/magejiCoder/magejiAoc/queue"
)

func main() {
	p1()
	// p2()
}

const (
	row = 11
	col = 11
)

const (
	maxStep = 100
)

type garden struct {
	m     matrix.Matrix[byte]
	paths queue.Queue[matrix.Point[byte]]
	all   map[matrix.Point[byte]]struct{}
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
		nbs := g.m.GetNeighbor(p.X, p.Y)
		for _, nb := range nbs {
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

func (g *garden) infiniteWalk() {
	g.load()
	for {
		if g.paths.Len() == 0 {
			break
		}
		p := g.paths.Pop()
		if p.X-1 < 0 {
			g.m = g.m.JoinTop(g.m)
			p.X = row
		}
		if p.X+1 >= g.m.Rows() {
			g.m = g.m.JoinBottom(g.m)
		}
		if p.Y-1 < 0 {
			g.m = g.m.JoinLeft(g.m)
			p.Y = col
		}
		if p.Y+1 >= g.m.Cols() {
			g.m = g.m.JoinRight(g.m)
		}
		nbs := g.m.GetNeighbor(p.X, p.Y)
		for _, nb := range nbs {
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
		g.print()
		println()
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
}
