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
  ret expr
  num int
  s string
}

%type <ret> expr
%token LPAREN RPAREN DOT PIPE COMMA
%token <num> INT
%token <s> VAR

%%

top: expr         { fmt.Println($1)     }

expr:
    INT           { $$ = constant($1)   }
  | VAR           { $$ = channel($1)    }


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
    case '(':
      return LPAREN
    case ')':
      return RPAREN
    case '.':
      return DOT
    case '|':
      return PIPE
    case ',':
      return COMMA
    case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
      return x.num(c, yylval)
    case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
      return x.string(c, yylval)
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

  yylval.num = v

  return INT
}

// Lex a string.
func (x *exprLex) string(c rune, yylval *exprSymType) int {
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
    case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
      add(&b, c)
    default:
      break L
    }
  }
  if c != eof {
    x.peek = c
  }

  yylval.s = b.String()

  return VAR
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
