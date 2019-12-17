package lib

import (
// log "kwiselog"
)

const (
	KW_PARAGRAPH_CAPACITY = 10
	KW_PARAGRAPH_EXTEND   = 5
)

var (
	ParaIdx = -1
	IncPara = func() {
		ParaIdx++
		SentIdx = -1
	}

	SentIdx = -1 /* Set to something stupid to shutup compiler */
	IncSent = func() {
		SentIdx++
		SrchIdx = -1
	}

	SrchIdx = -1
	IncSrch = func() {
		SrchIdx++
	}
)

func MakeParagraph(p *Paragraph) *Paragraph {
	p = new(Paragraph)
	p.Stc = make([]Pattern, KW_PARAGRAPH_EXTEND, KW_PARAGRAPH_CAPACITY)
	return p
}

// AddParagraph: Using a varidac function creates a slice for np. This
// AddParagraph: allows us to use copy to copy the contents of np.
//

func (p *ParagraphList) AddParagraph(np ...Paragraph) {
	IncPara()
	// log.Trace.Println("ADD PARAGRAPH ==> ", np[0].pna)
	contents := ParaIdx + len(np)
	if contents > cap(*p) {
		ns := contents*3/2 + 1
		ts := make(ParagraphList, contents, ns)
		copy(ts, *p)
		*p = ts
	}
	*p = (*p)[:contents]
	copy((*p)[ParaIdx:], np)
}

// func AddPattern(p []Pattern, i func(), idx *int, np ...Pattern) []Pattern {
func (p *PatternList) AddPattern(i func(), idx *int, np ...Pattern) {
	i()
	// log.Trace.Println("ADD PATTERN ==> ", *idx)
	contents := *idx + len(np)
	if contents > cap(*p) {
		ns := contents*3/2 + 1
		ts := make(PatternList, contents, ns)
		copy(ts, *p)
		*p = ts
	}
	*p = (*p)[:contents]
	copy((*p)[*idx:], np)
}

func (p *StringList) AddSearch(np ...string) {
	IncSrch()
	// log.Trace.Println("ADD Search ==> ", SrchIdx)
	contents := SrchIdx + len(np)
	if contents > cap(*p) {
		ns := contents*3/2 + 1
		ts := make(StringList, contents, ns)
		copy(ts, *p)
		*p = ts
	}
	*p = (*p)[:contents]
	copy((*p)[SrchIdx:], np)
}

// NewPage: Allocates storage for the page.
//
func NewPage() Page {
	var pg Page
	pg.Paras = make([]Paragraph, 5, 10)
	return pg
}

func ToString(x StringList) string {
	var s string

	for _, k := range x {
		s = s + string(k)
	}
	return s
}
