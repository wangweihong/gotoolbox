package workqueue

// RateLimitingInterface is an interface that rate limits items being added to the queue.
type RateLimitingInterface interface {
	DelayingInterface

	// AddRateLimited adds an item to the workqueue after the rate limiter says it's ok
	AddRateLimited(item interface{}) // 根据限速器算出对象新的等待插入到的时间(有些限速器会根据失败次数来进行退避,如每次失败会增加2倍延时)，加入等待队列等待加入工作队列

	// Forget indicates that an item is finished being retried.  Doesn't matter whether it's for perm failing
	// or for success, we'll stop the rate limiter from tracking it.  This only clears the `rateLimiter`, you
	// still have to call `Done` on the queue.
	Forget(item interface{}) // 从限速器中移除。下次插入则重新计算等待延时

	// NumRequeues returns back how many times the item was requeued
	NumRequeues(item interface{}) int // 对象已经被限制了多少次（重新尝试加入工作队列多少次)
}

// NewRateLimitingQueue constructs a new workqueue with rateLimited queuing ability
// Remember to call Forget!  If you don't, you may end up tracking failures forever.
func NewRateLimitingQueue(rateLimiter RateLimiter) RateLimitingInterface {
	return &rateLimitingType{
		DelayingInterface: NewDelayingQueue(),
		rateLimiter:       rateLimiter,
	}
}

func NewNamedRateLimitingQueue(rateLimiter RateLimiter, name string) RateLimitingInterface {
	return &rateLimitingType{
		DelayingInterface: NewNamedDelayingQueue(name),
		rateLimiter:       rateLimiter,
	}
}

// rateLimitingType wraps an Interface and provides rateLimited re-enquing.
type rateLimitingType struct {
	DelayingInterface

	rateLimiter RateLimiter // 限速器，用来计算对象等待插入到主工作队列的时间
}

// AddRateLimited AddAfter's the item based on the time when the rate limiter says it's ok
// 根据限速器算出对象新的等待插入到的时间(有些限速器会根据失败次数来进行退避,如每次失败会增加2倍延时)，加入等待队列等待加入工作队列.
func (q *rateLimitingType) AddRateLimited(item interface{}) {
	// 根据限速器算出对象新的等待插入到的时间(有些限速器会根据失败次数来进行退避,如每次失败会增加2倍延时)
	// 如果之前已经在等待队列中等待，则更新等待时间(如果新的等待时间更短)
	q.DelayingInterface.AddAfter(item, q.rateLimiter.When(item))
}

// 对象已经被限制了多少次（重新尝试加入工作队列多少次).
func (q *rateLimitingType) NumRequeues(item interface{}) int {
	return q.rateLimiter.NumRequeues(item)
}

// 从限速器中移除。下次插入则重新计算等待延时.
func (q *rateLimitingType) Forget(item interface{}) {
	q.rateLimiter.Forget(item)
}
