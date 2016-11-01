package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	splitFlag := flag.Uint("s", 0, "Split (size in KB).")
	joinFlag := flag.Bool("j", false, "Join.")
	flag.Parse()

	if (*splitFlag > 0) && *joinFlag {
		flag.Usage()
		panic(errors.New("Cannot specify both join and split flags."))
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		panic(errors.New("No file specified."))
	}
	if len(flag.Args()) > 1 {
		flag.Usage()
		panic(errors.New("More than one file specified."))
	}

	filename := flag.Args()[0]

	if *splitFlag > 0 {
		split(filename, *splitFlag)
	} else if *joinFlag {
		join(filename);
	} else {
		flag.Usage()
		panic(errors.New("Syntax error"))
	}
}

func split(filename string, size uint) {
	// TODO
}

func join(filename string) {
	dir := filepath.Dir(filename)
	name := filepath.Base(filename)
	p := len(name)
	for i := len(name) - 1; i >= 0; i-- {
		c := name[i]
		if c < '0' || c > '9' {
			p = i + 1
			break
		}
	}
	basename := name[0:p]
	startIndexStr := name[p:]
	startIndex, err := strconv.ParseUint(startIndexStr, 10, 32)
	if err != nil {
		panic(err)
	}
	if (startIndex != 0) && (startIndex != 1) {
		panic(errors.New("Invalid filename"))
	}
	index := startIndex
	formatString := "%s/%s%d"
	if len(startIndexStr) > 1 {
		formatString = fmt.Sprintf("%%s/%%s%%0%dd", len(startIndexStr))
	}
	outputFile, err := os.Create(dir + "/" + basename)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	for {
		path := fmt.Sprintf(formatString, dir, basename, index)
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			panic(err)
		}
		println(path)
		inputFile, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		io.Copy(outputFile, inputFile)
		defer inputFile.Close()

		inputFile.Close()
		index++
	}
	outputFile.Close()
}