//go:generate goyacc -o pilang.go -v pilang.output -p "expr" pilang.y

package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"fmt"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	for {
		if _, err := os.Stdout.WriteString("> "); err != nil {
			log.Fatalf("WriteString: %s", err)
		}
		line, err := in.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("ReadBytes: %s", err)
		}

		lex := &exprLex{line: line}
		exprParse(lex)

		ret := lex.ret
		eval(ret, env{})
		fmt.Println(ret)
	}
}
