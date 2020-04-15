//go:generate goyacc -o pilang.go -v pilang.output -p "expr" pilang.y

package main

import (
	"bufio"
	"os"
	"fmt"
  "sync"
)


func main() {
	in := bufio.NewReader(os.Stdin)

	lex := &exprLex{reader: in}
	exprParse(lex)

	ret := lex.ret

	fmt.Println(ret)

	// We want all goroutines to have finished before killing the process
	var wg sync.WaitGroup

	wg.Add(1)
	eval(ret, &env{}, &wg)
	wg.Wait() // We wait for all goroutines to terminate

}
