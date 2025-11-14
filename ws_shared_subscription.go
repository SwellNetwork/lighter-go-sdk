package lighter

import "sync"

type sharedSubscription struct {
	channel          string
	auth             string
	count            int64
	subscriberByID   map[string]callback
	subscriberFunc   func(channel string, auth string)
	unsubscriberFunc func(channel string, auth string)

	mu sync.RWMutex
}

func newSharedSubscription(
	channel string,
	auth string,
	subscriberFunc, unsubscriberFunc func(channel string, auth string),
) *sharedSubscription {
	return &sharedSubscription{
		channel:          channel,
		auth:             auth,
		count:            0,
		subscriberByID:   make(map[string]callback),
		subscriberFunc:   subscriberFunc,
		unsubscriberFunc: unsubscriberFunc,
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
	u.mu.Unlock()

	if c == 0 {
		u.unsubscriberFunc(u.channel, u.auth)
	}
}

func (u *sharedSubscription) dispatch(data any) {
	u.mu.RLock()
	callbacks := make([]callback, 0, len(u.subscriberByID))
	for _, cb := range u.subscriberByID {
		callbacks = append(callbacks, cb)
	}
	u.mu.RUnlock()

	for _, cb := range callbacks {
		cb(data)
	}
}

func (u *sharedSubscription) clear() {
	u.mu.Lock()
	defer u.mu.Unlock()

	for id := range u.subscriberByID {
		delete(u.subscriberByID, id)
	}
	u.count = 0
	u.unsubscriberFunc(u.channel, u.auth)
}
