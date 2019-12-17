package lib

import (
	//	"fmt"
	"go/scanner"
	"go/token"
	log "kwiselog"
	"os"
	"strings"
)

const (
	PBIT = iota
	SENT = iota
	PATB = iota
	PATC = iota
	RES  = iota
)

const (
	PAGE      = iota
	PARAGRAPH = iota
	SENTENCE  = iota
	PATTERN   = iota
	KW_SINGLE = "'"
	MBUFSIZE  = 512
)

type CollectList []Collect
type Mbuf []byte

type CollectQuote struct{}
type CollectLine struct{}
type CollectNext struct{}
type CollectFuncs Command

type Command interface {
	CollectResults(*os.File, *token.FileSet, *Mbuf, *CFToken, *Results, scanner.Scanner) scanner.Scanner
}

type Collect struct {
	st  [4]status        // Status counts
	dr  [4]Directive // Directives [page, paragraph, sentence, pattern]
	pid int              // Paragraph ID
	sid int              // Sentence ID
	pat StringList   // Pattern to match
	pm  StringList   // Next toke pattern to match
}

type Results struct {
	rs string
	nx *Results
}

type status int

const (
	COLLECTQUOTE = iota
	COLLECTLINE  = iota
	COLLECTNEXT  = iota
	ENDOFBUFFER  = "jxLmZJ7gPY"
)

var (
	Seqidx = 0
	SeqInc = func() {
		Seqidx++
	}

	ob   Mbuf // Overflow buffer
	tbuf = func(s []byte) {
		ob = append(ob, s...)
	}

	newline    = ""
	setnewline = func(i string) {
		newline = i
	}

	ttoken  CFToken
	nextTok = func(t CFToken) {
		ttoken = t
	}

	matchP   StringList
	MatchTok = func(s StringList) {
		matchP = matchP[0:0]
		matchP = append(matchP, s...)
	}
)

func NewMbuf() Mbuf {
	return make(Mbuf, MBUFSIZE+len(ENDOFBUFFER)+1)
}

func NewCollection(p int, s int, pt string, pgdir Directive, padir Directive, stdir Directive, patdir Directive, pm StringList) Collect {
	var (
		c Collect
		// cl    StringList
		// st    string
		// instr int
	)

	c.pid = p
	c.sid = s
	c.dr = [...]Directive{pgdir, padir, stdir, patdir}
	cl := MakeTokenString(pt)
	c.pat = append(c.pat, cl...)
	c.pm = append(c.pm, pm...)
	// log.Trace.Printf("NewCollection -- Pattern [%v] length [%d]", c.pat, len(c.pat))
	return (c)
}

func (c CollectLine) CollectResults(df *os.File,
	f *token.FileSet,
	b *Mbuf,
	t *CFToken,
	rs *Results,
	ds scanner.Scanner) scanner.Scanner {

	var (
		tstr string
	)

	log.Trace.Println("####### In CollectLine CollectResults t = ", t)
	eofcnt := 0
	for {
		f, ds = b.NextDataToken(df, f, t, ds)
		log.Trace.Printf("[%v]:- Collect Results [%v] [%d]", t.PLine, t, eofcnt)
		if t.PTok == token.EOF || t.PLit == ENDOFBUFFER {
			log.Trace.Println("####### In CollectLine EOF or EOB = [%v]", t)
			if eofcnt == 0 {
				f, ds = b.NextDataToken(df, f, t, ds)
				eofcnt++
				log.Trace.Printf("Engine: EOF token 1 [%v]\n", t)
				continue
			} else {
				log.Trace.Printf("Engine: EOF token\n")
				break
			}
		}
		eofcnt = 0

		if t.PTok == token.COMMENT {
			log.Trace.Printf("########## got SEMI-COLON [%v] Skipping....\n", t)
			log.Trace.Printf("########## CALLED BREAK [%s]", tstr)
			break
		}
		if t.PLit != "\n" && t.PLit != "\r" {
			tstr = tstr + t.PLit + " "
			if t.PLit == "" {
				tstr = tstr + t.PTok.String()
			}
		}
		// f, ds = b.NextDataToken(df, f, t, ds)
		// log.Trace.Printf("########## Processing [%v] tstr = [%s]", t, tstr)
	}
	log.Trace.Printf("########## **** Adding string:\n\t[%s]\n", tstr)
	rs.AddResults(tstr)
	tstr = ""
	log.Trace.Printf("########## Last toke [%v]", t)
	return ds
}

