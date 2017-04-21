package model

// Machine is a machine, obviously
type Machine struct {
	ID     int64  `jsonapi:"primary,machines"`
	Name   string `jsonapi:"attr,name" valid:"printableascii,length(4|50),required"`
	Image  string `jsonapi:"attr,image" valid:"printableascii,length(4|50),required"`
	Driver string `jsonapi:"attr,driver" valid:"alpha,required"`
	User  int
	Node   string
}
