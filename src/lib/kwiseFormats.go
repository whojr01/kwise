package lib

import (
	"fmt"
	"strconv"
)

func (da Directive) pageDirectives(t string) string {
	s := ""
	if da.IsDirectiveSet(INSEQUENCE) {
		s = s + fmt.Sprintf("%sINSEQUENCE\n", t)
	}

	if da.IsDirectiveSet(UNSEQUENCE) {
		s = s + fmt.Sprintf("%sUNSEQUENCE\n", t)
	}

	if da.IsDirectiveSet(PGMIXED) {
		s = s + fmt.Sprintf("%sMIXED CASE\n", t)
	}

	if da.IsDirectiveSet(PGLOWER) {
		s = s + fmt.Sprintf("%sLOWER CASE\n", t)
	}

	return s
}

func (da Directive) paraDirective(t string) string {

	s := ""
	if da.IsDirectiveSet(ONEOF) {
		s = s + fmt.Sprintf("%sONEOF\n", t)
	}

	if da.IsDirectiveSet(SOMEOF) {
		s = s + fmt.Sprintf("%sSOMEOF\n", t)
	}

	if da.IsDirectiveSet(ANYOF) {
		s = s + fmt.Sprintf("%sANYOF\n", t)
	}

	if da.IsDirectiveSet(ALLOF) {
		s = s + fmt.Sprintf("%sALLOF\n", t)
	}

	if da.IsDirectiveSet(BEGIN) {
		s = s + fmt.Sprintf("%sBEGIN\n", t)
	}

	if da.IsDirectiveSet(END) {
		s = s + fmt.Sprintf("%sEND\n", t)
	}

	if da.IsDirectiveSet(PAMIXED) {
		s = s + fmt.Sprintf("%sMIXED CASE\n", t)
	}

	if da.IsDirectiveSet(PALOWER) {
		s = s + fmt.Sprintf("%sLOWER CASE\n", t)
	}

	return s
}

func (da Directive) sentenceDirectives(t string) string {
	s := ""
	if da.IsDirectiveSet(ISTHERE) {
		s = s + fmt.Sprintf("%sISTHERE\n", t)
	}

	if da.IsDirectiveSet(EXACT) {
		s = s + fmt.Sprintf("%sEXACT\n", t)
	}

	if da.IsDirectiveSet(ANY) {
		s = s + fmt.Sprintf("%sANY\n", t)
	}

	if da.IsDirectiveSet(EXCLUDE) {
		s = s + fmt.Sprintf("%sEXCLUDE\n", t)
	}

	if da.IsDirectiveSet(FOLLOW) {
		s = s + fmt.Sprintf("%sFOLLOW\n", t)
	}

	if da.IsDirectiveSet(REGARDLESS) {
		s = s + fmt.Sprintf("%sREGARDLESS\n", t)
	}

	if da.IsDirectiveSet(STMIXED) {
		s = s + fmt.Sprintf("%sMIXED CASE\n", t)
	}

	if da.IsDirectiveSet(STLOWER) {
		s = s + fmt.Sprintf("%sLOWER CASE\n", t)
	}

	if da.IsDirectiveSet(START) {
		s = s + fmt.Sprintf("%sSTART\n", t)
	}
	return s
}

func (da Directive) searchDirectives(t string) string {
	s := ""
	if da.IsDirectiveSet(COLLECT_QUOTE) {
		s = s + fmt.Sprintf("%sCOLLECT QUOTE\n", t)
	}

	if da.IsDirectiveSet(COLLECT_LINE) {
		s = s + fmt.Sprintf("%sCOLLECT LINE\n", t)
	}

	if da.IsDirectiveSet(COLLECT_NEXT) {
		s = s + fmt.Sprintf("%sCOLLECT NEXT\n", t)
	}

	if da.IsDirectiveSet(COLLECT_COMMA) {
		s = s + fmt.Sprintf("%sCOLLECT COMMA\n", t)
	}

	if da.IsDirectiveSet(CONCAT_TAB) {
		s = s + fmt.Sprintf("%sCONCAT TAB\n", t)
	}

	if da.IsDirectiveSet(CONCAT) {
		s = s + fmt.Sprintf("%sCONCAT\n", t)
	}

	if da.IsDirectiveSet(COMPLETE) {
		s = s + fmt.Sprintf("%sCOMPLETE\n", t)
	}
	return s
}

func (da DirectiveList) patternDirectives(t string) string {

	var s string
	for _, k := range da {
		s = s + k.searchDirectives(t) + "\n"
	}
	return s
}

func (da DirectiveList) patternDirectiveTotal() Directive {
	var x Directive

	for _, k := range da {
		x = x + k
	}
	return x
}

func (pa Paragraph) String() string {
	s := "\n\tParagraph Dump\n\n"
	s = s + "\tName [%s]\n"
	s = s + "\tDirectives: [%d]\n\t===========\n%s\n"
	return fmt.Sprintf(s, pa.pna, pa.phdl, pa.phdl.paraDirective("\t"))
}

func (stc Pattern) String() string {
	s := "\n\t\tSentence Dump\n\n"
	s = s + "\t\tName [%s]\n"
	s = s + "\t\tDirectives: [%d]\n\t\t===========\n%s\n"
	s = s + "\n\t\tPattern to Match\n"
	s = s + "\t\t------------------------\n"
	for _, j := range stc.Pt {
		s = s + "\t\t" + j + "\n"
	}
	s = s + "\n\t\t\tPattern Directives: [%d]\n\t\t\t===========\n%s\n"
	s = s + "\t\t\tNext token match set to: %s\n"
	s = s + "\n\t\tNumber of Patterns [" + strconv.Itoa(stc.pi+1) + "]"
	return fmt.Sprintf(s, stc.sna, stc.stdl, stc.stdl.sentenceDirectives("\t\t"), stc.ptdl.patternDirectiveTotal(), stc.ptdl.patternDirectives("\t\t\t"), stc.Nt)
}

func (p Page) String() string {
	var (
		k Paragraph
		l Pattern
	)

	s := "\nPage Dump\n\n"
	s = s + "Directive: [%d]\n===========\n%s\n"
	s = s + "\tNumber of Paragraphs: [%d]\n"
	s = s + "\t*************************\n"
	for i := 0; i < p.pacnt+1; i++ {
		k = p.Paras[i]
		s = s + k.String()
		s = s + "\tParagraph Number [" + strconv.Itoa(i) + "] of [" + strconv.Itoa(p.pacnt) + "]\n"
		s = s + "\t------------------------\n"
		for j := 0; j < p.Paras[i].Stcnt+1; j++ {
			l = p.Paras[i].Stc[j]
			s = s + l.String()
			s = s + "\n\t\tSentence Number [" + strconv.Itoa(j) + "] of [" + strconv.Itoa(p.Paras[i].Stcnt) + "]"
			s = s + "\n\t\t------------------------\n"
		}
	}
	return fmt.Sprintf(s, p.pdl, p.pdl.pageDirectives("  "), p.pacnt)
}

// order= p.pdl,len(p.paras), p.pcnt
