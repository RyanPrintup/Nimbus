package nimbus

type Listener func(*Message)

func (c *Client) AddListener(event string, l Listener) {
	c.listeners[event] = append(c.listeners[event], l)
}

func (c *Client) Emit(event string, msg *Message) {
	for _, listener := range c.listeners[event] {
		go listener(msg)
	}
}
