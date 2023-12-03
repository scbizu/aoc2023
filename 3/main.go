package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/magejiCoder/set"
	"github.com/scbizu/aoc2022/helper/grid"
	"github.com/scbizu/aoc2022/helper/input"
	"github.com/scbizu/aoc2022/helper/matrix"
)

type vec struct {
	from matrix.Point[byte]
	to   matrix.Point[byte]
}

func main() {
	part1()
	part2()
}

func part2() {
	txt := input.NewTXTFile("input.txt")
	ctx := context.Background()
	m := matrix.NewMatrix[byte](140, 140)
	numbers := bytes.NewBuffer(nil)
	numberSet := set.New[vec]()
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		if line == "" {
			return nil
		}
		for j, v := range line {
			if isNumber(byte(v)) {
				numbers.WriteByte(byte(v))
				if j == len(line)-1 {
					numberSet.Add(vec{
						to: matrix.Point[byte]{
							X: i,
							Y: j + 1,
						},
						from: matrix.Point[byte]{
							X: i,
							Y: j - numbers.Len() + 1,
						},
					})
					numbers.Reset()
				}
			} else {
				if numbers.Len() > 0 {
					numberSet.Add(vec{
						to: matrix.Point[byte]{
							X: i,
							Y: j,
						},
						from: matrix.Point[byte]{
							X: i,
							Y: j - numbers.Len(),
						},
					})
					numbers.Reset()
				}
			}
			m.Add(i, j, byte(v))
		}
		return nil
	})
	var total int
	numberSet.Each(func(item vec) bool {
	LOOP_VEC:
		for i := item.from.Y; i < item.to.Y; i++ {
			for _, p := range m.GetNeighbors8(item.from.X, i) {
				if isSymbol(m.Get(p.X, p.Y)) {
					if isGear(m.Get(p.X, p.Y)) {
						gns := gearNumbers(m, numberSet, m.GetNeighbors8(p.X, p.Y))
						if len(gns) == 2 {
							total += gns[0] * gns[1]
						}
					}
					break LOOP_VEC
				}
			}
		}
		return true
	})
	fmt.Fprintf(os.Stdout, "p2: total: %d\n", total)
}

func part1() {
	txt := input.NewTXTFile("input.txt")
	ctx := context.Background()
	m := matrix.NewMatrix[byte](140, 140)
	numbers := bytes.NewBuffer(nil)
	var numberVecs []vec
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		if line == "" {
			return nil
		}
		for j, v := range line {
			if isNumber(byte(v)) {
				numbers.WriteByte(byte(v))
				if j == len(line)-1 {
					numberVecs = append(numberVecs, vec{
						to: matrix.Point[byte]{
							X: i,
							Y: j + 1,
						},
						from: matrix.Point[byte]{
							X: i,
							Y: j - numbers.Len() + 1,
						},
					})
					numbers.Reset()
				}
			} else {
				if numbers.Len() > 0 {
					numberVecs = append(numberVecs, vec{
						to: matrix.Point[byte]{
							X: i,
							Y: j,
						},
						from: matrix.Point[byte]{
							X: i,
							Y: j - numbers.Len(),
						},
					})
					numbers.Reset()
				}
			}
			m.Add(i, j, byte(v))
		}
		return nil
	})

	var total int
	for _, v := range numberVecs {
	LOOP_VEC:
		for i := v.from.Y; i < v.to.Y; i++ {
			for _, p := range m.GetNeighbors8(v.from.X, i) {
				if isSymbol(m.Get(p.X, p.Y)) {
					total += getVecNumber(m, v)
					break LOOP_VEC
				}
			}
		}
	}

	fmt.Fprintf(os.Stdout, "p1: total: %d\n", total)
}

func gearNumbers(
	m matrix.Matrix[byte],
	nbSet *set.Set[vec],
	pts []grid.Vec,
) []int {
	var ret []int
	for _, p := range pts {
		ns := findAndRemoveNumber(nbSet, p)
		for _, v := range ns {
			ret = append(ret, getVecNumber(m, v))
		}
	}
	return ret
}

func findAndRemoveNumber(nbSet *set.Set[vec], pt grid.Vec) []vec {
	numbers := make([]vec, 0, nbSet.Size())
	nbSet.Each(func(item vec) bool {
		if pointInVec(item, pt) {
			nbSet.Remove(item)
			numbers = append(numbers, item)
		}
		return true
	})
	return numbers
}

func pointInVec(v vec, p grid.Vec) bool {
	return p.X == v.from.X &&
		p.Y >= v.from.Y && p.Y < v.to.Y
}

func isNumber(v byte) bool {
	return v >= '0' && v <= '9'
}

func isSymbol(v byte) bool {
	return !isNumber(v) && v != '.'
}

func isGear(v byte) bool {
	return v == '*'
}

func getVecNumber(m matrix.Matrix[byte], v vec) int {
	bs := bytes.NewBuffer(nil)
	for i := v.from.Y; i < v.to.Y; i++ {
		bs.WriteByte(m.Get(v.from.X, i))
	}
	n, err := strconv.Atoi(bs.String())
	if err != nil {
		panic(err)
	}
	return n
}
