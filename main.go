//go:generate goyacc -o pilang.go -v pilang.output -p "expr" pilang.y

package main

import (
	"bufio"
	"os"
	"fmt"
	"sync"
	"flag"
)

func main() {
	outputCode := flag.Bool("showsrc", false, "Output the code before executing it")
	flag.Parse()

	in := bufio.NewReader(os.Stdin)

	lex := &exprLex{reader: in}
	exprParse(lex)

	ret := lex.ret
	if ret == nil {
		return
	}

	if *outputCode {
		fmt.Println(ret)
	} else {
		fmt.Println("no output")
	}

	// We want all goroutines to have finished before killing the process
	var wg sync.WaitGroup

	wg.Add(1)
	go eval(ret, &env{}, &wg)

	wg.Wait() // We wait for all goroutines to terminate
}
