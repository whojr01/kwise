package files

import (
	"fmt"
	log "kwiselog"
	"testing"
)

func TestOPenAndRead(t *testing.T) {

	var b Mbuf

	log.InitLog()
	fn := "testfile_alpha.txt"
	fooFile, size := OpenDataFile(fn)
	defer fooFile.Close()

	if size > 0 {
		log.Trace.Printf("The file opened successfully [%v] it is of size [%d]", fn, size)
	}
	for b.GetBuffer(fooFile) > 0 {
		// log.Trace.Printf("GETBUFFER RESULTS buf len [%d]\n[%s]\n\n",len(b), b)
		fmt.Printf("IN MAIN -- [%s]\n", string(b))
	}
}
