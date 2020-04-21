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


%type <ret> lambda litteral fundefinition application
%token FUN ARROW LPAREN RPAREN
%token <num> INT
%token <s> VAR


/*********** Parser ***********/
%%

top: lambda                                        { lambdalex.(*lambdaLex).ret = $1       }

lambda:
    application                                    { $$ = $1                               }
  | FUN VAR fundefinition                          { $$ = Lfun{$2, $3}                     }

fundefinition:
    ARROW lambda                                   { $$ = $2                               }
  | VAR fundefinition                              { $$ = Lfun{$1, $2}                     }

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
      case '-':
        c = x.getNextRune();
        if c == '>'{
          x.next()
          return ARROW
        }
        // return MINUS
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
