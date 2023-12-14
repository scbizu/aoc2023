package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
	"github.com/scbizu/aoc2022/helper/matrix"
)

type direction uint8

const (
	vertical direction = iota + 1
	horizontal
)

type mirrorMatrix struct {
	m matrix.Matrix[byte]
}

func (m mirrorMatrix) smudge(row, col int) mirrorMatrix {
	nm := matrix.NewMatrix[byte](m.m.Rows(), m.m.Cols())
	for i := 0; i < m.m.Rows(); i++ {
		for j := 0; j < m.m.Cols(); j++ {
			if i == row && j == col {
				if m.m[i][j] == '.' {
					nm.Add(i, j, '#')
				}
				if m.m[i][j] == '#' {
					nm.Add(i, j, '.')
				}
			} else {
				nm.Add(i, j, m.m[i][j])
			}
		}
	}
	return mirrorMatrix{m: nm}
}

func (m mirrorMatrix) getSmudge(lastDL directionLine) directionLine {
	for i := 0; i < m.m.Rows(); i++ {
		for j := 0; j < m.m.Cols(); j++ {
			newM := m.smudge(i, j)
			newM.m.PrintEx("%c")
			dls := newM.mirror(func(mg *mirrorGetter) {
				mg.isGetAll = true
			})
			for _, dl := range dls {
				if !dl.isValid() {
					continue
				}
				if dl.String() == lastDL.String() {
					continue
				}
				return dl
			}
		}
	}
	panic("no smudge")
}

type directionLine struct {
	d direction
	l int
}

func (dl directionLine) isValid() bool {
	return dl.d > 0
}

func (dl directionLine) String() string {
	switch dl.d {
	case vertical:
		return fmt.Sprintf("|: %d", dl.l)
	case horizontal:
		return fmt.Sprintf("-: %d", dl.l)
	default:
		panic("invalid direction")
	}
}

func (m mirrorMatrix) Row(i int) []byte {
	return m.m[i]
}

func (m mirrorMatrix) Col(i int) []byte {
	var cols []byte
	for j := 0; j < m.m.Rows(); j++ {
		cols = append(cols, m.m[j][i])
	}
	return cols
}

func (m mirrorMatrix) isMirror(d directionLine) bool {
	switch d.d {
	case vertical:
		var i int
		for {
			fmt.Printf("peers: %d,%d\n", d.l-1-i, d.l+2+i)
			if d.l-1-i < 0 || d.l+2+i >= m.m.Cols() {
				break
			}
			if string(m.Col(d.l-1-i)) != string(m.Col(d.l+2+i)) {
				return false
			}
			i++
		}
		fmt.Printf("isMirror: %s\n", d)
		return true
	case horizontal:
		var j int
		for {
			fmt.Printf("peers: %d,%d\n", d.l-1-j, d.l+2+j)
			if d.l-1-j < 0 || d.l+2+j >= m.m.Rows() {
				break
			}
			if string(m.Row(d.l-1-j)) != string(m.Row(d.l+2+j)) {
				return false
			}
			j++
		}
		fmt.Printf("isMirror: %s\n", d)
		return true
	default:
		panic("invalid direction")
	}
}

type mirrorGetter struct {
	isGetAll bool
}

type mirrorOption func(*mirrorGetter)

func (m mirrorMatrix) mirror(opts ...mirrorOption) []directionLine {
	if m.m.Rows() == 0 || m.m.Cols() == 0 {
		panic("invalid matrix")
	}

	var mg mirrorGetter
	for _, opt := range opts {
		opt(&mg)
	}

	var dis []directionLine

	i, j := 0, 1

	for {
		if j >= m.m.Rows() {
			break
		}
		fmt.Printf("line: (%d,%d)\n", i, j)
		if string(m.Row(i)) == string(m.Row(j)) {
			if m.isMirror(directionLine{d: horizontal, l: i}) {
				dis = append(dis, directionLine{d: horizontal, l: i})
			}
		}
		i++
		j++
	}

	i, j = 0, 1

	for {
		if j >= m.m.Cols() {
			break
		}
		fmt.Printf("col: (%d,%d)\n", i, j)
		if string(m.Col(i)) == string(m.Col(j)) {
			if m.isMirror(directionLine{d: vertical, l: i}) {
				dis = append(dis, directionLine{d: vertical, l: i})
			}
		}
		i++
		j++
	}

	if mg.isGetAll {
		return dis
	}

	return []directionLine{dis[0]}
}

func main() {
	// p1()
	p2()
}

func p2() {
	txt := input.NewTXTFile("./input.txt")

	var total int

	txt.ReadByBlock(context.Background(), "\n\n", func(block []string) error {
		if len(block) == 0 {
			return nil
		}
		for _, mt := range block {
			lines := strings.Split(mt, "\n")
			if lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
			m := matrix.NewMatrix[byte](len(lines), len(lines[0]))
			for i, line := range lines {
				for j, v := range line {
					m.Add(i, j, byte(v))
				}
			}
			mm := mirrorMatrix{m: m}
			dl := mm.getSmudge(mm.mirror()[0])
			fmt.Printf("mirror: %+v\n", dl)
			var score int
			switch dl.d {
			case vertical:
				score += dl.l + 1
			case horizontal:
				score += 100 * (dl.l + 1)
			}
			fmt.Printf("score: %d\n", score)
			total += score
		}
		return nil
	})
	fmt.Fprintf(os.Stdout, "p2: %d\n", total)
}

func p1() {
	txt := input.NewTXTFile("./input.txt")

	var total int

	txt.ReadByBlock(context.Background(), "\n\n", func(block []string) error {
		if len(block) == 0 {
			return nil
		}
		for _, mt := range block {
			lines := strings.Split(mt, "\n")
			if lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
			m := matrix.NewMatrix[byte](len(lines), len(lines[0]))
			for i, line := range lines {
				for j, v := range line {
					m.Add(i, j, byte(v))
				}
			}
			m.PrintEx("%c")
			mm := mirrorMatrix{m: m}
			dl := mm.mirror()[0]
			fmt.Printf("mirror: %+v\n", dl)
			var score int
			switch dl.d {
			case vertical:
				score += dl.l + 1
			case horizontal:
				score += 100 * (dl.l + 1)
			}
			fmt.Printf("score: %d\n", score)
			total += score
		}
		return nil
	})
	fmt.Fprintf(os.Stdout, "p1: %d\n", total)
}
