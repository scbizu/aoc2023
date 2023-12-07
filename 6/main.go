package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
)

func main() {
	txt := input.NewTXTFile("input.txt")
	var raceTimes []int64
	var dists []int64
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		// Time line
		if i == 0 {
			line = strings.TrimPrefix("Time: ", line)
			tms := strings.Split(line, " ")
			for _, tm := range tms {
				tm = removeBlank(tm)
				if tm == "" {
					continue
				}
				raceTimes = append(raceTimes, mustInt64(tm))
			}
		}
		if i == 1 {
			line = strings.TrimPrefix("Distance: ", line)
			dists := strings.Split(line, " ")
			for _, dist := range dists {
				dist = removeBlank(dist)
				if dist == "" {
					continue
				}
				dists = append(dists, mustInt64(dist))
			}
		}
		return nil
	})
}

func removeBlank(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func mustInt64(s string) int64 {
	i64, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return int64(i64)
}
