package lib

// kwisemap - creates a map of kwise commands/directives with values starting at 4020
import (
	//	"engine"
	"fmt"
	"go/scanner"
	"go/token"
	log "kwiselog"
	"os"
	str "strings"
)

const (
	KW_COLON   = ":"
	KW_LBRACE  = "{"
	KW_RBRACE  = "}"
	KW_DQUOTE  = "\""
	KW_SQUOTE  = "'"
	KW_DRSTART = KW_COLON
	KW_TFOUND  = true
	KW_NFOUND  = false
	KW_TOKSEP  = "-"
	KW_DEFINED = 1
	KW_SEMI    = ";"
	KW_EQUAL   = "="
)

// Page Constants
const (
	INSEQUENCE = 2 * 1 << iota
	UNSEQUENCE = 2 * 1 << iota
	PGMIXED    = 2 * 1 << iota
	PGLOWER    = 2 * 1 << iota
)

// Paragraph Constants
const (
	ONEOF   = 2 * 1 << iota
	SOMEOF  = 2 * 1 << iota
	ANYOF   = 2 * 1 << iota
	ALLOF   = 2 * 1 << iota
	BEGIN   = 2 * 1 << iota
	END     = 2 * 1 << iota
	PAMIXED = 2 * 1 << iota
	PALOWER = 2 * 1 << iota
)

// Sentence Constants
const (
	COLLECT_QUOTE = 2 * 1 << iota
	COLLECT_LINE  = 2 * 1 << iota
	COLLECT_NEXT  = 2 * 1 << iota
	COLLECT       = 2 * 1 << iota
	COLLECT_COMMA = 2 * 1 << iota
	CONCAT_TAB    = 2 * 1 << iota
	CONCAT        = 2 * 1 << iota
	COMPLETE      = 2 * 1 << iota
	ISTHERE       = 2 * 1 << iota
	EXACT         = 2 * 1 << iota
	ANY           = 2 * 1 << iota
	EXCLUDE       = 2 * 1 << iota
	FOLLOW        = 2 * 1 << iota
	REGARDLESS    = 2 * 1 << iota
	STMIXED       = 2 * 1 << iota
	STLOWER       = 2 * 1 << iota
	START         = 2 * 1 << iota
)

var (
	pageDirectives = command{
		":insequence": INSEQUENCE,
		":unsequence": UNSEQUENCE,
		":casemixed":  PGMIXED,
		":caselower":  PGLOWER,
	}

	// Decision
	paragraphDirectives = command{
		":oneof":     ONEOF,
		":someof":    SOMEOF,
		":anyof":     ANYOF,
		":allof":     ALLOF,
		":begin":     BEGIN,
		":end":       END,
		":casemixed": PAMIXED,
		":caselower": PALOWER,
	}

	sentenceDirectives = command{
		":isthere":    ISTHERE,
		":exact":      EXACT,
		":any":        ANY,
		":exclude":    EXCLUDE,
		":follow":     FOLLOW,
		":regardless": REGARDLESS,
		":casemixed":  STMIXED,
		":caselower":  STLOWER,
		":start":      START,
	}

	// Collection
	searchDirectives = command{
		":collect_quote": COLLECT_QUOTE,
		":collect_line":  COLLECT_LINE,
		":collect_next":  COLLECT_NEXT,
		":collect_comma": COLLECT_COMMA,
		":concat_tab":    CONCAT_TAB,
		":concat":        CONCAT,
		":complete":      COMPLETE,
	}

	Fset *token.FileSet
)

type CFToken struct {
	Pos   token.Pos
	PLine string
	PTok  token.Token
	PLit  string
}

type CFTErr struct {
	Err  string
	Line token.Pos
	Code int
}

type Directive uint
type command map[string]Directive
type PatternList []Pattern
type StringList []string
type DirectiveList []Directive

type Page struct {
	pdl   Directive // Page directive
	Paras ParagraphList
	pacnt int
}

type Paragraph struct {
	phdl  Directive   // Directives for paragraph
	pna   string      // Name
	Stc   PatternList // Sentences
	Stcnt int
}

type ParagraphList []Paragraph

