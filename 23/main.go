package main

import (
	"context"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/matrix"
)

const (
	line = 23
	col  = 23
)

type trails struct {
	start matrix.Point[byte]
	m     matrix.Matrix[byte]
	end   matrix.Point[byte]
}

func (t trails) walk() {
}

func main() {
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[byte](line, col)
	start := matrix.Point[byte]{}
	end := matrix.Point[byte]{}
	txt.ReadByLineEx(context.Background(), func(i int, l string) error {
		for j, v := range l {
			if i == 0 && v == '.' {
				start = matrix.Point[byte]{X: i, Y: j, Value: byte(v)}
			}
			if i == line-1 && v == '.' {
				end = matrix.Point[byte]{X: i, Y: j, Value: byte(v)}
			}
			m.Add(i, j, byte(v))
		}
		return nil
	})
	ts := trails{start: start, m: m, end: end}
	ts.m.PrintEx("%c")
}
