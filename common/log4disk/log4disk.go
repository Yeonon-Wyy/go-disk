package log4disk

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	verbose *log.Logger
	debug   *log.Logger
	info    *log.Logger
	warn    *log.Logger
	error   *log.Logger

	Writer io.Writer = os.Stdout
)

func init() {
	verbose = log.New(Writer, "[go-disk] [VERBOSE]", log.Ldate|log.Lmicroseconds|log.Llongfile)
	debug = log.New(Writer, "[go-disk] [DEBUG]", log.Ldate|log.Lmicroseconds|log.Llongfile)
	info = log.New(Writer, "[go-disk] [INFO]", log.Ldate|log.Lmicroseconds|log.Llongfile)
	warn = log.New(Writer, "[go-disk] [WARN]", log.Ldate|log.Lmicroseconds|log.Llongfile)
	error = log.New(Writer, "[go-disk] [ERROR]", log.Ldate|log.Lmicroseconds|log.Llongfile)

	//Verbose don't print to anywhere on default setting
	//We can use SetWriter to control
	SetWriterForVerbose(ioutil.Discard)
}

func SetWriterForVerbose(writer io.Writer) {
	verbose.SetOutput(writer)
}
func SetWriterForInfo(writer io.Writer) {
	info.SetOutput(writer)
}
func SetWriterForWarn(writer io.Writer) {
	warn.SetOutput(writer)
}
func SetWriterForDebug(writer io.Writer) {
	debug.SetOutput(writer)
}
func SetWriterForError(writer io.Writer) {
	error.SetOutput(writer)
}

func CloseAllLogger() {
	SetWriterForVerbose(ioutil.Discard)
	SetWriterForInfo(ioutil.Discard)
	SetWriterForDebug(ioutil.Discard)
	SetWriterForWarn(ioutil.Discard)
	SetWriterForError(ioutil.Discard)

}

func OpenAllLogger(writer io.Writer) {
	SetWriterForInfo(writer)
	SetWriterForDebug(writer)
	SetWriterForWarn(writer)
	SetWriterForError(writer)

	//Verbose need handled separately
}

func V(format string, args ...interface{}) {
	verbose.Output(2, fmt.Sprintf(format, args...))
}

func I(format string, args ...interface{}) {
	info.Output(2, fmt.Sprintf(format, args...))
}

func D(format string, args ...interface{}) {
	debug.Output(2, fmt.Sprintf(format, args...))
}

func W(format string, args ...interface{}) {
	warn.Output(2, fmt.Sprintf(format, args...))
}


func E(format string, args ...interface{}) {
	error.Output(2, fmt.Sprintf(format, args...))
}
