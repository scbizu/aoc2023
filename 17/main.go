package main

import (
	"context"
	"fmt"
	"os"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/matrix"
)

const (
	defaultLimit = 2
)

type direction uint8

const (
	up direction = iota + 1
	down
	left
	right
)

var caches = make(map[matrix.Point[int]]int)

type vec struct {
	matrix.Point[int] // point
	d                 directionLimit
}

type heatLossMap struct {
	m matrix.Matrix[int]
}

func (h heatLossMap) Print() {
	for i := 0; i < h.m.Rows(); i++ {
		for j := 0; j < h.m.Cols(); j++ {
			fmt.Printf("%d ", caches[matrix.Point[int]{
				X:     i,
				Y:     j,
				Value: h.m.Get(i, j),
			}])
		}
		fmt.Println()
	}
}

func (h heatLossMap) neighbor(x, y int) []matrix.Point[int] {
	var ns []matrix.Point[int]
	if y < h.m.Cols()-1 {
		ns = append(ns, matrix.Point[int]{X: x, Y: y + 1, Value: h.m.Get(x, y+1)})
	}
	if x < h.m.Rows()-1 {
		ns = append(ns, matrix.Point[int]{X: x + 1, Y: y, Value: h.m.Get(x+1, y)})
	}
	if x > 0 {
		ns = append(ns, matrix.Point[int]{X: x - 1, Y: y, Value: h.m.Get(x-1, y)})
	}
	if y > 0 {
		ns = append(ns, matrix.Point[int]{X: x, Y: y - 1, Value: h.m.Get(x, y-1)})
	}
	return ns
}

type directionValue struct {
	d direction
	v int
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
	points []vec,
) {
	if len(points) == 0 {
		return
	}
	p := points[0]
	dl := p.d
	value := caches[matrix.Point[int]{
		X:     p.X,
		Y:     p.Y,
		Value: p.Value,
	}]
	if p.X == h.m.Rows()-1 && p.Y == h.m.Cols()-1 {
		return
	}
	vecs := h.neighbor(p.X, p.Y)
	for _, v := range vecs {
		fmt.Printf("(%d,%d)[%d] -> (%d,%d)[%d]\n", p.X, p.Y, value, v.X, v.Y, v.Value)
		newdl := make(directionLimit)
		if limit, ok := dl[getDirection(p.Point, matrix.Point[int]{X: v.X, Y: v.Y})]; ok && limit == 0 {
			newdl[getDirection(p.Point, matrix.Point[int]{X: v.X, Y: v.Y})] = defaultLimit + 1
			continue
		}
		for d := range dl {
			if d == getDirection(p.Point, matrix.Point[int]{X: v.X, Y: v.Y}) {
				newdl[direction(d)] = dl[direction(d)] - 1
			} else {
				newdl[direction(d)] = defaultLimit + 1
			}
		}
		old, ok := caches[matrix.Point[int]{
			X:     v.X,
			Y:     v.Y,
			Value: v.Value,
		}]
		if !ok {
			// fmt.Printf("new: at (%d,%d),dl: %v,value: %d\n", v.X, v.Y, dl, value)
			caches[matrix.Point[int]{
				X:     v.X,
				Y:     v.Y,
				Value: v.Value,
			}] = value + v.Value
		} else {
			if value+v.Value <= old {
				// fmt.Printf("set: at (%d,%d),dl: %v,value: %d\n", v.X, v.Y, dl, value)
				caches[matrix.Point[int]{
					X:     v.X,
					Y:     v.Y,
					Value: v.Value,
				}] = value + v.Value
			} else {
				continue
			}
		}
		points = append(points, vec{
			Point: v,
			d:     newdl,
		})
	}
	h.fastPath(points[1:])
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
	caches[matrix.Point[int]{
		X:     0,
		Y:     0,
		Value: m.Get(0, 0),
	}] = m.Get(0, 0)
	v := vec{
		Point: matrix.Point[int]{
			X:     0,
			Y:     0,
			Value: m.Get(0, 0),
		},
		d: map[direction]int{
			up:    defaultLimit,
			down:  defaultLimit,
			left:  defaultLimit,
			right: defaultLimit,
		},
	}
	hl.fastPath([]vec{v})

	hl.Print()

	fmt.Fprintf(os.Stdout, "p1: %d\n", caches[matrix.Point[int]{
		X: m.Rows() - 1,
		Y: m.Cols() - 1,
		Value: m.Get(
			m.Rows()-1,
			m.Cols()-1,
		),
	}])
}

func main() {
	p1()
}
