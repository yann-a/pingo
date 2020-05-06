package typing

// Chaines de représentants
type Chain interface {
  isChain()
}

type Repr struct {
  Type Type
}
func (r Repr) isChain() { }

type Link struct {
  R *Chain
}
func (r Link) isChain() { }

// Types
type Type interface {
  isType()
}

type Int struct {}
func (r Int) isType() { }

type Channel struct {
  Inner *Chain
}
func (r Channel) isType() { }

type Pair struct {
  V1 *Chain
  V2 *Chain
}
func (r Pair) isType() { }

type Void struct {}
func (r Void) isType() { }

type Variable struct {}
func (r Variable) isType() { }
