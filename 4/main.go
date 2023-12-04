package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/magejiCoder/set"
	"github.com/scbizu/aoc2022/helper/input"
)

var cardInstance map[int32]int = make(map[int32]int)

type Card struct {
	index      uint32
	myNumbers  *set.Set[int32]
	winNumbers *set.Set[int32]
}

func ParseCard(raw string) Card {
	newMN := set.New[int32]()
	newWN := set.New[int32]()
	c := Card{
		myNumbers:  newMN,
		winNumbers: newWN,
	}
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		panic("invalid card")
	}
	idStr := strings.TrimPrefix(parts[0], "Card ")
	id, _ := strconv.Atoi(removeBlank(idStr))
	c.index = uint32(id)
	numberParts := strings.Split(parts[1], "|")
	if len(numberParts) != 2 {
		panic("invalid card")
	}
	myNumbers := strings.Split(numberParts[0], " ")
	for _, mn := range myNumbers {
		if mn == "" {
			continue
		}
		n, _ := strconv.Atoi(mn)
		c.myNumbers.Add(int32(n))
	}
	winNumbers := strings.Split(numberParts[1], " ")
	for _, wn := range winNumbers {
		if wn == "" {
			continue
		}
		n, _ := strconv.Atoi(wn)
		c.winNumbers.Add(int32(n))
	}

	return c
}

func removeBlank(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func (c Card) Point() int64 {
	iSet := set.Intersection[int32](c.myNumbers, c.winNumbers)
	return power2(int32(iSet.Size()) - 1)
}

func (c Card) SetInstanceMap() {
	// fmt.Printf("before: cardInstance: %+v\n", cardInstance)
	iSet := set.Intersection[int32](c.myNumbers, c.winNumbers)
	cardInstance[int32(c.index)] += 1
	count := cardInstance[int32(c.index)]
	// fmt.Printf("card %d: current count: %d, will set count: %d\n", c.index, count, iSet.Size())
	for i := 0; i < iSet.Size(); i++ {
		for j := 0; j < count; j++ {
			cardInstance[int32(c.index)+int32(i)+1] += 1
		}
	}
	// fmt.Printf("cardInstance: %+v\n", cardInstance)
}

func power2(n int32) int64 {
	if n < 0 {
		return 0
	}
	return 1 << n
}

func main() {
	txt := input.NewTXTFile("input.txt")
	var total int64
	txt.ReadByLineEx(
		context.Background(),
		func(_ int, line string) error {
			if line == "" {
				return nil
			}
			c := ParseCard(line)
			total += c.Point()
			c.SetInstanceMap()
			return nil
		})
	var cardTotal int64
	for _, v := range cardInstance {
		cardTotal += int64(v)
	}
	fmt.Fprintf(os.Stdout, "p1: total: %d\n", total)
	fmt.Fprintf(os.Stdout, "p2: total: %d\n", cardTotal)
}
