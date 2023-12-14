package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/magejiCoder/set"
	"github.com/scbizu/aoc2022/helper/input"
)

var runCache = make(map[string]int)

type sequence struct {
	seq       []byte
	formation []int
	resSeq    []byte
	resSet    *set.Set[string]
	total     int
}

func (s sequence) String() string {
	return fmt.Sprintf("seq: %s,  formation: %+v\n", string(s.seq), s.formation)
}

func (s *sequence) count() {
	// fmt.Printf("%s\n", s)
	// if _, ok := runCache[string(s.seq)]; !ok {
	// 	runCache[string(s.seq)] = s.resSet.Size()
	// } else {
	// 	return
	// }
	if len(s.formation) == 0 {
		if len(s.seq) > 0 {
			for i := 0; i < len(s.seq); i++ {
				if s.seq[i] == '#' {
					return
				}
				s.resSeq = append(s.resSeq, '.')
			}
		}
		if !s.resSet.Has(string(s.resSeq)) {
			fmt.Printf("res: %s\n", string(s.resSeq))
		}
		s.resSet.Add(string(s.resSeq))
		s.total++
		return
	}
	if len(s.seq) == 0 {
		return
	}
	n := s.formation[0]
	switch s.seq[0] {
	case '.':
		s.seq = s.seq[1:]
		s.resSeq = append(s.resSeq, '.')
		s.count()
	case '#':
		switch {
		case len(s.seq) == n && !strings.Contains(string(s.seq[:n]), "."):
			for i := 0; i < n; i++ {
				s.resSeq = append(s.resSeq, '#')
			}
			s.seq = s.seq[n:]
			s.formation = s.formation[1:]
			s.count()
		case len(s.seq) > n && !strings.Contains(string(s.seq[:n]), ".") && s.seq[n] != '#':
			s.seq = s.seq[n+1:]
			s.formation = s.formation[1:]
			for i := 0; i < n; i++ {
				s.resSeq = append(s.resSeq, '#')
			}
			s.resSeq = append(s.resSeq, '.')
			s.count()
		default:
			return
		}
	case '?':
		s1 := make([]byte, len(s.seq))
		copy(s1, s.seq)
		s1[0] = '.'
		s1q := sequence{
			seq:       s1,
			formation: s.formation,
			resSeq:    s.resSeq,
			resSet:    s.resSet,
			total:     s.total,
		}
		s1q.count()
		// fmt.Printf("? -> #: %s\n", s.seq)
		s2 := make([]byte, len(s.seq))
		copy(s2, s.seq)
		s2[0] = '#'
		s2q := sequence{
			seq:       s2,
			formation: s.formation,
			resSeq:    s.resSeq,
			resSet:    s.resSet,
			total:     s.total,
		}
		s2q.count()
	}
}

func sum(s []int) int {
	var sum int
	for idx, v := range s {
		sum += v
		if idx != len(s)-1 {
			sum += 1
		}
	}
	return sum
}

func main() {
	p1()
	// p2()
}

// func p2() {
// 	txt := input.NewTXTFile("./input.txt")
// 	var allTotal int
// 	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
// 		if line == "" {
// 			return nil
// 		}
// 		total = 0
// 		parts := strings.Split(line, " ")
// 		if len(parts) != 2 {
// 			panic("invalid input")
// 		}
// 		numStrs := strings.Split(parts[1], ",")
// 		seq := sequence{
// 			seq:       []byte(parts[0]),
// 			formation: numStrsToNums(numStrs),
// 			resSet:    set.New[string](),
// 		}
// 		seq.count()
// 		p1Total := total
// 		total = 0
// 		var s []byte
// 		s = append(s, []byte(parts[0])...)
// 		s = append(s, '?')
// 		s = append(s, []byte(parts[0])...)
// 		seq = sequence{
// 			seq:       s,
// 			formation: double(numStrsToNums(numStrs)),
// 			resSet:    set.New[string](),
// 		}
// 		seq.count()
// 		rate := total / p1Total
// 		allTotal += total * rate * rate * rate
// 		fmt.Printf("total: %d,rate: %d, p1Total: %d,allTotal: %d\n", total, rate, p1Total, total*rate*rate*rate)
// 		return nil
// 	})
// 	fmt.Fprintf(os.Stdout, "p2: %d\n", allTotal)
// }

func p2() {
	txt := input.NewTXTFile("./input.txt")
	var total int
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		if line == "" {
			return nil
		}
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic("invalid input")
		}
		numStrs := strings.Split(parts[1], ",")
		seq := sequence{
			seq:       unfoldSeq([]byte(parts[0])),
			formation: unfoldFormation(numStrsToNums(numStrs)),
			resSet:    set.New[string](),
		}
		seq.count()
		total += seq.resSet.Size()
		return nil
	})
	fmt.Fprintf(os.Stdout, "p2: %d\n", total)
}

func unfoldSeq(s []byte) []byte {
	var res []byte
	for i := 0; i < 5; i++ {
		res = append(res, s...)
		if i != 4 {
			res = append(res, '?')
		}
	}
	return res
}

func unfoldFormation(s []int) []int {
	var res []int
	for i := 0; i < 5; i++ {
		res = append(res, s...)
	}
	return res
}

// func double[T comparable](s []T) []T {
// 	var res []T
// 	for i := 0; i < 2; i++ {
// 		res = append(res, s...)
// 	}
// 	return res
// }

func p1() {
	txt := input.NewTXTFile("./input.txt")
	var total int
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		if line == "" {
			return nil
		}
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			panic("invalid input")
		}
		numStrs := strings.Split(parts[1], ",")
		seq := sequence{
			seq:       []byte(parts[0]),
			formation: numStrsToNums(numStrs),
			resSet:    set.New[string](),
		}
		seq.count()
		total += seq.total
		return nil
	})
	fmt.Fprintf(os.Stdout, "p1: %d\n", total)
}

func numStrsToNums(numStrs []string) []int {
	var nums []int
	for _, numStr := range numStrs {
		nums = append(nums, atoi(numStr))
	}
	return nums
}

func atoi(s string) int {
	var num int
	for _, v := range []byte(s) {
		num = num*10 + int(v-'0')
	}
	return num
}