func (c CollectList) CollectResults(df *os.File,
	f *token.FileSet,
	b *Mbuf,
	t *CFToken,
	rs *Results,
	ds scanner.Scanner) scanner.Scanner {
	log.Trace.Println("In CollectList CollectResults")
	return ds
}

func (c CollectNext) CollectResults(df *os.File,
	f *token.FileSet,
	b *Mbuf,
	t *CFToken,
	rs *Results,
	ds scanner.Scanner) scanner.Scanner {

	var (
		tstr string
		nstr string
	)

	log.Trace.Printf("[%v]:- ####### In CollectNext collecting up too: [%v] of type [%T]\n", t.PLine, matchP, matchP)
	log.Trace.Printf("matchP len [%d]", len(matchP))
	eofcnt := 0
	maxMatch := len(matchP)
	matchCNT := 0
	for {
		log.Trace.Printf("########## Match count [%d] looking for [%d]", matchCNT, maxMatch)
		if matchCNT == maxMatch {
			log.Trace.Printf("########## CollectNext MATCH FOUND CALLED BREAK [%s]", tstr)
			log.Trace.Printf("########## **** Adding string:\n\t[%s]\n", tstr)
			rs.AddResults(tstr)
			tstr = ""
			// log.Trace.Printf("[%v]:- ########## CollectNext Last toke [%v]", t.PLine, t)
			break
		}
		// log.Trace.Printf("########## Collect-Next Processing [%v] tstr = [%s]", t, tstr)

		f, ds = b.NextDataToken(df, f, t, ds)
		log.Trace.Printf("[%v]:- Collect Next -- Results [%v] [%d]", t.PLine, t, eofcnt)
		if t.PTok == token.EOF || t.PLit == ENDOFBUFFER {
			log.Trace.Printf("[%v]:- ####### In CollectNext EOF or EOB = [%v]", t.PLine, t)
			if eofcnt == 0 {
				f, ds = b.NextDataToken(df, f, t, ds)
				eofcnt++
				log.Trace.Printf("Engine: EOF token 1 [%v]\n", t)
				continue
			} else {
				log.Trace.Printf("Engine: EOF token\n")
				break
			}
		}
		eofcnt = 0

		if t.PTok == token.COMMENT {
			t.PLit = "\n"
		}

		nstr = t.PLit
		tstr = tstr + nstr + " "
		if t.PLit == "" {
			nstr = t.PTok.String()
			tstr = tstr + nstr
		}
		log.Trace.Printf("########## - Comparing matchP [%v] with nstr [%v]\n", matchP[matchCNT], nstr)
		if strings.Compare(strings.ToLower(matchP[matchCNT]), strings.ToLower(nstr)) == 0 {
			matchCNT++
		} else {
			matchCNT = 0
		}
	}
	return ds
}

func (c CollectQuote) CollectResults(df *os.File,
	f *token.FileSet,
	b *Mbuf,
	t *CFToken,
	rs *Results,
	ds scanner.Scanner) scanner.Scanner {
	log.Trace.Println("In CollectQuote CollectResults")
	return ds
}

func (o Collect) PutStr() {
	for _, k := range o.pat {
		log.Trace.Printf("[%s]\n", k)
	}
}

func (s Collect) IsParagraphComplete() bool {
	if s.st[PBIT] > 0 {
		return true
	}
	return false
}

func (s Collect) IsSentenceComplete() bool {
	if s.st[SENT] > 0 {
		return true
	}
	return false
}

func (s Collect) IsPatternComplete() bool {
	if s.st[PATB] > 0 {
		return true
	}
	return false
}

func (s Collect) GetPatternDirective(i int) Directive {
	log.Trace.Printf("Getting Pattern [%v] Directive for [%d]", s, i)
	return s.dr[i]
}

