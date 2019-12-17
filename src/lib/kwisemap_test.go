package lib

import (
	"fmt"
	"go/scanner"
	"go/token"
	log "kwiselog"
	str "strings"
	"testing"
)

func TestSentence(t *testing.T) {
	var (
		src = []byte("This is a test sentence along with attributes :collect-next and concat")
		s   scanner.Scanner
		fnd = false
		tok CFToken
	)

	log.InitLog()
	kw := KwiseParse{}
	fset := token.NewFileSet()
	nfile := fset.AddFile("", fset.Base(), len(src))
	s.Init(nfile, src, nil, scanner.ScanComments)
	for {
		tok, s = NextToken(kw, s)
		if tok.PTok == token.EOF {
			break
		}
		if tok.PTok.String() == "IDENT" && IsToken(kw, tok.PLit) && str.Compare(tok.PLit, "sentence") == 0 {
			fnd = true
			break
		}
	}
	if !fnd {
		t.Errorf("\nFailed to find token [sentence] in string")
	}
}

func TestMap(t *testing.T) {
	var (
		s     scanner.Scanner
		fnd   int
		nfnd  int
		total int
		tok   CFToken
		src   []byte
		st    = false
	)

	for k, _ := range kwiseMap {
		src = append(src, []byte(k)...)
		src = append(src, []byte(" ")...)
	}

	// Verify bad tokens are not recognized as good.
	src = append(src, []byte(":collect-seashells ")...) // Invalid token
	src = append(src, []byte("letter ")...)             // Invalid token

	total = len(kwiseMap)
	kw := KwiseParse{}
	fset := token.NewFileSet()
	nfile := fset.AddFile("", fset.Base(), len(src))
	s.Init(nfile, src, nil, scanner.ScanComments)
	for {
		tok, s = NextToken(kw, s)
		if tok.PTok == token.EOF {
			break
		}
		if tok.PTok.String() == "IDENT" && str.Compare(tok.PLit, "comma") != 0 && str.Compare(tok.PLit, "tab") != 0 {
			st = IsToken(kw, tok.PLit)
			if st {
				fnd++
			} else {
				t.Logf("\n** Token Failed [%v]", tok)
				nfnd++
			}
		}

		if tok.PTok.String() == ":" {
			s, st = BuildDirective(s, kw, &tok)
			if st {
				fnd++
			} else {
				t.Logf("\n** Directive Failed [%v]", tok)
				nfnd++
			}
		}
	}
	t.Logf("\n")
	if fnd != total {
		t.Errorf("\n\nfrom string [%v]\nFound [%v] tokens out of [%v]", string(src), fnd, total)
		fmt.Printf("Src [%v]", string(src))
	}
	t.Logf("\n\nTotal totkens processed: [%d]\nSuccesses: [%d]\nFailures: [%d]\nExpected failures: [2]\n\t:seashells\n\tletter\n", total, fnd, nfnd)
}
