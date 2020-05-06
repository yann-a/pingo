package typing

import "pingo/src/pi"

func evalRepr(chain *Chain) *Chain {
  switch v := (*chain).(type) {
  case Repr:
    return chain
  case Link:
    repr := evalRepr(v.R)
    v.R = repr // on pointe directement vers le réprésentant

    return repr
  default:
    panic("Unknown chain type")
  }
}

func unify(t1, t2 *Chain) {
  r1 := evalRepr(t1) // On détermine les représentants
  r2 := evalRepr(t2)

  type1 := (*r1).(Repr).Type
  type2 := (*r2).(Repr).Type

  // si t1 est une variable
  _, ok := type1.(Variable)
  if ok {
    *r1 = Link{r2} // on fait pointer r1 vers r2
    return
  }

  // si t2 est une variable
  _, ok = type2.(Variable)
  if ok {
    *r2 = Link{r1} // on fait pointer r2 vers r1
    return
  }

  switch v1 := type1.(type) {
  case Pair:
    p2 := type2.(Pair) // t2 doit être une paire

    unify(v1.V1, p2.V1)

  case Void:
    _, ok = type2.(Void)
    if !ok {
      panic("Unification failed")
    }

  case Channel:
    ch2 := type2.(Channel)
    unify(v1.Inner, ch2.Inner)

  case Int:
    _, ok = type2.(Int)
    if !ok {
      panic("Unification failed")
    }
  }
}

func typeTerminal(terminal pi.Terminal, env *env) *Chain {
  switch v := terminal.(type) {
  case pi.Constant:
    return createRepr(Int{})

  case pi.Variable:
    return env.get_type(v)

  case pi.Pair:
    return createRepr(Pair{typeTerminal(v.V1, env), typeTerminal(v.V2, env)})
  case pi.Nothing:
    return createRepr(Void{})

  case pi.Add:
    unify(typeTerminal(v.V1, env), createRepr(Int{}))
    unify(typeTerminal(v.V2, env), createRepr(Int{}))

    return createRepr(Int{})

  case pi.Sub:
    unify(typeTerminal(v.V1, env), createRepr(Int{}))
    unify(typeTerminal(v.V2, env), createRepr(Int{}))

    return createRepr(Int{})

  case pi.Mul:
    unify(typeTerminal(v.V1, env), createRepr(Int{}))
    unify(typeTerminal(v.V2, env), createRepr(Int{}))

    return createRepr(Int{})

  case pi.Div:
    unify(typeTerminal(v.V1, env), createRepr(Int{}))
    unify(typeTerminal(v.V2, env), createRepr(Int{}))

    return createRepr(Int{})

  default:
    panic("Unknown expr type")
  }
}

func TypeExpression(expr pi.Expr, env *env) {
  switch v := expr.(type) {
  case pi.Skip:
    return

  case pi.Parallel:
    for _, task := range v {
      TypeExpression(task, env)
    }

  case pi.Send:
    chantype := typeTerminal(pi.Variable(v.Channel), env)
    senttype := typeTerminal(v.Value, env)

    unify(chantype, createRepr(Channel{senttype}))

  case pi.ReceiveThen:
    subenv, argtype := env.type_from_pattern(v.Pattern)
    unify(env.get_type(pi.Variable(v.Channel)), createRepr(Channel{argtype}))

    TypeExpression(v.Then, subenv)

  case pi.Repl:
    subenv, argtype := env.type_from_pattern(v.Pattern)
    unify(env.get_type(pi.Variable(v.Channel)), createRepr(Channel{argtype}))

    TypeExpression(v.Then, subenv)

  case pi.Print:
    unify(typeTerminal(v.V, env), createRepr(Int{}))
    TypeExpression(v.Then, env)

  case pi.Privatize:
    newenv := env.privatize(pi.Variable(v.Channel))
    TypeExpression(v.Then, newenv)

  case pi.Choose:
    TypeExpression(v.E, env)
    TypeExpression(v.F, env)

  default:
    panic("Unknown expr type")
  }
}