func (s Collect) GetSentence() StringList {
	return s.pat
}

func (s Collect) GetMatchToken() StringList {
	return s.pm
}

// CheckParaGraphStatus - Checks the pattern bit status and if set then sets
// CheckParaGraphStatus - sentence/paragrpah bit.
func (l *Collect) CheckParaGraphStatus() {

	if l.st[SENT] == 0 {
		if l.st[PATB] > 0 {
			l.st[PBIT] = 1
			l.st[SENT] = 1
		}
	}
}

func (s *CollectList) ClearPatternStatus() {
	for i := 0; i < len(*s); i++ {
		if (*s)[i].st[PATB] == 0 {
			(*s)[i].st[PATC] = 0
		}
	}
}

func (rs *Results) PutResults() {
	log.Trace.Println("Putting the results")
	for l := rs; l != nil; l = l.nx {
		log.Trace.Printf("Results: [%s]\n", l.rs)
	}
}

func (r *Results) AddResults(rs string) {
	var (
		tp *Results
	)

	log.Trace.Printf("** Adding result [%v] **\n", rs)
	if r == nil {
		r = new(Results)
		r.rs = rs
		r.nx = nil
		log.Trace.Printf("** Added **\n")
		return
	}

	nd := new(Results)
	nd.rs = rs
	nd.nx = nil

	tp = r
	for tp.nx != nil {
		tp = tp.nx
	}
	tp.nx = nd
	log.Trace.Printf("** Added **\n")
}

func UpdateResults(r *Results, s string, u string) bool {

	if r == nil {
		return false
	}

	for tp := r; tp.nx != nil; tp = tp.nx {
		if strings.Compare(tp.rs, s) == 0 {
			log.Trace.Printf("UpdateResult: found [%s]\n", s)
			tp.rs = s
			return true
		}
	}
	return false
}

func (s *Collect) MatchCollectSentence(t CFToken) bool {

	match := t.PLit
	if t.PLit == "" {
		match = t.PTok.String()
	}
	// log.Trace.Printf("[%v]:-*** MatchCollectSentence - s.st[PATC] = [%d] s.pat [%s] len(s.pat) = [%d]", t.PLine, s.st[PATC], s.pat, len(s.pat))
	if int(s.st[PATC]) < len(s.pat) {
		// log.Trace.Printf("*** MATCHing [%s] with [%s]\n", s.pat[s.st[PATC]], match)
		if strings.Compare(strings.ToLower(s.pat[s.st[PATC]]), strings.ToLower(match)) == 0 {
			s.st[PATC]++
		}
	}
	if int(s.st[PATC]) == len(s.pat) {
		s.st[PATB] = 1
	}
	return s.st[PATB] == 1
}

func InitEngine(pg Page) CollectList {

	var (
		ops CollectList
	)

	ops = make(CollectList, 0, 100)
	for i := 0; i < pg.GetParagraphCnt()+1; i++ {
		for j := 0; j < pg.Paras[i].Stcnt+1; j++ {
			// log.Trace.Printf("page: [%T], Para: [%T], Sent: [%T]\nSent: [%v]\n", pg, pg.Paras[i], pg.Paras[i].Stc[j], pg.Paras[i].Stc[j]) // ,pg.Paras[i].Stc[j].stdl, pg.Paras[i].Stc[j].ptdl)
			for l, k := range pg.Paras[i].Stc[j].Pt {
				// log.Trace.Printf("***** INITENGINE -Paras[%d].Stc[%d].Pt = [%v] len(Pt) = [%d]", i, j, k, len(k))
				ops = append(ops, NewCollection(i,
					j,
					k,
					pg.GetPageDirective(),
					pg.Paras[i].GetParagraphDirective(),
					pg.Paras[i].Stc[j].GetSentenceDirective(),
					pg.Paras[i].Stc[j].GetPatternDirective(l),
					pg.Paras[i].Stc[j].GetMatchNextToken()))
			}
		}
	}
	return ops
}

