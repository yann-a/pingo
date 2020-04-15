package main


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


type env struct {
  name variable
  value value
  next *env
}
