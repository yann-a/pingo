//go:generate goyacc -o pilang.go -v pilang.output -p "expr" pilang.y
//go:generate goyacc -o lambda/lambda.go -v lambda/lambda.output -p "lambda" lambda/lambda.y

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
	"pingo/lambda"
)

func main() {
	_ = lambda.ShutUp

	// The options of the executable
	outputCode := flag.Bool("showsrc", false, "Output the code before executing it")
	flag.Parse()

	// Parsing
	nonFlagsArgs := flag.Args()
	buffer := os.Stdin
	// If a file is provided we try reading from it
	if len(nonFlagsArgs) > 0 {
		file, err := os.Open(nonFlagsArgs[0])
		if err != nil {
			fmt.Printf("Couldn't read from %s (%s). Reading from stdin\n", nonFlagsArgs[0], err)
		} else {
			buffer = file
		}
	}
	in := bufio.NewReader(buffer)

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
