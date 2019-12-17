package lib

import (
	"fmt"
)

func BuildToken(tok string, b string, e string) string {
	return fmt.Sprintf("%s%s%s", b, tok, e)
}
