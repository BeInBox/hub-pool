package POOL

type PoolStatus struct {
	Name     string
	CanRead  bool
	CarWrite bool
}

var PoolStatusInit = PoolStatus{"Init", true, false}
