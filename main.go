//go:generate goyacc -o pilang.go -v pilang.output -p "expr" pilang.y

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
)

func main() {
	// The options of the executable
	outputCode := flag.Bool("showsrc", false, "Output the code before executing it")
	flag.Parse()

	// Parsing
	in := bufio.NewReader(os.Stdin)

	lex := &exprLex{reader: in}
	exprParse(lex)

	ret := lex.ret
	if ret == nil {
		return
	}

	// We print the code if asked
	if *outputCode {
		fmt.Println(ret)
	}

	// We want all goroutines to have finished before killing the process
	var wg sync.WaitGroup

	wg.Add(1)
	// We launch the evaluation in a goroutine
	go eval(ret, &env{}, &wg)

	wg.Wait() // We wait for all goroutines to terminate
}
