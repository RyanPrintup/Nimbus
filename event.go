package nimbus

type Listener func(*Message)

func (c *Client) AddListener(e string, l Listener) {
    c.Listeners[e] = append(c.Listeners[e], l)
}

func (c *Client) Emit(e string, msg *Message) {
    for _, listener := range c.Listeners[e] {
        go listener(msg)
    }
}
