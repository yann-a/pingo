//go:generate goyacc -o pilang.go -v pilang.output -p "expr" pilang.y

package main

import (
	"bufio"
	"os"
	"fmt"
  "sync"
)

func waitWorker(wg *sync.WaitGroup) {

}

func main() {
	in := bufio.NewReader(os.Stdin)

	lex := &exprLex{reader: in}
	exprParse(lex)

	ret := lex.ret

	fmt.Println(ret)

	// We want all goroutines to have finished before killing the process
	var wg sync.WaitGroup

	wg.Add(1)
	go eval(ret, &env{}, &wg)

	go waitWorker(&wg) // We wait for all goroutines to terminate
										 // it is in a different goroutine in order to prevent "all goroutines are asleep" errors when there is only one parallel process
}
