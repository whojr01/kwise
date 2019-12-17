package lib

import (
	log "kwiselog"
	str "strings"
	"testing"
)

type tdata struct {
	s string
	i int
}

var (
	tlist = []tdata{
		{"Testing 1", 1},
		{"Testing 2", 2},
		{"Testing 3", 3},
		{"Testing 4", 4},
		{"Testing 5", 5},
		{"Testing 6", 6},
		{"Testing 7", 7},
		{"Testing 8", 8},
		{"Testing 9", 9},
		{"Testing 10", 10},
	}
)


func HasResult(r *Results, s string) bool {
	var (
		rs string
		// tp *Results
	)

	if r == nil {
		return false
	}

	for tp := r; tp.nx != nil; tp = tp.nx {
		rs = tp.rs
		// for _, sl := range rs {
			// if str.Compare(sl, s) == 0 {
			if str.Compare(rs, s) == 0 {
				log.Trace.Printf("Found string: [%s]\n", s)
				return true
			}
		// }
	}
	return false
}

func ListResults(r *Results) (int, int) {
	var (
		tp = r
		rc = 0 // Result count
		sc = 0 // slice count
	)

	if tp == nil {
		log.Trace.Printf("Print r == nil\n")
		return rc, sc
	}

	for {
		for m, n := range tp.rs {
			log.Trace.Printf("%d = [%s]\n", m, n)
		}
		rc++
		sc += len(tp.rs)
		log.Trace.Println()
		if tp.nx == nil {
			log.Trace.Printf("rs == nil\n")
			return rc, sc
		}
		tp = tp.nx
	}
}


func TestAddResult(t *testing.T) {
	var (
		r   *Results
		tsc = 0 // tstrings string count
	)

	tstrings := [][]string{
		{"Testing 0-0", "Testing 0-1", "Testing 0-2", "Testing 0-3", "Testing 0-4", "Testing 0-5", "Testing 0-6", "Testing 0-7", "Testing 0-8", "Testing 0-9"},
		{"Testing 1-0", "Testing 1-1", "Testing 1-2", "Testing 1-3", "Testing 1-4", "Testing 1-5", "Testing 1-6", "Testing 1-7", "Testing 1-8", "Testing 1-9"},
		{"Testing 2-0", "Testing 2-1", "Testing 2-2", "Testing 2-3", "Testing 2-4", "Testing 2-5", "Testing 2-6", "Testing 2-7", "Testing 2-8", "Testing 2-9"},
		{"Testing 3-0", "Testing 3-1", "Testing 3-2", "Testing 3-3", "Testing 3-4", "Testing 3-5", "Testing 3-6", "Testing 3-7", "Testing 3-8", "Testing 3-9"},
		{"Testing 4-0", "Testing 4-1", "Testing 4-2", "Testing 4-3", "Testing 4-4", "Testing 4-5", "Testing 4-6", "Testing 4-7", "Testing 4-8", "Testing 4-9"},
		{"Testing 5-0", "Testing 5-1", "Testing 5-2", "Testing 5-3", "Testing 5-4", "Testing 5-5", "Testing 5-6", "Testing 5-7", "Testing 5-8", "Testing 5-9"},
		{"Testing 6-0", "Testing 6-1", "Testing 6-2", "Testing 6-3", "Testing 6-4", "Testing 6-5", "Testing 6-6", "Testing 6-7", "Testing 6-8", "Testing 6-9"},
		{"Testing 7-0", "Testing 7-1", "Testing Yupper 7-2", "Testing 7-3", "Testing 7-4", "Testing 7-5", "Testing 7-6", "Testing 7-7", "Testing 7-8", "Testing 7-9"},
		{"Testing 8-0", "Testing 8-1", "Testing 8-2", "Testing 8-3", "Testing 8-4", "Testing 8-5", "Testing 8-6", "Testing 8-7", "Testing 8-8", "Testing 8-9"},
		{"Testing 9-0", "Testing 9-1", "Testing 9-2", "Testing 9-3", "Testing 9-4", "Testing 9-5", "Testing 9-6", "Testing 9-7", "Testing 9-8", "Testing 9-9"},
	}

	testString := "Scoobie Doobie Dio! Where are you we got so much to live for, Scoobie Doobie Dio what do you do when everyone shits on you."

	log.Trace.Printf("Dumping test list tstring\n\n")
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			log.Trace.Printf("%d-%d [%s]", i, j, tstrings[i][j])
		}
		log.Trace.Println()
		tsc += len(tstrings[i])
	}

	r = new(Results)

	for i := 0; i < 10; i++ {
		r.AddResults(tstrings[i][0])
	}

	log.Trace.Printf("** Listing Results pointer **\n")
	rc, sc := ListResults(r)

	if sc != tsc || rc != len(tstrings) {
		t.Errorf("\n\nInvalid string count:\nResults Added: %d\nStrings Added: %d\n\nResults Listed: %d\nStrings Listed: %d\n\n", len(tstrings), tsc, rc, sc)
	}

	if !UpdateResults(r, "Testing 4-2", testString) {
		t.Errorf("\n\nFailed to update results")
	}

	rc, sc = ListResults(r)
	if sc != tsc || rc != len(tstrings) {
		t.Errorf("\n\nUpdate Results Failed:\nResults Added: %d\nStrings Added: %d\n\nResults Listed: %d\nStrings Listed: %d\n\n", len(tstrings), tsc, rc, sc)
	}

	if !HasResult(r, testString) {
		t.Errorf("Failed to find added string: %s", testString)
	}

	log.Trace.Printf("\n\nString count:\nResults Added: %d\nStrings Added: %d\n\nResults Listed: %d\nStrings Listed: %d\n\n", len(tstrings), tsc, rc, sc)
}
