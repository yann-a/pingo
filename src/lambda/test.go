package lambda

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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
			defer file.Close()
			buffer = file
		}
	}
	in := bufio.NewReader(buffer)

	lex := &lambdaLex{reader: in}
	if lambdaParse(lex) == 1 {
		panic("Parsing error")
	}

	fmt.Printf("Understood : %v\n", lex.ret)
}

func GetParsedInput() Lambda {
	nonFlagsArgs := flag.Args()
	buffer := os.Stdin
	// If a file is provided we try reading from it
	if len(nonFlagsArgs) > 0 {
		file, err := os.Open(nonFlagsArgs[0])
		if err != nil {
			fmt.Printf("Couldn't read from %s (%s). Reading from stdin\n", nonFlagsArgs[0], err)
		} else {
			defer file.Close()
			buffer = file
		}
	}
	in := bufio.NewReader(buffer)

	lex := &lambdaLex{reader: in}
	if lambdaParse(lex) == 1 {
		panic("Parsing error")
	}

	return lex.ret
}
