package nimbus

type Listener func(*Message)

func (c *Client) AddListener(event string, l Listener) {
    c.Listeners[event] = append(c.Listeners[event], l)
}

func (c *Client) Emit(event string, msg *Message) {
    for _, listener := range c.Listeners[event] {
        go listener(msg)
    }
}