type Pattern struct {
	pi   int           // Pattern index
	stdl Directive     // Directives for Sentences
	ptdl DirectiveList // Directives for patterns
	sna  string        // sentence
	Pt   StringList    // pattern to match
	Nt   StringList    // Match next token
	st   bool          // status
}

func (c command) printCommand() {
	for k, v := range c {
		fmt.Println("Key = Value ", k, v)
	}
}

func (c command) GetDirective(k string) (bool, Directive) {
	return c[k] != 0, c[k]
}

func (d Directive) SetDirective(v Directive) Directive {
	return d ^ v
}

func (d Directive) ClearDirective(v Directive) Directive {
	return d &^ v
}

func (d Directive) IsDirectiveSet(v Directive) bool {
	//	log.Trace.Printf("!!!!!!!!! Directive [%d] Checking [%v] is set [%t]", d, v, d&v == v)
	return d&v == v
}

// IsToken: Returns true if string represents a token of the map.
func (c command) IsToken(t string) bool {
	return c[t] != 0
}

// NextToken: Takes a scanner object and returns the next token
// NextToken: and updated Scanner position.
func NextToken(f *token.FileSet, s scanner.Scanner) (CFToken, scanner.Scanner) {
	var t CFToken

	t.Pos, t.PTok, t.PLit = s.Scan()
	t.PLine = f.Position(t.Pos).String()

	return t, s
}

func (pg Page) GetParagraphCnt() int {
	return pg.pacnt
}

func (pg Page) GetPageDirective() Directive {
	return pg.pdl
}

func (para Paragraph) GetParagraphDirective() Directive {
	return para.phdl
}

func (sent Pattern) GetSentenceDirective() Directive {
	return sent.stdl
}

func (pt Pattern) GetPatternDirective(i int) Directive {
	return pt.ptdl[i]
}

func (pt Pattern) GetMatchNextToken() StringList {
	return pt.Nt
}

func (pt Pattern) GetPatternDirectiveList() DirectiveList {
	return pt.ptdl
}

func (pat Pattern) PatternDirectiveTotal() Directive {
	var x Directive

	for _, k := range pat.ptdl {
		x = x + k
	}
	return x
}

func (sent Pattern) GetSentenceName() string {
	return sent.sna
}

func (pg Page) CheckDirectives() bool {

	var (
		drlist DirectiveList
		tstr   string
	)

	x := 0
	for i := 0; i < pg.GetParagraphCnt()+1; i++ {
		for j := 0; j < pg.Paras[i].Stcnt+1; j++ {
			drlist = pg.Paras[i].Stc[j].GetPatternDirectiveList()
			for l := 0; l < len(drlist)-l; l++ {
				if drlist[l] > 0 && len(drlist)-1 > 0 {
					tstr = "***** Error: Sentence [" + pg.Paras[i].Stc[j].sna + "]\nsearch directive " + pg.Paras[i].Stc[j].Pt[l] + "\ndefines directive [" + drlist[l].searchDirectives("") + "]"
					log.Trace.Printf("%s\ndrlist [%v]\nl = [%d]\ndrlist len(%d)\n", tstr, drlist, l, len(drlist)-1)
					x = 1
				}
			}
		}
	}
	return x == 0
}

