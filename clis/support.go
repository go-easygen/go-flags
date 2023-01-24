////////////////////////////////////////////////////////////////////////////
// Program: support.go
// Purpose: cli boilerplate runtime support functions
// Authors: Tong Sun (c) 2015-2019, All rights reserved
////////////////////////////////////////////////////////////////////////////

package clis

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/labstack/gommon/color"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

// The OptsT type defines all the configurable options from cli.
type OptsT struct {
	Verbose int
}

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var (
	progname = "wireframe"
	Opts     OptsT
)

func Setup(p string, v int) {
	progname = p
	Opts.Verbose = v
}

// SUPPORT-FUNCTIONS
//==========================================================================
// support functions

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Basename returns the file name without extension.
func Basename(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}

// IsExist checks if the given file exist
func IsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// ReadInput reads from the given file name (or "-" for stdin)
// and abort the program if error occurs.
func ReadInput(fileName string) []byte {
	var data []byte
	var err error
	if fileName == "-" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = ioutil.ReadFile(fileName)
	}
	AbortOn("ReadInput", err)
	return data
}

// GetInputStream will return the file reading handler to the given file,
// or stdin for "-".
// It'll abort the program if error occurs.
func GetInputStream(fileName string) io.ReadCloser {
	if fileName == "-" {
		return os.Stdin
	}
	// open input file instead
	file, err := os.Open(fileName)
	AbortOn("GetInputStream", err)
	return file
}

// GetOutputStream will return the file writing handler to the given file,
// or stdout for "-".
// It'll abort the program if error occurs.
func GetOutputStream(fileName string) io.WriteCloser {
	if fileName == "-" {
		return os.Stdout
	}
	// open output file instead
	file, err := os.Create(fileName)
	AbortOn("GetOutputStream", err)
	return file
}

// Warning will print the given string as a Warning
func Warning(m string) {
	fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", progname, color.Yellow("Warning"), m)
}

// WarnOn will print the error message as a Warning, if applicable,
// and retur true if so.
// For a suggested format of "ActionName, step name", the output would be
//
//   [progname] Warning: ActionName, step name, sql: Rows are closed
func WarnOn(errCase string, e error) bool {
	if e != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s: %s, %v\n",
			color.White(progname), color.Yellow("Warning"), errCase, e)
		return true
	}
	return false
}

// AbortOn will quit on the anticipated error gracefully without stack trace
func AbortOn(errCase string, e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "[%s] %s: %s, %v\n",
			color.White(progname), color.Red("Error"), errCase, e)
		os.Exit(1)
	}
}

// Verbose will print info to stderr according to the verbose level setting
func Verbose(levelSet int, format string, args ...interface{}) {
	if Opts.Verbose >= levelSet {
		fmt.Fprintf(os.Stderr, "[%s] ", color.White(progname))
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
