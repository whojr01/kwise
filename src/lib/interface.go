package lib

import (
	"go/scanner"
)

type CFParser interface {
	NextToken(s scanner.Scanner) (CFToken, scanner.Scanner)
	IsToken(s string) bool
}