// ProcessPage takes a token and determines if it contains the start to the page
// command and then proccesses the command. ProcessPage returns true if successful.
// Otherwise ProcessPage returns false.
//
func (pg *Page) ProcessPage(f *token.FileSet, s scanner.Scanner, t CFToken) (scanner.Scanner, bool) {

	var (
		tokCnt  = 0
		tokLine = t.PLine
		st      bool
		lc      = false
	)

	// log.Trace.Printf("Page -- Process start.\n")
	if pg.pdl.IsDirectiveSet(KW_DEFINED) {
		return s, false
	}

	pg.pdl.SetDirective(KW_DEFINED)
	for !lc {
		t, s = NextToken(f, s)
		//		log.Trace.Printf("PAGE: parsing token [%v].\n", t)
		switch {
		case t.PTok == token.EOF:
			log.Trace.Printf("[%v]:- Page EOF found.\n", t.PLine)
			lc, st = true, true
			break

		case str.Compare(t.PTok.String(), KW_LBRACE) == 0:
			// log.Trace.Printf("[%v]:- Page KW_LBRACE found.\n", t.PLine)
			lc, st = true, false
			break

		case str.Compare(t.PTok.String(), KW_DRSTART) == 0:
			//			log.Trace.Printf("[%v]:- Page KW_DRSTART found.\n", t.PLine)
			t, s = NextToken(f, s)
			//			log.Trace.Printf("[%v]:- Next token after KW_DRSTART is [%v]\n", t.PLine, t)
			if pageDirectives.IsToken(str.ToLower(KW_DRSTART + t.PLit)) {
				t.PLit = str.ToLower(KW_DRSTART + t.PLit)
				//				log.Trace.Printf("[%v]:- Page: [%v] directive found.\n", t.PLine, t)

				switch {
				case str.Compare(t.PLit, ":insequence") == 0:
					pg.pdl = pg.pdl.SetDirective(pageDirectives[t.PLit])

				case str.Compare(t.PLit, ":unsequence") == 0:
					pg.pdl = pg.pdl.SetDirective(pageDirectives[t.PLit])

				case str.Compare(t.PLit, ":caselower") == 0:
					pg.pdl = pg.pdl.SetDirective(pageDirectives[t.PLit])

				case str.Compare(t.PLit, ":casemixed") == 0:
					pg.pdl = pg.pdl.SetDirective(pageDirectives[t.PLit])

				default:
					log.Trace.Printf("[%v]:- Page OOPS!!!!  ** UNACCOUNTED ** Directive [%v] found.\n", t.PLine, t)
					lc = true
				} // End switch
			} else { // end Build Directives
				log.Trace.Printf("[%v]:- Page - Invalid directive [%s]. ", t.PLine, t.PLit)
				lc, st = true, true
			}

		case str.Compare(t.PTok.String(), ";") == 0: /* Need to skip this. Noop */
			//			log.Trace.Println("Skipping ;", t)

		default:
			log.Trace.Printf("[%v]:- Page Unknown token found [%v]. plit - [%v]", t.PLine, t)
			lc, st = true, true
			break
		} // End Switch
		tokCnt++
		if tokCnt == 100 {
			log.Trace.Printf("[%v]: Page missing closing brace.\n", t.PLine)
			lc, st = true, true
			t.PLine = tokLine
			break
		}
	} // End For
	return s, st
}

func (pg *Page) ProcessParagraph(f *token.FileSet, s scanner.Scanner, t CFToken) (scanner.Scanner, bool) {
	var (
		tp      Paragraph
		tokCnt  = 0
		tokLine = t.PLine
		st      bool
		lc      = false
	)

	// log.Trace.Printf("Paragraph -- Process start.\n\n")
	for !lc {
		t, s = NextToken(f, s)
		// log.Trace.Printf("Para: parsing token [%v]\n", t)
		switch {
		case t.PTok == token.EOF:
			log.Trace.Printf("[%v]:- Para EOF found.\n", t.PLine)
			lc, st = true, true
			break

		case str.Compare(t.PTok.String(), KW_LBRACE) == 0:
			// log.Trace.Printf("[%v]:- Para KW_LBRACE found.\n", t.PLine)
			lc, st = true, false
			break

		case str.Compare(t.PTok.String(), KW_DRSTART) == 0:
			// log.Trace.Printf("[%v]:- Para KW_DRSTART found.\n", t.PLine)
			t, s = NextToken(f, s)
			if paragraphDirectives.IsToken(str.ToLower(KW_DRSTART + t.PLit)) {
				t.PLit = str.ToLower(KW_DRSTART + t.PLit)
				// log.Trace.Printf("[%v]:- Para: [%v] directive found [%v].\n", t.PLine, t, t.PLit)
				switch {
				case str.Compare(t.PLit, ":casemixed") == 0:
					tp.phdl = tp.phdl.SetDirective(paragraphDirectives[t.PLit])

				case str.Compare(t.PLit, ":caselower") == 0:
					tp.phdl = tp.phdl.SetDirective(paragraphDirectives[t.PLit])

				case str.Compare(t.PLit, ":oneof") == 0:
					tp.phdl = tp.phdl.SetDirective(paragraphDirectives[t.PLit])

				case str.Compare(t.PLit, ":someof") == 0:
					tp.phdl = tp.phdl.SetDirective(paragraphDirectives[t.PLit])

				case str.Compare(t.PLit, ":anyof") == 0:
					tp.phdl = tp.phdl.SetDirective(paragraphDirectives[t.PLit])

				case str.Compare(t.PLit, ":allof") == 0:
					tp.phdl = tp.phdl.SetDirective(paragraphDirectives[t.PLit])

				case str.Compare(t.PLit, ":begin") == 0:
					tp.phdl = tp.phdl.SetDirective(paragraphDirectives[t.PLit])

				case str.Compare(t.PLit, ":end") == 0:
					tp.phdl = tp.phdl.SetDirective(paragraphDirectives[t.PLit])

				default:
					log.Trace.Printf("[%v]:- Para Invalid Directive [%v] found.\n", t.PLine, t)
					lc = true
				} // end directives switch
			} // end Build Directives

		case str.Compare(t.PTok.String(), ";") == 0: /* Need to skip this. Noop */
			log.Trace.Println("Skipping ;")

		default:
			//log.Trace.Printf("[%v]:- Para name found [%v].\n", t.PLine, t)
			tp.pna = t.PLit
		} // End Switch
		tokCnt++
		if tokCnt == 100 {
			log.Trace.Printf("[%v]: Paragraph missing closing brace.\n", t.PLine)
			lc, st = true, true
			t.PLine = tokLine
			break
		}
	} // End For

	if !st {
		pg.Paras.AddParagraph(tp)
		pg.pacnt = ParaIdx
		// log.Trace.Printf("[%v]: Adding Paragraph to page @ index [%d].", t.PLine, ParaIdx)
	}
	return s, st
}

