package typing

func evalRepr(chain *Chain) *Chain {
  switch v := (*chain).(type) {
  case Repr:
    return chain;
  case Link:
    repr := evalRepr(v.R)
    *chain = Link{repr} // on pointe directement vers le réprésentant

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
