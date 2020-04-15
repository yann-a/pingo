package main

// Concrete values
type value interface {
  isValue()
}

type vchannel chan value
func (c vchannel) isValue() { }

type vinteger int
func (c vinteger) isValue() { }

type vpair struct {
  v1 value
  v2 value
}
func (c vpair) isValue() { }


// Environment management
type env struct {
  name variable
  value value
  next *env
}

func (e *env) set_value(x variable, v value) env {
  return env{x, v, e}
}

func (e *env) get_value(x variable) value {
  if (e == nil) {
    channel := make(vchannel)
    e = &env{x, channel, nil} // On ajoute le nouveau channel dans l'espace global en le mettant à la fin de l'environnement

    return channel
  }

  if e.name == x {
    return e.value
  }

  if e.next == nil { // on a atteint la fin de l'environnement
    channel := make(vchannel)
    e.next = &env{x, channel, nil} // On ajoute le nouveau channel dans l'espace global en le mettant à la fin de l'environnement

    return channel
  }

  return e.next.get_value(x)
}