func MakeTokenString(tokestr string) StringList {

	var src []byte

	if len(tokestr) > 1 && str.Compare(string(tokestr[0]), KW_DQUOTE) == 0 && str.Compare(string(tokestr[len(tokestr)-1]), KW_DQUOTE) == 0 {
		log.Trace.Printf("MakeTokenString str [%v] stripped [%v]\n", tokestr, tokestr[1:len(tokestr)-1])
		src = []byte(tokestr[1 : len(tokestr)-1])
	} else {
		log.Trace.Printf("MakeTokenString str [%v]\n", tokestr)
		src = []byte(tokestr)
	}
	mslice := make(StringList, 0, 20)

	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, 0)

	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		// log.Trace.Printf("tok [%v] lit [%v]\n", tok, lit)
		if lit == "" {
			mslice = append(mslice, tok.String())
		} else if tok.String() == ";" && lit == ";" {
			mslice = append(mslice, lit)
		} else if tok.String() != ";" {
			mslice = append(mslice, lit)
		}
	}
	return mslice
}

func (pa *Pattern) getCollectNextToken(f *token.FileSet, s scanner.Scanner, t CFToken) scanner.Scanner {

	t, s = NextToken(f, s)
	// log.Trace.Printf("[%v]: ^^^^^^^In getCollectNextToken\n", t.PLine)
	if t.PTok == token.EOF {
		log.Trace.Printf("getCollectNextToken: EOF found\n")
		return s
	}

	if str.Compare(t.PTok.String(), KW_EQUAL) == 0 {
		//	log.Trace.Printf("^^^^^^^ Got equal start [%s]\n", t.PTok.String())
		t, s = NextToken(f, s)
		if str.Compare(t.PTok.String(), "STRING") != 0 {
			log.Trace.Printf("[%v]: getCollectNextToken - Invalid next token [%v] aborting...\n", t.PLine, t)
			os.Exit(1)
		}
	}
	pa.Nt = append(pa.Nt, t.PLit)
	// log.Trace.Printf("[%v]: ^^^^^ Got Collect till next token token [%v]\n", t.PLine, pa.Nt)
	return s
}

