package nimbus

//import "fmt"
//import "strconv"

type Listener func(*Message)

type Handle struct {
	listener Listener
	done     chan bool
}

func (h *Handle) Run(msg *Message) {
	h.done <- true
	h.listener(msg)
}

func (c *Client) AddListener(event string, l Listener) {
	c.listeners[event] = append(c.listeners[event], l)
}

func (c *Client) GetListeners(event string) []Listener {
	return c.listeners[event]
}

func (c *Client) Emit(event string, msg *Message) {
	for _, listener := range c.listeners[event] {
		h := Handle{listener, make(chan bool)}
		go h.Run(msg)
		<- h.done
	}
}
