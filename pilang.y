%{

package main

import (
	"bytes"
	"fmt"
	"log"
	"unicode/utf8"
	"strconv"
)

%}

%union {
	num expr
}

%type	<num>	expr expr1 expr2 expr3

%token '+' '-' '*' '/' '(' ')'

%token	<num>	NUM

%%

top:
	expr
	{
		fmt.Println($1)
	}

expr:
	expr1
|	'+' expr
	{
		$$ = $2
	}
|	'-' expr
	{
		$$ = - $2
	}

expr1:
	expr2
|	expr1 '+' expr2
	{
		$$ = $1 + $3
	}
|	expr1 '-' expr2
	{
		$$ = $1 - $3
	}

expr2:
	expr3
|	expr2 '*' expr3
	{
		$$ = $1 * $3
	}
|	expr2 '/' expr3
	{
		$$ = $1 / $3
	}

expr3:
	NUM
|	'(' expr ')'
	{
		$$ = $2
	}


%%

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
	line []byte
	peek rune
}

// The parser calls this method to get each new token. This
// implementation returns operators and NUM.
func (x *exprLex) Lex(yylval *exprSymType) int {
	for {
		c := x.next()
		switch c {
		case eof:
			return eof
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return x.num(c, yylval)
		case '+', '-', '*', '/', '(', ')':
			return int(c)

		// Recognize Unicode multiplication and division
		// symbols, returning what the parser expects.
		case '×':
			return '*'
		case '÷':
			return '/'

		case ' ', '\t', '\n', '\r':
		default:
			log.Printf("unrecognized character %q", c)
		}
	}
}

// Lex a number.
func (x *exprLex) num(c rune, yylval *exprSymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			log.Fatalf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	add(&b, c)
	L: for {
		c = x.next()
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e', 'E':
			add(&b, c)
		default:
			break L
		}
	}
	if c != eof {
		x.peek = c
	}

	v, err := strconv.Atoi(b.String())
	if err != nil {
		log.Printf("bad number %q", b.String())
		return eof
	}

	yylval.num = expr(v)

	return NUM
}

// Return the next rune for the lexer.
func (x *exprLex) next() rune {
	if x.peek != eof {
		r := x.peek
		x.peek = eof
		return r
	}
	if len(x.line) == 0 {
		return eof
	}
	c, size := utf8.DecodeRune(x.line)
	x.line = x.line[size:]
	if c == utf8.RuneError && size == 1 {
		log.Print("invalid utf8")
		return x.next()
	}
	return c
}

// The parser calls this method on a parse error.
func (x *exprLex) Error(s string) {
	log.Printf("parse error: %s", s)
}
