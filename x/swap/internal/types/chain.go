package types

type Chain struct {
	Name   string
	Active bool
}

func NewChain(name string, active bool) Chain {
	return Chain{Name: name, Active: active}
}
