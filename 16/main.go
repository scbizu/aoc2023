package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/magejiCoder/set"
	"github.com/scbizu/aoc2022/helper/input"
	"github.com/scbizu/aoc2022/helper/matrix"
)

type direction int

const (
	up direction = iota + 1
	down
	left
	right
)

type contraption struct {
	m          matrix.Matrix[byte]
	traces     *set.Set[vec]
	at         matrix.Point[byte]
	direction  direction
	traceCount int
}

func countTrace(traces *set.Set[vec]) int {
	pmap := make(map[matrix.Point[byte]]struct{})
	traces.Each(func(item vec) bool {
		pmap[item.point] = struct{}{}
		return true
	})
	return len(pmap)
}

type vec struct {
	point matrix.Point[byte]
	dir   direction
}

func (c *contraption) applyTrace() {
	c.traces.Each(func(item vec) bool {
		c.m[item.point.X][item.point.Y] = '#'
		return true
	})
}

func (c contraption) isPointOut(p matrix.Point[byte]) bool {
	return p.X < 0 || p.Y < 0 || p.X >= c.m.Rows() || p.Y >= c.m.Cols()
}

func (c *contraption) goThrough() {
	if c.isPointOut(c.at) {
		return
	}
	if c.traces.Has(vec{
		point: c.at,
		dir:   c.direction,
	}) {
		return
	}
	c.traces.Add(vec{
		point: c.at,
		dir:   c.direction,
	})

	switch c.direction {
	default:
		panic("unknown direction")
	case up:
		next := c.at
		switch c.m.Get(next.X, next.Y) {
		case '.':
			c.at = matrix.Point[byte]{X: next.X - 1, Y: next.Y}
			c.direction = up
			c.goThrough()
		case '|':
			if c.isPointOut(matrix.Point[byte]{X: next.X - 1, Y: next.Y}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X - 1, Y: next.Y}
			c.direction = up
			c.goThrough()
		case '/':
			if c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y + 1}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X, Y: next.Y + 1}
			c.direction = right
			c.goThrough()
		case '\\':
			if c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y - 1}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X, Y: next.Y - 1}
			c.direction = left
			c.goThrough()
		case '-':
			if !c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y - 1}) {
				c.at = matrix.Point[byte]{X: next.X, Y: next.Y - 1}
				c.direction = left
				c.goThrough()
			}
			if !c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y + 1}) {
				c.at = matrix.Point[byte]{X: next.X, Y: next.Y + 1}
				c.direction = right
				c.goThrough()
			}
		}
	case down:
		next := c.at
		switch c.m.Get(next.X, next.Y) {
		case '.':
			c.at = matrix.Point[byte]{X: next.X + 1, Y: next.Y}
			c.direction = down
			c.goThrough()
		case '|':
			if c.isPointOut(matrix.Point[byte]{X: next.X + 1, Y: next.Y}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X + 1, Y: next.Y}
			c.direction = down
			c.goThrough()
		case '/':
			if c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y - 1}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X, Y: next.Y - 1}
			c.direction = left
			c.goThrough()
		case '\\':
			if c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y + 1}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X, Y: next.Y + 1}
			c.direction = right
			c.goThrough()
		case '-':
			if !c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y - 1}) {
				c.at = matrix.Point[byte]{X: next.X, Y: next.Y - 1}
				c.direction = left
				c.goThrough()
			}
			if !c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y + 1}) {
				c.at = matrix.Point[byte]{X: next.X, Y: next.Y + 1}
				c.direction = right
				c.goThrough()
			}
		}
	case left:
		next := c.at
		switch c.m.Get(next.X, next.Y) {
		case '.':
			c.at = matrix.Point[byte]{X: next.X, Y: next.Y - 1}
			c.goThrough()
		case '|':
			if !c.isPointOut(matrix.Point[byte]{X: next.X + 1, Y: next.Y}) {
				c.at = matrix.Point[byte]{X: next.X + 1, Y: next.Y}
				c.direction = down
				c.goThrough()
			}
			if !c.isPointOut(matrix.Point[byte]{X: next.X - 1, Y: next.Y}) {
				c.at = matrix.Point[byte]{X: next.X - 1, Y: next.Y}
				c.direction = up
				c.goThrough()
			}
		case '/':
			if c.isPointOut(matrix.Point[byte]{X: next.X + 1, Y: next.Y}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X + 1, Y: next.Y}
			c.direction = down
			c.goThrough()
		case '\\':
			if c.isPointOut(matrix.Point[byte]{X: next.X - 1, Y: next.Y}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X - 1, Y: next.Y}
			c.direction = up
			c.goThrough()
		case '-':
			if c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y - 1}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X, Y: next.Y - 1}
			c.goThrough()
		}
	case right:
		next := c.at
		switch c.m.Get(next.X, next.Y) {
		case '.':
			c.at = matrix.Point[byte]{X: next.X, Y: next.Y + 1}
			c.goThrough()
		case '|':
			if !c.isPointOut(matrix.Point[byte]{X: next.X + 1, Y: next.Y}) {
				c.at = matrix.Point[byte]{X: next.X + 1, Y: next.Y}
				c.direction = down
				c.goThrough()
			}
			if !c.isPointOut(matrix.Point[byte]{X: next.X - 1, Y: next.Y}) {
				c.at = matrix.Point[byte]{X: next.X - 1, Y: next.Y}
				c.direction = up
				c.goThrough()
			}
		case '/':
			if c.isPointOut(matrix.Point[byte]{X: next.X - 1, Y: next.Y}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X - 1, Y: next.Y}
			c.direction = up
			c.goThrough()
		case '\\':
			if c.isPointOut(matrix.Point[byte]{X: next.X + 1, Y: next.Y}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X + 1, Y: next.Y}
			c.direction = down
			c.goThrough()
		case '-':
			if c.isPointOut(matrix.Point[byte]{X: next.X, Y: next.Y + 1}) {
				return
			}
			c.at = matrix.Point[byte]{X: next.X, Y: next.Y + 1}
			c.goThrough()
		}
	}
}

