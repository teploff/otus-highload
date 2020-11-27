package domain

type Syncer interface {
	Run()
	Stop()
}