func (pa *Pattern) processPattern(f *token.FileSet, s scanner.Scanner, t CFToken) (scanner.Scanner, bool) {
	var (
		sr string
		st bool
		lc = false
		dr Directive
	)

	dr = 0
	// log.Trace.Printf("Pattern -- Collection start for Sentence.")
	for !lc {
		t, s = NextToken(f, s)
		// log.Trace.Printf("[%v]:- Pattern: parsing token [%v]\n\ntoken.PTOk.String() [%s] token.PLit [%s]", t.PLine, t, t.PTok.String(), t.PLit)
		if str.Compare(t.PTok.String(), KW_RBRACE) == 0 {
			// log.Trace.Printf("[%v]:- Pattern: KW_RBRACE found.\n", t.PLine)
			if sr != "" {
				//log.Trace.Printf("[%v]:- Pattern: Recording pattern before leaving: [%s]", t.PLine, sr)
				pa.Pt.AddSearch(str.TrimSpace(sr))
				pa.ptdl = append(pa.ptdl, dr)
				dr = 0
				//	log.Trace.Printf("[%v]:- Added to Paragraph [%v], Sentence [%v], Pattern [%v]", t.PLine, ParaIdx, SentIdx, SrchIdx)
			}
			lc, st = true, false
			break
		}

		if str.Compare(t.PTok.String(), KW_SEMI) == 0 {
			//	log.Trace.Printf("[%v]:- Pattern: found SENI-COLON - recording string [%s]", t.PLine, sr)
			pa.Pt.AddSearch(str.TrimSpace(sr))
			pa.pi = SrchIdx
			pa.ptdl = append(pa.ptdl, dr)
			dr = 0
			sr = ""
			//	log.Trace.Printf("[%v]:- Added to Paragraph [%v], Sentence [%v], Pattern [%v]", t.PLine, ParaIdx, SentIdx, SrchIdx)
			continue
		}

		if str.Compare(t.PTok.String(), KW_DRSTART) == 0 {
			//		log.Trace.Printf("[%v]:- Para KW_DRSTART found.\n", t.PLine)
			t, s = NextToken(f, s)
			if searchDirectives.IsToken(str.ToLower(KW_DRSTART + t.PLit)) {
				t.PLit = str.ToLower(KW_DRSTART + t.PLit)
				//			log.Trace.Printf("[%v]:-@@@@@@@@@@@ Para: [%v] directive found [%v].\n", t.PLine, t, t.PLit)
				switch {
				case str.Compare(t.PLit, ":collect_quote") == 0:
					dr = dr ^ searchDirectives[t.PLit]
					continue

				case str.Compare(t.PLit, ":collect_line") == 0:
					//				log.Trace.Printf("[%v]:- @@@@@@@@@@@ setting directive for :collect_line [%d]", t.PLine, searchDirectives[t.PLit])
					dr = dr ^ searchDirectives[t.PLit]
					continue

				case str.Compare(t.PLit, ":collect_next") == 0:
					//				log.Trace.Printf("[%v]:- @@@@@@@@@@@ setting directive for :collect_next [%d]", t.PLine, searchDirectives[t.PLit])
					dr = dr ^ searchDirectives[t.PLit]
					s = pa.getCollectNextToken(f, s, t)
					continue

				case str.Compare(t.PLit, ":collect_comma") == 0:
					dr = dr ^ searchDirectives[t.PLit]
					continue

				case str.Compare(t.PLit, ":concat_tab") == 0:
					dr = dr ^ searchDirectives[t.PLit]
					continue

				case str.Compare(t.PLit, ":concat") == 0:
					dr = dr ^ searchDirectives[t.PLit]
					continue

				case str.Compare(t.PLit, ":complete") == 0:
					dr = dr ^ searchDirectives[t.PLit]
					continue

				default: /* Invalid directive or a valid search string. Either way... Keep it moving! */
					sr = sr + t.PLit + " "
					continue
				}
			} else { /* Invalid directive or a valid search string. Either way... Keep it moving! */
				sr = sr + KW_DRSTART + t.PLit + " "
				continue
			}
		} /* If directive */

		if t.PLit == "" {
			// log.Trace.Printf("GOT Punctuation - [%v] PLit [%v]", t.PTok.String(), t.PLit)
			t.PLit = t.PTok.String()
		}
		sr = sr + t.PLit + " "
	}
	return s, st
}

