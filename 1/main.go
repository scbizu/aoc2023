package main

import (
	"bytes"
	"context"
	"fmt"
	"index/suffixarray"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
)

func main() {
	phase1()
	phase2()
}

func phase1() {
	txt := input.NewTXTFile("./input.txt")
	ctx := context.Background()
	var total int
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		if line == "" {
			return nil
		}
		number := bytes.NewBuffer(nil)
		for _, c := range line {
			if '0' < c && c <= '9' {
				number.WriteByte(byte(c))
			}
		}
		numberStr := bytes.NewBuffer(nil)
		for i := 0; i < number.Len(); i++ {
			if i == 0 {
				numberStr.WriteByte(number.Bytes()[i])
				continue
			}
			if i == number.Len()-1 {
				numberStr.WriteByte(number.Bytes()[i])
				continue
			}
		}
		if numberStr.Len() == 1 {
			numberStr.WriteByte(numberStr.Bytes()[0])
		}
		n, _ := strconv.Atoi(numberStr.String())
		total += n
		return nil
	})
	fmt.Fprintf(os.Stdout, "phase 1: %d\n", total)
}

var digitMap = map[string]byte{
	"one":   '1',
	"two":   '2',
	"three": '3',
	"four":  '4',
	"five":  '5',
	"six":   '6',
	"seven": '7',
	"eight": '8',
	"nine":  '9',
}

type digit struct {
	index int
	len   int
	value byte
}

func phase2() {
	txt := input.NewTXTFile("./input.txt")
	ctx := context.Background()
	var total int
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		if line == "" {
			return nil
		}
		var indexes []digit
		for k, v := range digitMap {
			if !strings.Contains(line, k) {
				continue
			}
			index := suffixarray.New([]byte(line))
			res := index.Lookup([]byte(k), -1)
			for _, r := range res {
				indexes = append(indexes, digit{
					index: r,
					len:   len(k),
					value: v,
				})
			}
		}
		for index, c := range line {
			if '0' < c && c <= '9' {
				indexes = append(indexes, digit{
					index: index,
					len:   1,
					value: byte(c),
				})
			}
		}
		sort.Slice(indexes, func(i, j int) bool {
			return indexes[i].index < indexes[j].index
		})
		number := bytes.NewBuffer(nil)
		for i := 0; i < len(indexes); i++ {
			number.WriteByte(indexes[i].value)
		}
		numberStr := bytes.NewBuffer(nil)
		for i := 0; i < number.Len(); i++ {
			if i == 0 {
				numberStr.WriteByte(number.Bytes()[i])
				continue
			}
			if i == number.Len()-1 {
				numberStr.WriteByte(number.Bytes()[i])
				continue
			}
		}
		if numberStr.Len() == 1 {
			numberStr.WriteByte(numberStr.Bytes()[0])
		}
		n, err := strconv.Atoi(numberStr.String())
		if err != nil {
			panic(err)
		}
		fmt.Printf("numberStr: %s, n: %d\n", number.String(), n)
		total += n
		return nil
	})
	fmt.Fprintf(os.Stdout, "phase 2: %d\n", total)
}
