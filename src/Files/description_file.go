package files

import (
	"go/scanner"
	"go/token"
	log "kwiselog"
	"os"
)

// This function returns the description file in a
// byte array for processing or returns nil if it
// can't open the file.
//
func OpenDescriptionFile(f string) ([]byte, error) {
	file, err := os.Open(f)
	if err != nil {
		// handle the error here
		return nil, err
	}
	defer file.Close()

	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	// read the file
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return nil, err
	}
	return bs, err
}

// OpenDataFile: Opens the file and returns the file pointer or dies trying.
// OpenDataFile: returns the size of the file along with the file pointer.
func OpenDataFile(fn string) (*os.File, int64) {
	file, err := os.Open(fn)
	if err != nil {
		log.Error.Printf("Error: OpenDataFile file error [%v]", err)
	}
	stat, err := file.Stat()
	if err != nil {
		log.Error.Printf("Error: OpenDataFile error Can't stat file [%s]", fn)
	}
	size := stat.Size()
	return file, size
}

func NewScanner(s scanner.Scanner, bs []byte) (*token.FileSet, scanner.Scanner) {
	fset := token.NewFileSet()
	//	log.Trace.Printf("Buffer len [%d]", len(bs))
	nfile := fset.AddFile("", fset.Base(), len(bs))
	s.Init(nfile, bs, nil /* no error handler */, 0)
	return fset, s
}
