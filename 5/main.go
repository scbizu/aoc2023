package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
)

type almanac struct {
	seeds, soils, ferts, waters, lights, temps, hums, locs []int64
}

func (a *almanac) String() string {
	return fmt.Sprintf("seeds: %v\nsoils: %v\nferts: %v\nwaters: %v\nlights: %v\ntemps: %v\nhums: %v\nlocs: %v\n", a.seeds, a.soils, a.ferts, a.waters, a.lights, a.temps, a.hums, a.locs)
}

type i64Vec struct {
	from, to int64
}

type rangeAlmanac struct {
	seeds, soils, ferts, waters, lights, temps, hums, locs []i64Vec
	almanacFn
}

func (ra *rangeAlmanac) reset() {
	ra.seeds = []i64Vec{}
	ra.soils = []i64Vec{}
	ra.ferts = []i64Vec{}
	ra.waters = []i64Vec{}
	ra.lights = []i64Vec{}
	ra.temps = []i64Vec{}
	ra.hums = []i64Vec{}
	ra.locs = []i64Vec{}
}

func (ra *rangeAlmanac) AddSeeds(ss []int64) {
	for _, sv := range buildSeedsVecs(ss) {
		ra.seeds = append(ra.seeds, splitRange(sv, 0)...)
	}
}

func (ra *rangeAlmanac) minLoc() int64 {
	return ra.locs[0].from
}

func (ra *rangeAlmanac) Do() {
	fmt.Printf("seeds: %v\n", ra.seeds)
	for _, seed := range ra.seeds {
		for i := seed.from; i <= seed.to; i++ {
			if includes(ra.soils, ra.seed2Soil(i)) {
				continue
			}
			ra.soils = merge(ra.soils, ra.seed2Soil(i))
		}
		fmt.Printf("soils: %v\n", ra.soils)
	}
	ra.soils = splitVecs(ra.soils)
	for _, soil := range ra.soils {
		for i := soil.from; i <= soil.to; i++ {
			if includes(ra.ferts, ra.soil2Fert(i)) {
				continue
			}
			ra.ferts = merge(ra.ferts, ra.soil2Fert(i))
		}
		fmt.Printf("ferts: %v\n", ra.ferts)
	}
	ra.ferts = splitVecs(ra.ferts)
	for _, fert := range ra.ferts {
		for i := fert.from; i <= fert.to; i++ {
			if includes(ra.waters, ra.fert2Water(i)) {
				continue
			}
			ra.waters = merge(ra.waters, ra.fert2Water(i))
		}
	}
	fmt.Printf("waters: %v\n", ra.waters)
	ra.waters = splitVecs(ra.waters)
	for _, water := range ra.waters {
		for i := water.from; i <= water.to; i++ {
			if includes(ra.lights, ra.water2Light(i)) {
				continue
			}
			ra.lights = merge(ra.lights, ra.water2Light(i))
		}
	}
	ra.lights = splitVecs(ra.lights)
	fmt.Printf("lights: %v\n", ra.lights)
	for _, light := range ra.lights {
		for i := light.from; i <= light.to; i++ {
			if includes(ra.temps, ra.light2Temp(i)) {
				continue
			}
			ra.temps = merge(ra.temps, ra.light2Temp(i))
		}
	}
	ra.temps = splitVecs(ra.temps)
	fmt.Printf("temps: %v\n", ra.temps)
	for _, temp := range ra.temps {
		for i := temp.from; i <= temp.to; i++ {
			if includes(ra.hums, ra.temp2Hum(i)) {
				continue
			}
			ra.hums = merge(ra.hums, ra.temp2Hum(i))
		}
	}
	ra.hums = splitVecs(ra.hums)
	fmt.Printf("hums: %v\n", ra.hums)
	for _, hum := range ra.hums {
		for i := hum.from; i <= hum.to; i++ {
			if includes(ra.locs, ra.hum2Loc(i)) {
				continue
			}
			ra.locs = merge(ra.locs, ra.hum2Loc(i))
		}
	}
	fmt.Printf("locs: %v\n", ra.locs)
}

func includes(vecs []i64Vec, target int64) bool {
	for _, vec := range vecs {
		if vec.from <= target && target <= vec.to {
			return true
		}
	}
	return false
}

func merge(vecs []i64Vec, n int64) []i64Vec {
	if len(vecs) == 0 {
		return []i64Vec{{n, n}}
	}
	var newVecs []i64Vec
	var isMerged bool
	for i := 0; i < len(vecs); i++ {
		// expand
		if n == vecs[i].from-1 {
			newVecs = append(newVecs, i64Vec{n, vecs[i].to})
			isMerged = true
			continue
		}
		// merge
		if n == vecs[i].to+1 && i != len(vecs)-1 && n == vecs[i+1].from-1 {
			newVecs = append(newVecs, i64Vec{vecs[i].from, vecs[i+1].to})
			isMerged = true
			i++
			continue
		}
		// expand
		if n == vecs[i].to+1 {
			newVecs = append(newVecs, i64Vec{vecs[i].from, n})
			isMerged = true
			continue
		}
		newVecs = append(newVecs, vecs[i])
	}
	if !isMerged {
		newVecs = append(newVecs, i64Vec{n, n})
	}
	sort.Slice(newVecs, func(i, j int) bool {
		return newVecs[i].from < newVecs[j].from
	})
	return newVecs
}

