%{

package lambda

import (
  "bufio"
  "bytes"
  "log"
  "unicode" // To tell apart letters from numbers
  "unicode/utf8"
  "strconv"
  "fmt"
)

%}

%union {
  ret Lambda
  num Lconst
  s Lvar
}


%type <ret> lambda litteral fundefinition application value
%token FUN LPAREN RPAREN
%token PLUS MINUS TIMES DIV
%token GT LT
%token SEMICOLON COLON EQUAL
%token <num> INT
%token <s> VAR

%left PLUS
%left MINUS
%left TIMES
%left DIV

/*********** Parser ***********/
%%

top: lambda                                        { lambdalex.(*lambdaLex).ret = $1       }

lambda:
    value                                          { $$ = $1                               }
  | FUN VAR fundefinition                          { $$ = Lfun{$2, $3}                     }
  | VAR LT MINUS VAR SEMICOLON lambda              { $$ = Read{$1, $4, $6}                 }
  | VAR COLON EQUAL litteral SEMICOLON lambda      { $$ = Write{$1, $4, $6}                }
  | VAR LT MINUS VAR COLON EQUAL litteral SEMICOLON lambda { $$ = Swap{$1, $4, $7, $9}     }

fundefinition:
    MINUS GT lambda                                { $$ = $3                               }
  | VAR fundefinition                              { $$ = Lfun{$1, $2}                     }

value:
    application                                    { $$ = $1                               }
  | value PLUS value                               { $$ = Add{$1, $3}                      }
  | value MINUS value                              { $$ = Sub{$1, $3}                      }
  | value TIMES value                              { $$ = Mult{$1, $3}                     }
  | value DIV value                                { $$ = Div{$1, $3}                      }

application:
    litteral                                       { $$ = $1                               }
  | application litteral                           { $$ = Lapp{$1, $2}                     }

litteral:
    INT                                            { $$ = $1                               }
  | VAR                                            { $$ = $1                               }
  | LPAREN lambda RPAREN                           { $$ = $2                               }


/**********     Lexer     ***********/
%%

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type lambdaLex struct {
  ret Lambda
  reader *bufio.Reader
}

func Parse(in *bufio.Reader) Lambda {
  lex := &lambdaLex{reader: in}
  if lambdaParse(lex) == 1 {
    panic("Parsing error")
  }

  return lex.ret
}

// The parser calls this method to get each new token. This
// implementation returns operators and NUM.
func (x *lambdaLex) Lex(yylval *lambdaSymType) int {
  for {
    c := x.next()
    switch c {
      case eof:
        return eof
      case '(':
        return LPAREN
      case ')':
        return RPAREN
      case '+':
        return PLUS
      case '-':
        return MINUS
      case '>':
        return GT
      case '<':
        return LT
      case ';':
        return SEMICOLON
      case ':':
        return COLON
      case '=':
        return EQUAL
      case '*':
        return TIMES
      case '/':
        return DIV
      case ' ', '\t', '\n', '\r':
      default:
        if unicode.IsLetter(c) {
          return x.string(c, yylval)
        } else if unicode.IsNumber(c) {
          return x.num(c, yylval)
        }

        panic(fmt.Sprintf("unrecognized character %q", c))
      }
  }
}

// Lex a number.
func (x *lambdaLex) num(c rune, yylval *lambdaSymType) int {
  add := func(b *bytes.Buffer, c rune) {
    if _, err := b.WriteRune(c); err != nil {
      log.Fatalf("WriteRune: %s", err)
    }
  }

  var b bytes.Buffer
  add(&b, c)

  L: for {
    c = x.next()
    switch {
      case unicode.IsNumber(c):
        add(&b, c)
      default:
        if c != eof {
          x.reader.UnreadRune()
        }

        break L
    }
  }

  v, err := strconv.Atoi(b.String())
  if err != nil {
    log.Printf("bad number %q", b.String())
    return eof
  }

  yylval.num = Lconst(v)

  return INT
}

// Lex a string.
func (x *lambdaLex) string(c rune, yylval *lambdaSymType) int {
  add := func(b *bytes.Buffer, c rune) {
    if _, err := b.WriteRune(c); err != nil {
      log.Fatalf("WriteRune: %s", err)
    }
  }

  var b bytes.Buffer
  add(&b, c)

  L: for {
    c = x.next()
    switch {
      case unicode.IsLetter(c):
        add(&b, c)
      default:
        if c != eof {
          x.reader.UnreadRune()
        }

        break L
    }
  }

  yylval.s = Lvar(b.String())

  if b.String() == "fun" { return FUN }

  return VAR
}

// Return the next rune for the lexer.
func (x *lambdaLex) next() rune {
  c, size, err := x.reader.ReadRune()

  if c == utf8.RuneError && size == 1 {
    log.Print("invalid utf8")
    return x.next()
  }

  if err != nil {
    return eof
  }

  return c
}

func (x *lambdaLex) getNextRune() rune {
  c, _, err := x.reader.ReadRune()

  if err != nil {
    return eof
  }

  x.reader.UnreadRune()
  return c
}

// The parser calls this method on a parse error.
func (x *lambdaLex) Error(s string) {
  log.Printf("parse error: %s", s)
}