func (pg *Page) ProcessSentence(f *token.FileSet, s scanner.Scanner, t CFToken) (scanner.Scanner, bool) {
	var (
		tp      Pattern
		tokCnt  = 0
		tokLine = t.PLine
		st      bool
		lc      = false
	)

	// log.Trace.Printf("Sentence -- Process start.\n\n")
	tp.ptdl = make(DirectiveList, 0, 10)
	for !lc {
		t, s = NextToken(f, s)
		// log.Trace.Printf("Sent: parsing token [%v]\n", t)
		switch {
		case t.PTok == token.EOF:
			// log.Trace.Printf("[%v]:- Sent EOF found.\n", t.PLine)
			lc, st = true, true
			break

		case str.Compare(t.PTok.String(), KW_LBRACE) == 0:
			// log.Trace.Printf("[%v]:- Sent KW_LBRACE found.\n", t.PLine)
			s, st = tp.processPattern(f, s, t)
			lc = true
			break

		case str.Compare(t.PTok.String(), KW_RBRACE) == 0:
			// log.Trace.Printf("[%v]:- Sent KW_RBRACE found.\n", t.PLine)
			lc, st = true, false
			break

		case str.Compare(t.PTok.String(), KW_DRSTART) == 0:
			// log.Trace.Printf("[%v]:- Sent KW_DRSTART found.\n", t.PLine)
			t, s = NextToken(f, s)
			if sentenceDirectives.IsToken(str.ToLower(KW_DRSTART + t.PLit)) {
				t.PLit = str.ToLower(KW_DRSTART + t.PLit)
				// log.Trace.Printf("[%v]:- Sent: [%v] directive found [%v].\n", t.PLine, t, t.PLit)
				switch {
				case str.Compare(t.PLit, ":isthere") == 0:
					tp.stdl = tp.stdl.SetDirective(sentenceDirectives[t.PLit])

				case str.Compare(t.PLit, ":exact") == 0:
					tp.stdl = tp.stdl.SetDirective(sentenceDirectives[t.PLit])

				case str.Compare(t.PLit, ":any") == 0:
					tp.stdl = tp.stdl.SetDirective(sentenceDirectives[t.PLit])

				case str.Compare(t.PLit, ":exclude") == 0:
					tp.stdl = tp.stdl.SetDirective(sentenceDirectives[t.PLit])

				case str.Compare(t.PLit, ":follow") == 0:
					tp.stdl = tp.stdl.SetDirective(sentenceDirectives[t.PLit])

				case str.Compare(t.PLit, ":regardless") == 0:
					tp.stdl = tp.stdl.SetDirective(sentenceDirectives[t.PLit])

				case str.Compare(t.PLit, ":casemixed") == 0:
					tp.stdl = tp.stdl.SetDirective(sentenceDirectives[t.PLit])

				case str.Compare(t.PLit, ":caselower") == 0:
					tp.stdl = tp.stdl.SetDirective(sentenceDirectives[t.PLit])

				case str.Compare(t.PLit, ":start") == 0:
					tp.stdl = tp.stdl.SetDirective(sentenceDirectives[t.PLit])

				default:
					log.Trace.Printf("[%v]:- Sent Invalid Directive [%v] found.\n", t.PLine, t)
					lc = true
				} // end directives switch
			} // end Build Directives

		case str.Compare(t.PTok.String(), ";") == 0: /* Need to skip this. Noop */
			// log.Trace.Println("Skipping ;")

		default:
			// log.Trace.Printf("[%v]:- Sent name found [%v].\n", t.PLine, t)
			tp.sna = t.PLit
		} // End Switch
		tokCnt++
		if tokCnt == 100 {
			log.Trace.Printf("[%v]: Sentence missing closing brace.\n", t.PLine)
			lc, st = true, true
			t.PLine = tokLine
			break
		}
	} // End For

	if !st {
		pg.Paras[ParaIdx].Stc.AddPattern(IncSent, &(SentIdx), tp)
		pg.Paras[ParaIdx].Stcnt = SentIdx
		// log.Trace.Printf("[%v]: Added @ Paragraph index [%d]. Sentence index [%d] Done.", t.PLine, ParaIdx, SentIdx)
	}
	return s, st
}
