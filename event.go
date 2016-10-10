package nimbus

import "gopkg.in/sorcix/irc.v1"

// Listener represents a function to be added as a listener to client.
type Listener func(*irc.Message)

// Handle is used to sequentially run listeners.
type Handle struct {
	listener Listener
	done     chan bool
}

// Run will call back on the done channel for handle and thereafter run its listener function.
func (h *Handle) Run(msg *irc.Message) {
	h.done <- true
	h.listener(msg)
}

// AddListener adds a listener to an event.
func (c *Client) AddListener(event string, l Listener) {
	c.listeners[event] = append(c.listeners[event], l)
}

// GetListeners returns all listeners for a given event.
func (c *Client) GetListeners(event string) []Listener {
	return c.listeners[event]
}

// Emit will run each listener concurrently. It uses a handle to run the listeners in sequential
// order. This means that listeners that were added first, are ran before listeners added later.
// This should solve a problem where a user wants to run two listeners in parallel, both which will
// acquire some kind of lock. If they aren't run sequentially, and the user want's priority for one
// of the listeners, then that can't be guaranteed since one of the two listeners will be randomly
// chosen for the lock. To solve this, the user should add the listener who should receive priority
// first.
func (c *Client) Emit(event string, msg *irc.Message) {
	for _, listener := range c.listeners[event] {
		h := Handle{listener, make(chan bool)}
		go h.Run(msg)
		<-h.done
	}
}
