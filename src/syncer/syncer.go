package syncer

type Action func(arg interface{}) error

type Service interface {
	Run(stop <-chan struct{})
	Trigger(arg interface{}, auto bool) bool
	GetRecords() []SyncInfo
}

var (
	_ Service = (*OneWorkerSyncer)(nil)
	_ Service = (*WorkequeueSyncer)(nil)
)
