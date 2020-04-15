package main

import "fmt"

func (e parallel)    String() (st string) {
  st = "("
  for i,v := range e {
    if i>0 { st = st+" | "}
    st = st + v.String()
  }
  st += ")"
  return
}

func (e send)        String() string      {
  return fmt.Sprintf("(^%v %v)", e.channel, e.value)
}

func (e receiveThen) String() string      {
  return fmt.Sprintf("(%v(%v).%v)", e.channel, e.pattern, e.then)
}

func (e privatize)   String() string      {
  return fmt.Sprintf("((%v).%v)", e.channel, e.then)
}

func (e print)       String() string      {
  return fmt.Sprintf("(print %v; %v)", e.v, e.then)
}

func (e skip)        String() string      {
  return "skip"
}