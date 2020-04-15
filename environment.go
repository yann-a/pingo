package main

type value interface {
  isValue()
}

type channel chan value
func (c channel) isValue() { }

type integer int
func (c integer) isValue() { }

type env struct {

}
