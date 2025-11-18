package lighter

import (
	"sync"
	"time"
)

type subscriptionRetryConfig struct {
	maxAttempts int
	backoff     func(attempt int) time.Duration
	shouldRetry func(err error) bool
}

type sharedSubscription struct {
	channel          string
	auth             string
	count            int64
	subscriberByID   map[string]callback
	subscriberFunc   func(channel string, auth string)
	unsubscriberFunc func(channel string, auth string)
	retryConfig      subscriptionRetryConfig
	retryAttempt     int
	retryTimer       *time.Timer

	mu sync.RWMutex
}

func newSharedSubscription(
	channel string,
	auth string,
	subscriberFunc, unsubscriberFunc func(channel string, auth string),
	retryConfig subscriptionRetryConfig,
) *sharedSubscription {
	return &sharedSubscription{
		channel:          channel,
		auth:             auth,
		count:            0,
		subscriberByID:   make(map[string]callback),
		subscriberFunc:   subscriberFunc,
		unsubscriberFunc: unsubscriberFunc,
		retryConfig:      retryConfig,
	}
}

func (u *sharedSubscription) addSubscriber(id string, cb callback) {
	u.mu.Lock()
	if _, exists := u.subscriberByID[id]; exists {
		u.mu.Unlock()
		return
	}
	u.subscriberByID[id] = cb
	u.count++
	c := u.count
	u.mu.Unlock()

	if c == 1 {
		u.subscriberFunc(u.channel, u.auth)
	}
}

func (u *sharedSubscription) removeSubscriber(id string) {
	u.mu.Lock()
	if _, exists := u.subscriberByID[id]; !exists {
		u.mu.Unlock()
		return
	}
	delete(u.subscriberByID, id)
	c := u.count - 1
	u.count = c
	if c == 0 {
		u.stopRetryLocked()
	}
	u.mu.Unlock()

	if c == 0 {
		u.unsubscriberFunc(u.channel, u.auth)
	}
}

func (u *sharedSubscription) dispatch(data any) {
	callbacks := u.snapshotCallbacks(true)

	for _, cb := range callbacks {
		cb(data)
	}
}

func (u *sharedSubscription) dispatchError(err error) {
	if u.tryScheduleRetry(err) {
		return
	}

	callbacks := u.snapshotCallbacks(false)
	for _, cb := range callbacks {
		cb(err)
	}
}

func (u *sharedSubscription) clear() {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.stopRetryLocked()

	for id := range u.subscriberByID {
		delete(u.subscriberByID, id)
	}
	u.count = 0
	u.unsubscriberFunc(u.channel, u.auth)
}

func (u *sharedSubscription) snapshotCallbacks(resetRetry bool) []callback {
	u.mu.Lock()
	defer u.mu.Unlock()

	if resetRetry {
		u.stopRetryLocked()
	}

	callbacks := make([]callback, 0, len(u.subscriberByID))
	for _, cb := range u.subscriberByID {
		callbacks = append(callbacks, cb)
	}

	return callbacks
}

func (u *sharedSubscription) tryScheduleRetry(err error) bool {
	cfg := u.retryConfig

	if cfg.backoff == nil || cfg.shouldRetry == nil || !cfg.shouldRetry(err) {
		return false
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if len(u.subscriberByID) == 0 {
		return false
	}

	nextAttempt := u.retryAttempt + 1
	if cfg.maxAttempts > 0 && nextAttempt > cfg.maxAttempts {
		return false
	}

	wait := cfg.backoff(nextAttempt)
	if wait <= 0 {
		wait = time.Millisecond
	}

	if u.retryTimer != nil {
		u.retryTimer.Stop()
	}

	u.retryAttempt = nextAttempt
	u.retryTimer = time.AfterFunc(wait, func() {
		u.subscriberFunc(u.channel, u.auth)
	})

	return true
}

func (u *sharedSubscription) stopRetry() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.stopRetryLocked()
}

func (u *sharedSubscription) stopRetryLocked() {
	if u.retryTimer != nil {
		u.retryTimer.Stop()
		u.retryTimer = nil
	}
	u.retryAttempt = 0
}
