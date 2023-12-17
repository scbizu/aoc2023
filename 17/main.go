package main

import (
	"context"
	"fmt"
	"os"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/matrix"
)

const (
	defaultLimit = 3
)

type direction uint8

const (
	up direction = iota + 1
	down
	left
	right
)

func (d direction) String() string {
	switch d {
	case up:
		return "👆"
	case down:
		return "👇"
	case left:
		return "<-"
	case right:
		return "->"
	}
	panic("invalid direction")
}

var caches = make(map[matrix.Point[int]]int)

type vec struct {
	row, col int
	d        direction
}

type heatLossMap struct {
	m matrix.Matrix[int]
}

func getDirection(p1, p2 matrix.Point[int]) direction {
	if p1.X == p2.X {
		if p1.Y > p2.Y {
			return left
		} else {
			return right
		}
	} else {
		if p1.X > p2.X {
			return up
		} else {
			return down
		}
	}
}

type directionLimit map[direction]int

func (h *heatLossMap) fastPath(
	p matrix.Point[int],
	value int,
	dl directionLimit,
) {
	if p.X == h.m.Rows()-1 && p.Y == h.m.Cols()-1 {
		return
	}
	vecs := h.m.GetNeighbor(p.X, p.Y)
	for _, v := range vecs {
		if limit, ok := dl[getDirection(p, matrix.Point[int]{X: v.X, Y: v.Y})]; ok && limit == 0 {
			fmt.Printf("skip at %d,%d\n", v.X, v.Y)
			continue
		}
		old, ok := caches[matrix.Point[int]{X: v.X, Y: v.Y}]
		if !ok {
			old := h.m.Get(v.X, v.Y)
			caches[matrix.Point[int]{X: v.X, Y: v.Y}] = value + old
		} else {
			if value+h.m.Get(v.X, v.Y) > old {
				continue
			} else {
				caches[matrix.Point[int]{X: v.X, Y: v.Y}] = value + h.m.Get(v.X, v.Y)
			}
		}
		for d := range dl {
			if d == getDirection(p, matrix.Point[int]{X: v.X, Y: v.Y}) {
				dl[direction(d)] -= 1
			} else {
				dl[direction(d)] = defaultLimit
			}
		}
		// fmt.Printf("go to %d,%d,dl: %v\n", v.X, v.Y, dl)
		h.fastPath(matrix.Point[int]{X: v.X, Y: v.Y}, value+h.m.Get(v.X, v.Y), dl)
	}
}

const (
	row = 13
	col = 13
)

func p1() {
	caches = make(map[matrix.Point[int]]int)
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[int](row, col)
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		for j, v := range line {
			m.Add(i, j, int(v-'0'))
		}
		return nil
	})
	m.PrintEx("%d")
	hl := &heatLossMap{
		m: m,
	}
	caches[matrix.Point[int]{X: 0, Y: 0}] = m.Get(0, 0)
	hl.fastPath(matrix.Point[int]{X: 0, Y: 0}, m.Get(0, 0), map[direction]int{
		up:    defaultLimit,
		down:  defaultLimit,
		left:  defaultLimit,
		right: defaultLimit,
	})
	fmt.Printf("%+v\n", caches)
	fmt.Fprintf(os.Stdout, "p1: %d\n", caches[matrix.Point[int]{X: row - 1, Y: col - 1}])
}

func main() {
	p1()
}
