package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
)

type race struct {
	time   int64
	record int64
}

func quickPath(n int64, record int64) []int64 {
	var records []int64
	x := (n - n/2)
	y := (n / 2)
	for {
		mul := x * y
		// fmt.Printf("x: %d, y: %d,mul: %d\n", x, y, mul)
		if mul <= record {
			break
		}
		if x == y {
			records = append(records, mul)
		} else {
			records = append(records, mul, mul)
		}
		x += 1
		y -= 1
	}
	return records
}

func (r *race) winRecords() []int64 {
	return quickPath(r.time, r.record)
}

func main() {
	txt := input.NewTXTFile("input.txt")
	var raceTimes []int64
	var raceDistances []int64
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		// Time line
		if i == 0 {
			line = strings.TrimPrefix(line, "Time: ")
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
			line = strings.TrimPrefix(line, "Distance: ")
			dists := strings.Split(line, " ")
			for _, dist := range dists {
				dist = removeBlank(dist)
				if dist == "" {
					continue
				}
				raceDistances = append(raceDistances, mustInt64(dist))
			}
		}
		return nil
	})
	multiple := 1
	for index, raceTime := range raceTimes {
		r := &race{
			time:   raceTime,
			record: raceDistances[index],
		}
		records := r.winRecords()
		println(len(records))
		if len(records) > 0 {
			multiple *= len(records)
		}
	}
	fmt.Fprintf(os.Stdout, "p1: %d\n", multiple)

	raceTimes = []int64{}
	raceDistances = []int64{}
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		// Time line
		if i == 0 {
			line = strings.TrimPrefix(line, "Time: ")
			line = removeBlank(line)
			raceTimes = append(raceTimes, mustInt64(line))
		}
		if i == 1 {
			line = strings.TrimPrefix(line, "Distance: ")
			line = removeBlank(line)
			println(line)
			raceDistances = append(raceDistances, mustInt64(line))
		}
		return nil
	})
	r := &race{
		time:   raceTimes[0],
		record: raceDistances[0],
	}
	records := r.winRecords()
	fmt.Fprintf(os.Stdout, "p2: %d\n", len(records))
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
