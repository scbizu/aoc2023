package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/magejiCoder/set"
	"github.com/scbizu/aoc2022/helper/input"
)

type sequence struct {
	query []int64
}

func (s sequence) expand(isReverse bool) sequence {
	// fmt.Printf("seq: %v\n", s.query)
	if !diff(s.query) {
		return sequence{
			query: []int64{s.query[1]},
		}
	}
	var newQuery []int64

	if isReverse {
		newQuery = append(newQuery, s.query[0]-s.query[1])
		for i := 2; i < len(s.query); i++ {
			subAbs := s.query[i-1] - s.query[i]
			newQuery = append(newQuery, subAbs)
		}
	} else {
		newQuery = append(newQuery, s.query[1]-s.query[0])
		for i := 2; i < len(s.query); i++ {
			subAbs := s.query[i] - s.query[i-1]
			newQuery = append(newQuery, subAbs)
		}
	}
	if diff(newQuery) {
		ss := sequence{
			query: newQuery,
		}
		next := ss.expand(isReverse).query
		newQuery = append(newQuery, next[len(next)-1])
	}
	if isReverse {
		s.query = append(s.query, s.query[len(s.query)-1]-newQuery[len(newQuery)-1])
	} else {
		s.query = append(s.query, s.query[len(s.query)-1]+newQuery[len(newQuery)-1])
	}
	// fmt.Printf("new seq: %v\n", s.query)
	return s
}

func diff(query []int64) bool {
	dset := set.New[int64]()
	for _, q := range query {
		dset.Add(q)
	}
	return dset.Size() > 1
}

func main() {
	p1()
	p2()
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	var total int64
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		if line == "" {
			return nil
		}
		nums := strings.Split(line, " ")
		seq := sequence{
			query: make([]int64, 0, len(nums)),
		}
		for i := len(nums) - 1; i >= 0; i-- {
			seq.query = append(seq.query, MustAtoi(nums[i]))
		}
		seq = seq.expand(true)
		// fmt.Printf("result seq: %v\n", seq.query)
		total += seq.query[len(seq.query)-1]
		return nil
	})
	fmt.Printf("p2: %d\n", total)
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	var total int64
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		if line == "" {
			return nil
		}
		nums := strings.Split(line, " ")
		seq := sequence{
			query: make([]int64, 0, len(nums)),
		}
		for _, num := range nums {
			seq.query = append(seq.query, MustAtoi(num))
		}
		seq = seq.expand(false)
		fmt.Printf("result seq: %v\n", seq.query)
		total += seq.query[len(seq.query)-1]
		return nil
	})
	fmt.Printf("p1: %d\n", total)
}

func MustAtoi(s string) int64 {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return int64(n)
}
