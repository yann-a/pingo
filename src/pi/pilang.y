%{

package pi

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
  ret Expr
  v Terminal
  num int
  s string
}


%type <ret> expr innerexpression chooseExpression receiveThen
%type <v> pattern value literal arithvalue arith
%token LPAREN RPAREN DOT PIPE COMMA SEMICOLON
%token <num> INT
%token <s> VAR
%token OUTPUT PRINT
%token EQUAL BRA KET REPL
%token PLUS MINUS TIMES DIV

%left PLUS /* + et * sont associatifs à gauche; et * est distributif sur + */
%left MINUS
%left TIMES
%left DIV

%left CAPTUREPLUS

%left COMMA
%left SEMICOLON
%left PIPE

%nonassoc RPAREN



/*********** Parser ***********/
%%

top: expr                                          { exprlex.(*exprLex).ret = $1           }

expr:
    chooseExpression
  | expr PIPE expr                                 {
                                                      switch v := $1.(type) {
                                                        /* We put all piped processes at the same level */
                                                        case Parallel:
                                                          $$ = Parallel(append([]Expr(v), $3))
                                                        default:
                                                          $$ = Parallel{$1, $3}
                                                      }
                                                   }

chooseExpression:
    innerexpression
  | receiveThen PLUS receiveThen                   { $$ = Choose{$1.(ReceiveThen), $3.(ReceiveThen)} }

innerexpression:
    INT                                            { $$ = Skip($1)                         }
  | LPAREN expr RPAREN                             { $$ = $2                               }
  | LPAREN VAR RPAREN innerexpression              { $$ = Privatize{$2, $4}                }
  | REPL VAR pattern DOT innerexpression           { $$ = Repl{$2, $3, $5}                 }
  | receiveThen                                    { $$ = $1                               }
  | OUTPUT VAR value                               { $$ = Send{$2, $3}                     }
  | PRINT value                                    { $$ = Print{$2, Skip(0)}               }
  | PRINT value SEMICOLON innerexpression          { $$ = Print{$2, $4}                    }
  | BRA value EQUAL value KET innerexpression      { $$ = Conditional{$2, true, $4, $6}    }
  | BRA value REPL EQUAL value KET innerexpression { $$ = Conditional{$2, false, $5, $7}   }

receiveThen:
  VAR pattern DOT innerexpression                { $$ = ReceiveThen{$1, $2, $4}          }

pattern: /* for reception */
    VAR                                            { $$ = Variable($1)                     }
  | VAR COMMA VAR                                  { $$ = Pair{Variable($1), Variable($3)} }
  | LPAREN pattern RPAREN                          { $$ = $2                               }

value: /* for sending */
    literal
  | literal COMMA literal                          { $$ = Pair{$1, $3}                     }
  | LPAREN arithvalue RPAREN                       { $$ = $2                               }

/* Arithmetic operations must be between parenthesis */
arithvalue:
    arith
  | arith COMMA arith                              { $$ = Pair{$1, $3}                     }
  | LPAREN arithvalue RPAREN                       { $$ = $2                               }

arith:
    literal
  | arith PLUS arith                               { $$ = Add{$1, $3}                      }
  | arith MINUS arith                              { $$ = Sub{$1, $3}                      }
  | arith TIMES arith                              { $$ = Mul{$1, $3}                      }
  | arith DIV arith                                { $$ = Div{$1, $3}                      }

literal:
    INT                                            { $$ = Constant($1)                     }
  | VAR                                            { $$ = Variable($1)                     }


/**********     Lexer     ***********/
%%

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
  ret Expr
  reader *bufio.Reader
}

func Parse(in *bufio.Reader) Expr {
  lex := &exprLex{reader: in}
  if exprParse(lex) == 1 {
    panic("Parsing error")
  }

  return lex.ret
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
        return SEMICOLON
      case '^':
        return OUTPUT
      case '+':
        return PLUS
      case '-':
        return MINUS
      case '/':
        return DIV
      case '*':
        return TIMES
      case '=':
        return EQUAL
      case '[':
        return BRA
      case ']':
        return KET
      case '!':
        return REPL
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
      case unicode.IsDigit(c):
        add(&b, c)
      default:
        if c != eof {
          x.reader.UnreadRune()
        }

        break L
    }
  }

  yylval.s = b.String()

  if yylval.s == "print" {
    return PRINT
  }

  return VAR
}

// Return the next rune for the lexer.
func (x *exprLex) next() rune {
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

// The parser calls this method on a parse error.
func (x *exprLex) Error(s string) {
  log.Printf("parse error: %s", s)
}
