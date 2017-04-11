package model

// Machine is a machine, obviously
type Machine struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Driver string `json:"driver"`
	Owner  string `json:"-"`
	Node   string `json:"-"`
}
