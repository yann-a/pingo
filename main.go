//go:generate goyacc -o pilang.go -v pilang.output -p "expr" pilang.y

package main

import (
	"bufio"
	"os"
	"fmt"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	lex := &exprLex{reader: in}
	exprParse(lex)

	ret := lex.ret

	eval(ret, &env{})
	fmt.Println(ret)
}
