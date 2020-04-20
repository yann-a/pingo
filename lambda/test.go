package lambda

import (
	"bufio"
	"os"
	"fmt"
	"flag"
)

func Test() {
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

	lex := &lambdaLex{reader: in}
	lambdaParse(lex)

	fmt.Printf("Understood : %v\n", lex.ret)
}