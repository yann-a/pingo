%{

package main

import (
  "bytes"
  "log"
  "fmt"
  "unicode" // différencie les lettres des nombres
  "unicode/utf8"
  "strconv"
)

%}

%union {
  ret expr
  v terminal
  num int
  s string
}


%type <ret> top expr innerexpression
%type <v> pattern value literal
%token LPAREN RPAREN DOT PIPE COMMA COLON
%token <num> INT
%token <s> VAR
%token OUTPUT PRINT

%left COMMA
%left COLON
%left PIPE

%nonassoc RPAREN

%%

top: expr                                          { fmt.Println($1); $$ = $1              }

expr:
    innerexpression
  | expr PIPE expr                                 {
                                                      switch v := $1.(type) {
                                                        case parallel: /* on met tous les process parallèles au même niveau */
                                                          $$ = parallel(append([]expr(v), $3))
                                                        default:
                                                          $$ = parallel{$1, $3}
                                                      }
                                                   }

innerexpression:
    INT                                            { $$ = skip($1)                         }
  | LPAREN expr RPAREN                             { $$ = $2                               }
  | LPAREN VAR RPAREN innerexpression              { $$ = privatize{$2, $4}                }
  | VAR LPAREN pattern RPAREN DOT innerexpression  { $$ = receiveThen{$1, $3, $6}          }
  | OUTPUT VAR value                               { $$ = send{$2, $3}                     }
  | PRINT value                                    { $$ = print{$2, skip(0)}               }
  | PRINT value COLON innerexpression              { $$ = print{$2, $4}                    }

pattern: /* for reception */
    VAR                                            { $$ = variable($1)                     }
  | VAR COMMA VAR                                  { $$ = pair{variable($1), variable($3)} }

value: /* for sending */
    literal
  | literal COMMA literal                          { $$ = pair{$1, $3}         }

literal:
    INT                                            { $$ = constant($1)         }
  | VAR                                            { $$ = variable($1)         }




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
    case ';':
      return COLON
    case '^':
      return OUTPUT
    case ' ', '\t', '\n', '\r':
    default:
      if unicode.IsLetter(c) {
        return x.string(c, yylval)
      } else if unicode.IsNumber(c) {
        return x.num(c, yylval)
      }

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
    switch {
    case unicode.IsNumber(c):
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
    switch {
    case unicode.IsLetter(c):
      add(&b, c)
    default:
      break L
    }
  }
  if c != eof {
    x.peek = c
  }

  yylval.s = b.String()

  if yylval.s == "print" {
    return PRINT
  }

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
