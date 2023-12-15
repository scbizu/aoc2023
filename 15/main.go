package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
)

type sequence struct {
	s     []byte
	value int
}

func (s sequence) String() string {
	return string(s.s)
}

func (s sequence) HASH() int {
	for _, v := range s.s {
		s.value += int(v)
		s.value *= 17
		s.value %= 256
	}
	return s.value
}

func p1() {
	txt := input.NewTXTFile("./input.txt")
	var seqs []sequence
	txt.ReadByBlock(context.Background(), ",", func(block []string) error {
		for _, v := range block {
			v = strings.TrimSuffix(v, "\n")
			seqs = append(seqs, sequence{s: []byte(v)})
		}
		return nil
	})
	var total int
	for _, v := range seqs {
		total += v.HASH()
	}
	fmt.Fprintf(os.Stdout, "p1: %d\n", total)
}

type box struct {
	lens map[string]Len
	cap  int
}

type Len struct {
	label    string
	focalLen int
	index    int
}

func (b *box) remove(l Len) {
	delete(b.lens, l.label)
}

func (b *box) add(l Len) {
	if _, ok := b.lens[l.label]; ok {
		b.lens[l.label] = Len{
			label:    l.label,
			focalLen: l.focalLen,
			index:    b.lens[l.label].index,
		}
		b.cap++
		return
	}
	b.lens[l.label] = Len{
		label:    l.label,
		focalLen: l.focalLen,
		index:    b.cap + 1,
	}
	b.cap++
}

func (b *box) String() string {
	var lens []Len
	for _, l := range b.lens {
		lens = append(lens, l)
	}
	sort.Slice(lens, func(i, j int) bool {
		return lens[i].index < lens[j].index
	})
	var s []string
	for _, l := range lens {
		s = append(s, fmt.Sprintf("%s:%d:%d", l.label, l.focalLen, l.index))
	}
	return strings.Join(s, ",")
}

func p2() {
	txt := input.NewTXTFile("./input.txt")
	var strs []string
	txt.ReadByBlock(context.Background(), ",", func(block []string) error {
		for _, v := range block {
			v = strings.TrimSuffix(v, "\n")
			strs = append(strs, v)
		}
		return nil
	})
	boxes := make(map[int]*box)
	for _, s := range strs {
		switch {
		case strings.Contains(s, "="):
			parts := strings.Split(s, "=")
			seq := sequence{s: []byte(parts[0])}
			if _, ok := boxes[seq.HASH()]; !ok {
				boxes[seq.HASH()] = &box{
					lens: make(map[string]Len),
				}
				boxes[seq.HASH()].add(Len{
					label:    parts[0],
					focalLen: mustAtoi(parts[1]),
				})
			} else {
				boxes[seq.HASH()].add(Len{
					label:    parts[0],
					focalLen: mustAtoi(parts[1]),
				})
			}
		case strings.Contains(s, "-"):
			s := strings.TrimSuffix(s, "-")
			seq := sequence{s: []byte(s)}
			if _, ok := boxes[seq.HASH()]; ok {
				boxes[seq.HASH()].remove(Len{
					label: s,
				})
				if len(boxes[seq.HASH()].lens) == 0 {
					delete(boxes, seq.HASH())
				}
			}
		default:
			panic("invalid input")
		}
	}
	var total int
	for boxIndex, v := range boxes {
		var lens []Len
		for _, l := range v.lens {
			lens = append(lens, l)
		}
		sort.Slice(lens, func(i, j int) bool {
			return lens[i].index < lens[j].index
		})
		for lenIndex, l := range lens {
			value := (boxIndex + 1) * (lenIndex + 1) * l.focalLen
			// fmt.Fprintf(os.Stdout, "%s:%d:%d:%d\n", l.label, (boxIndex + 1), (lenIndex + 1), l.focalLen)
			total += value
		}
	}
	// fmt.Printf("%+v\n", boxes)
	fmt.Fprintf(os.Stdout, "p2: %d\n", total)
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func main() {
	p1()
	p2()
}
