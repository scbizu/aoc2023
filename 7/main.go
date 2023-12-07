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

type kind uint8

const (
	kindHighCard kind = iota + 1
	kindOnePair
	kindTwoPair
	kindThreeOfAKind
	kindFullHouse
	kindFourOfAKind
	kindFiveOfAKind
)

type hand struct {
	cards []byte
	bid   int64
}

func (h hand) String() string {
	return fmt.Sprintf("%s\n", string(h.cards))
}

func (h hand) Kind() kind {
	var counts []int64
	for _, c := range h.cards {
		counts = append(counts, int64(strings.Count(string(h.cards), string(c))))
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i] > counts[j]
	})
	switch {
	case counts[0] == 5:
		return kindFiveOfAKind
	case counts[0] == 4:
		return kindFourOfAKind
	case counts[0] == 3 && counts[len(counts)-1] == 2:
		return kindFullHouse
	case counts[0] == 3 && counts[len(counts)-1] == 1:
		return kindThreeOfAKind
	case counts[0] == 2 && counts[len(counts)-2] == 2 && counts[len(counts)-1] == 1:
		return kindTwoPair
	case counts[0] == 2 && counts[len(counts)-2] == 1 && counts[len(counts)-1] == 1:
		return kindOnePair
	default:
		return kindHighCard
	}
}

func (h *hand) KindWithJoker() kind {
	var counts []int64
	var jokers int64
	for _, c := range h.cards {
		if c == 'J' {
			jokers++
			continue
		}
		counts = append(counts, int64(strings.Count(string(h.cards), string(c))))
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i] > counts[j]
	})
	if len(counts) == 0 && int(jokers) == len(h.cards) {
		return kindFiveOfAKind
	}
	counts[0] += jokers
	switch {
	case counts[0] == 5:
		return kindFiveOfAKind
	case counts[0] == 4:
		return kindFourOfAKind
	case counts[0] == 3 && counts[len(counts)-1] == 2:
		return kindFullHouse
	case counts[0] == 3 && counts[len(counts)-1] == 1:
		return kindThreeOfAKind
	case counts[0] == 2 && counts[len(counts)-2] == 2 && counts[len(counts)-1] == 1:
		return kindTwoPair
	case counts[0] == 2 && counts[len(counts)-2] == 1 && counts[len(counts)-1] == 1:
		return kindOnePair
	default:
		return kindHighCard
	}
}

func pt(c byte) int64 {
	switch c {
	case 'A':
		return 14
	case 'K':
		return 13
	case 'Q':
		return 12
	case 'J':
		return 11
	case 'T':
		return 10
	default:
		return int64(c - '0')
	}
}

func jpt(c byte) int64 {
	switch c {
	case 'A':
		return 14
	case 'K':
		return 13
	case 'Q':
		return 12
	case 'J':
		return 0
	case 'T':
		return 10
	default:
		return int64(c - '0')
	}
}

func fallbackCompare(h1, h2 []byte) bool {
	for idx := range h1 {
		if pt(h1[idx]) > pt(h2[idx]) {
			return true
		}
		if pt(h1[idx]) < pt(h2[idx]) {
			return false
		}
	}
	panic("hands are equal")
}

func fallbackCompareWithJoker(h1, h2 []byte) bool {
	for idx := range h1 {
		if jpt(h1[idx]) > jpt(h2[idx]) {
			return true
		}
		if jpt(h1[idx]) < jpt(h2[idx]) {
			return false
		}
	}
	panic("hands are equal")
}

func main() {
	p1()
	p2()
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	var hands []hand
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		if line == "" {
			return nil
		}
		parts := strings.Split(line, " ")
		hands = append(hands, hand{
			cards: []byte(parts[0]),
			bid:   int64(mustInt64(parts[1])),
		})
		return nil
	})
	sort.Slice(hands, func(i, j int) bool {
		if hands[i].Kind() > hands[j].Kind() {
			return true
		}
		if hands[i].Kind() < hands[j].Kind() {
			return false
		}
		return fallbackCompare(hands[i].cards, hands[j].cards)
	})
	var total int64
	for idx, h := range hands {
		total += h.bid * (int64(len(hands)) - int64(idx))
	}
	fmt.Fprintf(os.Stdout, "p1: %d\n", total)
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	var hands []hand
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		if line == "" {
			return nil
		}
		parts := strings.Split(line, " ")
		hands = append(hands, hand{
			cards: []byte(parts[0]),
			bid:   int64(mustInt64(parts[1])),
		})
		return nil
	})
	sort.Slice(hands, func(i, j int) bool {
		if hands[i].KindWithJoker() > hands[j].KindWithJoker() {
			return true
		}
		if hands[i].KindWithJoker() < hands[j].KindWithJoker() {
			return false
		}
		return fallbackCompareWithJoker(hands[i].cards, hands[j].cards)
	})
	var total int64
	for idx, h := range hands {
		total += h.bid * (int64(len(hands)) - int64(idx))
	}
	fmt.Fprintf(os.Stdout, "p2: %d\n", total)
}

func mustInt64(s string) int64 {
	i64, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return int64(i64)
}
