package main

import (
	"files"
	"flag"
	"fmt"
	"go/scanner"
	"go/token"
	log "kwiselog"
	"lib"
	"os"
	str "strings"
)

var (
	df    *os.File
	dFile string
	tFile string
	tr    bool
	Mp    []string
)

func init() {
	flag.StringVar(&tFile, "pageDescription", "", "Location of the page description file that describes the data to be collected from the specified data file.")
	flag.StringVar(&dFile, "dataFile", "", "Location of the data file to scan.")
	flag.BoolVar(&tr, "trace", false, "Turn on debugging.")
}

func main() {
	var (
		s   scanner.Scanner
		ds  scanner.Scanner
		bs  []byte
		err error
		f   *token.FileSet
		// b      files.Mbuf = make(files.Mbuf, files.MBUFSIZE)
		t              lib.CFToken
		st             bool
		b              lib.Mbuf
		rs             lib.Results
		Cq             lib.CollectQuote
		Cl             lib.CollectLine
		Cn             lib.CollectNext
		CollectRoutine = []lib.CollectFuncs{
			Cq, Cl, Cn,
		}
	)

	flag.Parse()
	if _, err = os.Stat(dFile); os.IsNotExist(err) {
		fmt.Printf("Error: %s\n", str.Split(err.Error(), "CreateFile ")[1])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if _, err = os.Stat(tFile); os.IsNotExist(err) {
		fmt.Printf("Error: %s\n", str.Split(err.Error(), "CreateFile ")[1])
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.InitLog()

	log.Info.Printf("\n** Program run start **")
	if tr {
		fmt.Printf("Trace on\n")
		log.Trace.Printf("tFile - stub [%v]\n", tFile)
	}
	bs, err = files.OpenDescriptionFile(tFile)
	if err != nil {
		log.Error.Fatalf("\n%s", err)
	}
	f, s = files.NewScanner(s, bs)
	pg := lib.NewPage()
	lc := false

	for !lc {
		t, s = lib.NextToken(f, s)
		if t.PTok == token.EOF {
			log.Trace.Printf("[%v]:- Main EOF found.\n", t.PLine)
			break
		}
		// log.Trace.Printf("Kwise parsing token [%v]\n\n", t)
		switch {
		case t.PTok.String() == "IDENT" && str.Compare(str.ToLower(t.PLit), "page") == 0:
			s, st = pg.ProcessPage(f, s, t)
			if st {
				log.Trace.Printf("[%v]:- Error processing page.\n", t.PLine)
				lc = true
			}

		case t.PTok.String() == "IDENT" && str.Compare(str.ToLower(t.PLit), "paragraph") == 0:
			s, st = pg.ProcessParagraph(f, s, t)
			if st {
				log.Trace.Printf("[%v]:- Error processing paragraph\n", t.PLine)
				lc = true
			}

		case t.PTok.String() == "IDENT" && str.Compare(str.ToLower(t.PLit), "sentence") == 0:
			// log.Trace.Printf("[%v]:- Calling Process Sentence", t.PLine)
			s, st = pg.ProcessSentence(f, s, t)
			if st {
				log.Trace.Printf("[%v]:- Error processing sentence\n", t.PLine)
				lc = true
			}
			// log.Trace.Printf("[%v]:- Process Sentence Complete.", t.PLine)

		default:
			/*	if str.Compare(t.PTok.String(), ";") == 0 {
					log.Trace.Printf("[%v]:- In Main skipping Semi-Colon", t.PLine)
				} else {
					log.Trace.Printf("[%v]:-Unhandled token [%v]\n", t.PLine, t)
				} */
			// lc = true
		}

		if lc {
			log.Trace.Println("!!!!!!! ABORTING LC is TRUE !!!!!!!!!!")
			os.Exit(1)
		}
	} // end of for

	if !pg.CheckDirectives() {
		log.Trace.Printf("** Invalid directive ** ")
		os.Exit(1)
	}
	log.Trace.Println(fmt.Sprintf("%s", pg))

	ops := lib.InitEngine(pg)

	for _, o := range ops {
		log.Trace.Printf("[%v]\n", o)
		//	o.PutStr()
	}

	// b =b
	// dFile = dFile

	sequential := 0 // Just dummy this up for now. Unsequenced
	// idx = 0
	df, _ = files.OpenDataFile(dFile)
	b = lib.NewMbuf()
	b.GetBuffer(df)
	f, ds = b.NewKwScanner(ds)
	t.PTok = token.IDENT // Need to prime the pump.
	eofcnt := 0
	for {
		f, ds = b.NextDataToken(df, f, &t, ds)

		if t.PLit == lib.ENDOFBUFFER {
			log.Trace.Printf("MAIN: End of Buffer found\n")
			continue
		}

		if t.PTok == token.EOF {
			if eofcnt == 0 {
				log.Trace.Printf("MAIN: EOF token 1\n")
				eofcnt++
				continue
			} else {
				log.Trace.Printf("MAIN: EOF token FINAL\n")
				break
			}
		}
		eofcnt = 0
		log.Trace.Printf("MAIN TOKENS [%v]", t)
		if t.PTok == token.COMMENT { // This token repurposed to represent a newline
			log.Trace.Printf("End of line reached Clearing status fields\n")
			ops.ClearPatternStatus()
			continue
		}

		if sequential == 0 {
			// log.Trace.Printf("In Sequential [%v]\n", t)
			for o := 0; o < len(ops); o++ {
				// log.Trace.Printf("In Sequential loop index [%d] token [%v]\n", o, t)
				if ops[o].IsParagraphComplete() {
					// log.Trace.Printf("In Sequential - Paragraph is complete [%v]\n", t)
					continue
				}
				if ops[o].MatchCollectSentence(t) {
					ops.ClearPatternStatus()
					// log.Trace.Printf("<******> Matched Sentence# [%d] [%v]\n", o, ops[o].GetSentence())
					// log.Trace.Printf("**** Collecting Sentence# [%d] data directive [%v]\n", o, ops[o].GetPatternDirective(lib.PATTERN))
					// log.Trace.Printf("**** Directives: COLLECT_LINE [%d] COLLECT_NEXT [%d]", lib.COLLECT_LINE, lib.COLLECT_NEXT)
					if int(ops[o].GetPatternDirective(lib.PATTERN)) > 0 {
						if ops[o].GetPatternDirective(lib.PATTERN).IsDirectiveSet(lib.COLLECT_LINE) {
							// log.Trace.Println("########## COLLECT_LINE set", o, ops[o])
							ds = CollectRoutine[lib.COLLECTLINE].CollectResults(df, f, &b, &t, &rs, ds)
							// log.Trace.Println("########## COLLECT_LINE Done", o, ops[o])
							ops[o].CheckParaGraphStatus()
							break
						}
						// log.Trace.Printf("**** CHECKING FOR COLLECT_NEXT")
						if ops[o].GetPatternDirective(lib.PATTERN).IsDirectiveSet(lib.COLLECT_NEXT) {
							// dah =
							// log.Trace.Printf("%%%%%%%% - [%v] len [%d]", dah, len(dah))
							lib.MatchTok(lib.MakeTokenString(lib.ToString(ops[o].GetMatchToken())))
							// log.Trace.Println("########## COLLECT_NEXT set", o, ops[o])
							ds = CollectRoutine[lib.COLLECTNEXT].CollectResults(df, f, &b, &t, &rs, ds)
							// log.Trace.Println("########## COLLECT_NEXT Done", o, ops[o])
							ops[o].CheckParaGraphStatus()
							break
						}
					}
				}

				ops[o].CheckParaGraphStatus()
				// log.Trace.Printf("Loop bottom Sequential index [%d]", o)
			}
			// log.Trace.Printf("Left Sequential [%v]\n", t)
		} else {
			log.Trace.Printf("FOUND MISSING TOKENS [%v]", t)
		}
		// obs.ProcessCompleteParas()
	}

	log.Trace.Printf("OPS Length %d Cap %d", len(ops), cap(ops))
	for k, o := range ops {
		log.Trace.Printf("[%d] [%v]\n", k, o)
		// o.PutStr()
	}

	rs.PutResults()
	log.Info.Printf("** Program completed **\n")
}
