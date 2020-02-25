package common

type StoreType int

const (
	_ StoreType = iota

	//local
	StoreLocal
	//ceph
	StoreCeph
)
