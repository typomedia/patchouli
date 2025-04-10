package structs

import (
	"sort"
	"strings"
	"time"
)

type lessFunc func(s1, s2 *System) bool

type multiSorter struct {
	systems []System
	less    []lessFunc
}

func (ms *multiSorter) Len() int      { return len(ms.systems) }
func (ms *multiSorter) Swap(i, j int) { ms.systems[i], ms.systems[j] = ms.systems[j], ms.systems[i] }
func (ms *multiSorter) Sort(systems []System) {
	ms.systems = systems
	sort.Sort(ms)
}
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.systems[i], &ms.systems[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

type Systems []System

type System struct {
	Id           string
	Name         string
	LTS          bool
	EOL          string
	MachineCount int
}

func (s System) IsEOL() bool {
	eolTime, err := time.Parse(time.DateOnly, s.EOL)
	if err != nil {
		return false
	}
	return eolTime.UTC().Before(time.Now().UTC())
}

func (s Systems) OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (s Systems) ByName() lessFunc {
	return func(s1, s2 *System) bool {
		aSlice := strings.Split(s1.Name, " ")
		bSlice := strings.Split(s2.Name, " ")
		return aSlice[0] < bSlice[0]
	}
}

func (s Systems) ByEOL() lessFunc {
	return func(s1, s2 *System) bool {
		time1, _ := time.Parse(time.DateOnly, s1.EOL)
		time2, _ := time.Parse(time.DateOnly, s2.EOL)
		return time1.UTC().After(time2.UTC())
	}
}