type almanacFn struct {
	seed2Soil, soil2Fert, fert2Water, water2Light, light2Temp, temp2Hum, hum2Loc func(int64) int64
}

func parseAlmanacFn(lines []string) (*almanacFn, error) {
	af := &almanacFn{}
	// soils
	sos := strings.Split(lines[1], "\n")
	soilRs := make([]rangeDirection, 0, len(sos)-1)
	for _, soilString := range sos[1:] {
		parts := strings.Split(soilString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid soil: %s", soilString)
		}
		soilRule := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		soilRs = append(soilRs, soilRule)
	}
	af.seed2Soil = func(seed int64) int64 {
		return doMapping(soilRs, seed)
	}
	// fertilizers
	ferts := strings.Split(lines[2], "\n")
	fertRules := make([]rangeDirection, 0, len(ferts)-1)
	for _, fertString := range ferts[1:] {
		parts := strings.Split(fertString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid fertilizer: %s", fertString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		fertRules = append(fertRules, rd)
	}
	af.soil2Fert = func(soil int64) int64 {
		return doMapping(fertRules, soil)
	}
	// water
	waters := strings.Split(lines[3], "\n")
	waterRules := make([]rangeDirection, 0, len(waters)-1)
	for _, waterString := range waters[1:] {

		parts := strings.Split(waterString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid water: %s", waterString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		waterRules = append(waterRules, rd)
	}

	af.fert2Water = func(fert int64) int64 {
		return doMapping(waterRules, fert)
	}

	// light
	lights := strings.Split(lines[4], "\n")
	lightRules := make([]rangeDirection, 0, len(lights)-1)
	for _, lightString := range lights[1:] {
		parts := strings.Split(lightString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid light: %s", lightString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		lightRules = append(lightRules, rd)
	}

	af.water2Light = func(water int64) int64 {
		return doMapping(lightRules, water)
	}

	// temperature
	temperatures := strings.Split(lines[5], "\n")
	tempRules := make([]rangeDirection, 0, len(temperatures)-1)
	for _, temperatureString := range temperatures[1:] {
		parts := strings.Split(temperatureString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid temperature: %s", temperatureString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		tempRules = append(tempRules, rd)
	}

	af.light2Temp = func(light int64) int64 {
		return doMapping(tempRules, light)
	}

	// humidity
	humidities := strings.Split(lines[6], "\n")
	humRules := make([]rangeDirection, 0, len(humidities)-1)
	for _, humidityString := range humidities[1:] {
		parts := strings.Split(humidityString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid humidity: %s", humidityString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		humRules = append(humRules, rd)
	}

	af.temp2Hum = func(temp int64) int64 {
		return doMapping(humRules, temp)
	}

	// location
	locations := strings.Split(lines[7], "\n")
	locRules := make([]rangeDirection, 0, len(locations)-1)
	for _, locationString := range locations[1:] {
		if locationString == "" {
			continue
		}
		parts := strings.Split(locationString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid location: %s", locationString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		locRules = append(locRules, rd)
	}

	af.hum2Loc = func(hum int64) int64 {
		return doMapping(locRules, hum)
	}
	return af, nil
}

func parseAlmanac(lines []string) (*almanac, error) {
	a := &almanac{}
	if len(lines) != 8 {
		return nil, fmt.Errorf("invalid almanac")
	}
	// seeds
	seedsString := strings.TrimPrefix(lines[0], "seeds: ")
	ss, err := splitString2Int64(seedsString, " ")
	if err != nil {
		return nil, fmt.Errorf("failed to parse seeds: %w", err)
	}
	a.seeds = ss
	// soils
	sos := strings.Split(lines[1], "\n")
	rds := make([]rangeDirection, 0, len(sos)-1)
	for _, soilString := range sos[1:] {
		parts := strings.Split(soilString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid soil: %s", soilString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}

	for _, seed := range a.seeds {
		a.soils = append(a.soils, doMapping(rds, seed))
	}
	// fertilizers
	ferts := strings.Split(lines[2], "\n")
	rds = make([]rangeDirection, 0, len(ferts)-1)
	for _, fertString := range ferts[1:] {
		parts := strings.Split(fertString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid fertilizer: %s", fertString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, soil := range a.soils {
		a.ferts = append(a.ferts, doMapping(rds, soil))
	}
	// water
	waters := strings.Split(lines[3], "\n")
	rds = make([]rangeDirection, 0, len(waters)-1)
	for _, waterString := range waters[1:] {
		parts := strings.Split(waterString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid water: %s", waterString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, fert := range a.ferts {
		a.waters = append(a.waters, doMapping(rds, fert))
	}
	// light
	lights := strings.Split(lines[4], "\n")
	rds = make([]rangeDirection, 0, len(lights)-1)
	for _, lightString := range lights[1:] {
		parts := strings.Split(lightString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid light: %s", lightString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, water := range a.waters {
		a.lights = append(a.lights, doMapping(rds, water))
	}
	// temperature
	temperatures := strings.Split(lines[5], "\n")
	rds = make([]rangeDirection, 0, len(temperatures)-1)
	for _, temperatureString := range temperatures[1:] {
		parts := strings.Split(temperatureString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid temperature: %s", temperatureString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, light := range a.lights {
		a.temps = append(a.temps, doMapping(rds, light))
	}
	// humidity
	humidities := strings.Split(lines[6], "\n")
	rds = make([]rangeDirection, 0, len(humidities)-1)
	for _, humidityString := range humidities[1:] {
		parts := strings.Split(humidityString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid humidity: %s", humidityString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, temperature := range a.temps {
		a.hums = append(a.hums, doMapping(rds, temperature))
	}
	// location
	locations := strings.Split(lines[7], "\n")
	rds = make([]rangeDirection, 0, len(locations)-1)
	for _, locationString := range locations[1:] {
		if locationString == "" {
			continue
		}
		parts := strings.Split(locationString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid location: %s", locationString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, humidity := range a.hums {
		a.locs = append(a.locs, doMapping(rds, humidity))
	}

	return a, nil
}

type rangeDirection struct {
	source int64
	target int64
	rngLen int64
}

func (rd rangeDirection) mapTo(source int64) int64 {
	if source > rd.source+rd.rngLen || source < rd.source {
		return source
	}
	return rd.target + (source - rd.source)
}

func doMapping(rds []rangeDirection, source int64) int64 {
	for _, rd := range rds {
		target := rd.mapTo(source)
		if target != source {
			return target
		}
	}
	return source
}

func main() {
	txt := input.NewTXTFile("input.txt")
	txt.ReadByBlock(context.Background(), "\n\n", func(lines []string) error {
		a, err := parseAlmanac(lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse almanac: %v\n", err)
			return nil
		}
		minLoc := math.MaxInt64
		for _, loc := range a.locs {
			if loc < int64(minLoc) {
				minLoc = int(loc)
			}
		}
		fmt.Fprintf(os.Stdout, "p1: %d\n", minLoc)
		return nil
	})

	txt.ReadByBlock(context.Background(), "\n\n", func(lines []string) error {
		af, err := parseAlmanacFn(lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse almanac: %v\n", err)
			return nil
		}
		ra := &rangeAlmanac{
			almanacFn: *af,
		}

		seedsString := strings.TrimPrefix(lines[0], "seeds: ")
		ss, err := splitString2Int64(seedsString, " ")
		if err != nil {
			panic(err)
		}
		var minLoc int64 = math.MaxInt64
		for i := 0; i < len(ss); i += 2 {
			// TODO(scnace): use goroutine and waitGroup to speed up
			start := ss[i]
			length := ss[i+1]
			ra.AddSeeds([]int64{start, start + length - 1})
			ra.Do()
			if ra.minLoc() < minLoc {
				minLoc = ra.minLoc()
			}
		}
		fmt.Fprintf(os.Stdout, "p2: %d\n", ra.minLoc())
		return nil
	})
}

func splitString2Int64(raw string, sep string) ([]int64, error) {
	raws := strings.Split(raw, sep)
	ints := make([]int64, len(raws))
	for i, raw := range raws {
		if _, err := fmt.Sscanf(raw, "%d", &ints[i]); err != nil {
			return nil, fmt.Errorf("failed to parse string to int64: %w", err)
		}
	}
	return ints, nil
}

func atoi64(raw string) int64 {
	var i int64
	if _, err := fmt.Sscanf(raw, "%d", &i); err != nil {
		panic(err)
	}
	return i
}

const splitSize = 8

func buildSeedsVecs(ss []int64) []i64Vec {
	var seeds []i64Vec
	for i := 0; i < len(ss); i += 2 {
		start := ss[i]
		length := ss[i+1]
		seeds = append(seeds, i64Vec{start, start + length - 1})
	}
	return seeds
}

func splitVecs(vecs []i64Vec) []i64Vec {
	var seedVecs []i64Vec
	for _, vec := range vecs {
		seedVecs = append(seedVecs, splitRange(vec, 0)...)
	}
	return seedVecs
}

func splitRange(vec i64Vec, count uint32) []i64Vec {
	var seedVecs []i64Vec
	if count < splitSize {
		seedVecs = append(seedVecs, splitRange(i64Vec{
			from: vec.from,
			to:   vec.from + (vec.to-vec.from)/2,
		}, count+1)...)
		seedVecs = append(seedVecs, splitRange(i64Vec{
			from: vec.from + (vec.to-vec.from)/2,
			to:   vec.to,
		}, count+1)...)
	} else {
		seedVecs = append(seedVecs, vec)
	}
	return seedVecs
}
