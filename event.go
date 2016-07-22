package nimbus

import "gopkg.in/sorcix/irc.v1"

type Listener func(*irc.Message)

type Handle struct {
	listener Listener
	done     chan bool
}

func (h *Handle) Run(msg *irc.Message) {
	h.done <- true
	h.listener(msg)
}

func (c *Client) AddListener(event string, l Listener) {
	c.listeners[event] = append(c.listeners[event], l)
}

func (c *Client) GetListeners(event string) []Listener {
	return c.listeners[event]
}

func (c *Client) Emit(event string, msg *irc.Message) {
	for _, listener := range c.listeners[event] {
		h := Handle{listener, make(chan bool)}
		go h.Run(msg)
		<-h.done
	}
}
