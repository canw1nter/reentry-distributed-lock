package dlock

type DistributedLock interface {
	Lock()
	Unlock()
}
