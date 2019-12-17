package kwiselog

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

const logfile = "kwise.log"

//
// There is an awesome article on go logging. Check it out:
// https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html
//

// ResetOutput - Changes the output stream for the logger. Typically used
// to turn on and off trace throughout kwise.
//
func ResetOutput(l *log.Logger, w io.Writer) {
	l.SetOutput(w)
}

func InitLog() {

	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", os.Stderr, ":", err)
	}

	multiHandle := io.MultiWriter(file, os.Stdout)

	Trace = log.New(multiHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(multiHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(multiHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(multiHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