func (by Mbuf) IsWhiteSpace(b byte) bool {
	return b == byte('\t') ||
		b == byte('\n') ||
		b == byte('\v') ||
		b == byte('\f') ||
		b == byte('\r') ||
		b == byte(' ')
}

// GetBuffer: Reads a buffer full of data and returns the buffer ending at
// GetBuffer: a whitespace character. It returns 0 if the end of file is reached.
func (by *Mbuf) GetBuffer(f *os.File) int {
	var (
		x, s int
		err  error
	)
	// log.Trace.Println("In Get buffer")
	*by = NewMbuf()
	if len(ob) > 0 {
		// log.Trace.Println("###>>>Getbuf Overflow len", len(ob))
		*by = (*by)[:len(*by)-len(ob)] // Set the number of chars to read
		s, err = f.Read(*by)
		tb := make(Mbuf, 0, MBUFSIZE+len(ENDOFBUFFER)+1)
		tb = append(tb, ob...)
		tb = append(tb, *by...)
		*by = tb[:len(ob)+s]
		s = len(*by)
		ob = ob[:0]
	} else {
		s, err = f.Read(*by)
		// log.Trace.Println("###>>>Getbuf Read len", len(*by))
	}
	*by = (*by)[:s]
	// Break the buffer at a whitespace character
	if err == nil {
		for x = len(*by) - 1; x > -1; x-- {
			if by.IsWhiteSpace((*by)[x]) {
				break
			}
		}
		// Exclude partial double qutoes since the tokenizer returns them as one token.
		// Yeah -- Could be problematic with wicked long quoted sentence.
		for ; x > -1 && strings.Count(string((*by)[:x]), KW_DQUOTE)%2 > 0; x-- {
			if strings.Count(string((*by)[:x]), KW_DQUOTE)%2 == 0 {
				break
			}
		}
		s = x
		tbuf((*by)[s:])
		*by = (*by)[:s]
		*by = append(*by, []byte(" "+ENDOFBUFFER)...)
	}
	// log.Trace.Println("###>>>Getbuf returning len", len(*by))
	return s
}

func (bs Mbuf) NewKwScanner(s scanner.Scanner) (*token.FileSet, scanner.Scanner) {
	f := token.NewFileSet()
	// log.Trace.Printf("Buffer size == [%d]", len(bs))
	nfile := f.AddFile("", f.Base(), len(bs))
	s.Init(nfile, bs, nil /* no error handler */, 0)
	return f, s
}

// NextDataToken: Takes a scanner object and returns the next token
// NextDataToken: and updated Scanner position. If the line number
// NextDataToken: changes then a NewLine token is returned and the
// NextDataToken: subsequent call returns the cached token.
func (b *Mbuf) NextDataToken(df *os.File, f *token.FileSet, t *CFToken, s scanner.Scanner) (*token.FileSet, scanner.Scanner) {

	if t.PTok == token.EOF || t.PLit == ENDOFBUFFER {
		log.Trace.Printf("In NextDataToken found EOF or ENDOFBUFFER")
		if b.GetBuffer(df) > 0 {
			f, s = b.NewKwScanner(s)
		} else {
			return f, s
		}
	}

	if newline == "y" {
		//setnewline(strings.Split(t.PLine, ":")[0])
		setnewline(strings.Split(ttoken.PLine, ":")[0])
		log.Trace.Printf("[%v]:- ttoken.PLIne [%v] In NextDataToken Returning Buffered token [%v] newline [%v]", t.PLine, ttoken.PLine, ttoken, newline)
		*t = ttoken
		return f, s
	}

	t.Pos, t.PTok, t.PLit = s.Scan()
	t.PLine = f.Position(t.Pos).String()
	ln := strings.Split(t.PLine, ":")[0]
	// log.Trace.Printf("Called scan token [%v]\n", t)

	if t.PLit == "\n" || t.PLit == "\r" || newline != ln {
		// log.Trace.Printf("[%v]:- NextDataToken newline toke [%v] newline [%v] ln[%v]\n", t.PLine, t, newline, ln)
		setnewline("y")
		nextTok(*t)
		t.PTok = token.COMMENT
		t.PLit = ""
	}
	return f, s
}
