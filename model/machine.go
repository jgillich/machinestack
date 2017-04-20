package model

// Machine is a machine, obviously
type Machine struct {
	ID     int64  `jsonapi:"primary,machines"`
	Name   string `jsonapi:"attr,name"`
	Image  string `jsonapi:"attr,image"`
	Driver string `jsonapi:"attr,driver"`
	Owner  int64
	Node   string
}
