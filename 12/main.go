package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
)

var total int

type sequence struct {
	seq       []byte
	formation []int
}

func (s sequence) String() string {
	return fmt.Sprintf("seq: %s,  formation: %+v", string(s.seq), s.formation)
}

func (s sequence) count() sequence {
	fmt.Printf("%s\n", s)
	if len(s.formation) == 0 {
		total++
		return s
	}
	if len(s.seq) == 0 {
		return s
	}
	switch s.seq[0] {
	case '.':
		newSeq := sequence{
			seq:       s.seq[1:],
			formation: s.formation,
		}
		newSeq = newSeq.count()
		return newSeq
	case '#':
		n := s.formation[0]
		if n == 1 {
			switch {
			case len(s.seq) > 1 && s.seq[1] != '#':
				newSeq := sequence{
					seq:       s.seq[2:],
					formation: s.formation[1:],
				}
				newSeq = newSeq.count()
				return newSeq
			case len(s.seq) == 1:
				newSeq := sequence{
					seq:       s.seq[1:],
					formation: s.formation[1:],
				}
				newSeq = newSeq.count()
				return newSeq
			}
		} else {
			newSeq := sequence{
				seq: s.seq[1:],
			}
			newSeq.formation = append(newSeq.formation, n-1)
			if len(s.formation) > 1 {
				newSeq.formation = append(newSeq.formation, s.formation[1:]...)
			}
			newSeq = newSeq.count()
			return newSeq
		}
	case '?':
		tmpSeq := make([]byte, len(s.seq))
		copy(tmpSeq, s.seq)
		tmpSeq[0] = '.'
		newSeq := sequence{
			seq:       tmpSeq,
			formation: s.formation,
		}
		newSeq.count()
		tmpSeq = make([]byte, len(s.seq))
		copy(tmpSeq, s.seq)
		tmpSeq[0] = '#'
		newSeq = sequence{
			seq:       tmpSeq,
			formation: s.formation,
		}
		newSeq.count()
	}
	return s
	// var i int
	// for {
	// 	n := s.formation[0]
	// 	if n > len(s.seq) {
	// 		break
	// 	}
	// 	if i == n {
	// 		if len(s.seq) > 0 {
	// 			s.seq = s.seq[1:]
	// 		}
	// 		if len(s.formation) > 0 {
	// 			s.formation = s.formation[1:]
	// 		}
	// 		newSeq := s.count()
	// 		// fmt.Printf("newSeq: %s\n", newSeq)
	// 		s.arranges += newSeq.arranges
	// 		return s
	// 	}
	// 	seq := s.seq[0]
	// 	if seq == '.' {
	// 		s.seq = s.seq[1:]
	// 		continue
	// 	}
	// 	if seq == '?' && len(s.seq) > sum(s.formation) {
	// 		dot := sequence{
	// 			seq:       s.seq[1:],
	// 			arranges:  s.arranges,
	// 			formation: s.formation,
	// 		}
	// 		s.arranges += dot.count().arranges
	// 		sharp := sequence{
	// 			seq:       s.seq[1:],
	// 			arranges: s.arranges,
	// 			formation: s.formation[1:],
	// 		}
	// 		s.arranges += sharp.count().arranges
	// 	}
	// 	if seq == '#'  {
	// 		i++
	// 		s.seq = s.seq[1:]
	// 	}
	// }
	return s
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
	txt := input.NewTXTFile("./input.txt")
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
		}
		seq.count()
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
