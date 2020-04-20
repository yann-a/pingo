package lambda

import (
	"bufio"
	"os"
	"fmt"
)

func Test() {
	in := bufio.NewReader(os.Stdin)

	lex := &lambdaLex{reader: in}
	lambdaParse(lex)

	fmt.Printf("Understood : %v\n", lex.ret)
}