const (
	row, line = 110, 110
)

var runCache = make(map[vec]*set.Set[vec])

func p1() {
	runCache = make(map[vec]*set.Set[vec])
	txt := input.NewTXTFile("./input.txt")
	m := matrix.NewMatrix[byte](row, line)
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		for j, v := range line {
			m.Add(i, j, byte(v))
		}
		return nil
	})
	trace := set.New[vec]()
	c := &contraption{
		m:         m,
		traces:    trace,
		at:        matrix.Point[byte]{X: 0, Y: 0},
		direction: right,
	}
	c.goThrough()
	// for debug
	// c.applyTrace()
	// c.m.PrintEx("%c")
	fmt.Fprintf(os.Stdout, "p1: %v\n", countTrace(c.traces))
}

func p2() {
	runCache = make(map[vec]*set.Set[vec])
	txt := input.NewTXTFile("./input.txt")
	m := matrix.NewMatrix[byte](row, line)
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		for j, v := range line {
			m.Add(i, j, byte(v))
		}
		return nil
	})
	var counts []int
	m.ForEach(func(x, y int, _ byte) {
		if x != 0 && y != 0 && x != m.Rows()-1 && y != m.Cols()-1 {
			return
		}
		if x == 0 {
			trace := set.New[vec]()
			c := &contraption{
				m:         m,
				traces:    trace,
				at:        matrix.Point[byte]{X: x, Y: y},
				direction: down,
			}
			c.goThrough()
			counts = append(counts, countTrace(c.traces))
		}
		if y == 0 {
			trace := set.New[vec]()
			c := &contraption{
				m:         m,
				traces:    trace,
				at:        matrix.Point[byte]{X: x, Y: y},
				direction: right,
			}
			c.goThrough()
			counts = append(counts, countTrace(c.traces))
		}
		if x == m.Rows()-1 {
			trace := set.New[vec]()
			c := &contraption{
				m:         m,
				traces:    trace,
				at:        matrix.Point[byte]{X: x, Y: y},
				direction: up,
			}
			c.goThrough()
			counts = append(counts, countTrace(c.traces))
		}
		if y == m.Cols()-1 {
			trace := set.New[vec]()
			c := &contraption{
				m:         m,
				traces:    trace,
				at:        matrix.Point[byte]{X: x, Y: y},
				direction: left,
			}
			c.goThrough()
			counts = append(counts, countTrace(c.traces))
		}
	})
	sort.Slice(counts, func(i, j int) bool {
		return counts[i] > counts[j]
	})
	fmt.Fprintf(os.Stdout, "p2: %v\n", counts[0])
}

func main() {
	p1()
	p2()
}
