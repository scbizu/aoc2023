package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/scbizu/aoc2022/helper/input"
	"github.com/scbizu/aoc2022/helper/matrix"
)

var tiltCache = make(map[string]matrix.Matrix[byte])

type panel struct {
	m matrix.Matrix[byte]
}

func (p panel) String() string {
	panelBytes := bytes.NewBuffer(nil)
	for i := 0; i < p.m.Rows(); i++ {
		for j := 0; j < p.m.Cols(); j++ {
			panelBytes.WriteByte(p.m[i][j])
		}
		panelBytes.WriteByte('\n')
	}
	return panelBytes.String()
}

func writeTiltCache(panel panel, next direction) {
	k := fmt.Sprintf("%s%d", panel.String(), next)
	tiltCache[k] = panel.m
}

type direction uint8

const (
	north direction = iota + 1
	west
	south
	east
)

func (p *panel) tilt(d direction) {
	if v, ok := tiltCache[fmt.Sprintf("%s%d", p.String(), d)]; ok {
		p.m = v
		return
	}
	defer writeTiltCache(*p, d)
	switch d {
	case north:
		for j := 0; j < p.m.Cols(); j++ {
			for i := 1; i < p.m.Rows(); i++ {
				if p.m[i][j] != 'O' {
					continue
				}
				up := 1
				for {
					if i-up < 0 {
						break
					}
					if p.m[i-up][j] == '#' || p.m[i-up][j] == 'O' {
						break
					}
					// forward
					p.m[i-up][j] = 'O'
					// trace
					p.m[i-up+1][j] = '.'
					up++
				}
			}
		}
	case west:
		for i := 0; i < p.m.Rows(); i++ {
			for j := 1; j < p.m.Cols(); j++ {
				if p.m[i][j] != 'O' {
					continue
				}
				left := 1
				for {
					if j-left < 0 {
						break
					}
					if p.m[i][j-left] == '#' || p.m[i][j-left] == 'O' {
						break
					}
					// forward
					p.m[i][j-left] = 'O'
					// trace
					p.m[i][j-left+1] = '.'
					left++
				}
			}
		}
	case south:
		for j := 0; j < p.m.Cols(); j++ {
			for i := p.m.Rows() - 2; i >= 0; i-- {
				if p.m[i][j] != 'O' {
					continue
				}
				down := 1
				for {
					if i+down >= p.m.Rows() {
						break
					}
					if p.m[i+down][j] == '#' || p.m[i+down][j] == 'O' {
						break
					}
					// forward
					p.m[i+down][j] = 'O'
					// trace
					p.m[i+down-1][j] = '.'
					down++
				}
			}
		}
	case east:
		for i := 0; i < p.m.Rows(); i++ {
			for j := p.m.Cols() - 2; j >= 0; j-- {
				if p.m[i][j] != 'O' {
					continue
				}
				right := 1
				for {
					if j+right >= p.m.Cols() {
						break
					}
					if p.m[i][j+right] == '#' || p.m[i][j+right] == 'O' {
						break
					}
					// forward
					p.m[i][j+right] = 'O'
					// trace
					p.m[i][j+right-1] = '.'
					right++
				}
			}
		}
	}
}

func (p panel) count() int {
	var c int
	p.m.ForEach(func(x, _ int, v byte) {
		if v == 'O' {
			c += (p.m.Rows()) - x
		}
	})
	return c
}

const (
	row, col = 100, 100
)

const (
	cycles = 1000000000
)

type circleSeq struct {
	seq []int
}

// findCircle uses **Floyd's tortoise and hare** algorithm (fixed to avoid the one-step trap).
func (c *circleSeq) findCircle() (int, int) {
	if len(c.seq) < 2 {
		return -1, 0
	}
	x0, x1 := -1, -1
	for i := 0; i < len(c.seq); i++ {
		for j := i + 1; j < len(c.seq); j++ {
			if c.seq[i] == c.seq[j] {
				if j < len(c.seq)-2 && c.seq[i+2] == c.seq[j+2] {
					x0, x1 = i, j
					break
				}
			}
		}
	}
	return x0, x1 - x0
}

func p2() {
	tiltCache = make(map[string]matrix.Matrix[byte])
	txt := input.NewTXTFile("./input.txt")
	m := matrix.NewMatrix[byte](row, col)
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		for j, c := range line {
			m.Add(i, j, byte(c))
		}
		return nil
	})
	p := &panel{m: m}
	dirs := []direction{north, west, south, east}
	cs := &circleSeq{}
	for i := 0; i < cycles; i++ {
		for _, d := range dirs {
			p.tilt(d)
		}
		cs.seq = append(cs.seq, p.count())
		// fmt.Printf("%d:%d\n", i, p.count())
		start, gap := cs.findCircle()
		if start != -1 {
			// fmt.Printf("start:%d,gap:%d\n", start, gap)
			fmt.Fprintf(os.Stdout, "p2: %d\n", cs.seq[start+(cycles-start)%gap-1])
			return
		}
	}
	fmt.Fprintf(os.Stdout, "p2: %d\n", p.count())
}

func p1() {
	tiltCache = make(map[string]matrix.Matrix[byte])
	txt := input.NewTXTFile("./input.txt")
	m := matrix.NewMatrix[byte](row, col)
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		for j, c := range line {
			m.Add(i, j, byte(c))
		}
		return nil
	})
	p := &panel{m: m}
	p.tilt(north)
	fmt.Fprintf(os.Stdout, "p1: %d\n", p.count())
}

func main() {
	p1()
	p2()
}
