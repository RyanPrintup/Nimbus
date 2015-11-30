package nimbus

type Listener func(*Message)

func (c *Client) AddListener(e string, l Listener) {
    c.Listeners[e] = append(c.Listeners[e], l)
}

<<<<<<< HEAD
func (c *Client) emit(e string, msg *irc.Message) {
=======
func (c *Client) Emit(e string, msg *Message) {
>>>>>>> 71b2d405597943afa6f49fceaf6a718d6eb8e99f
    for _, listener := range c.Listeners[e] {
        go listener(msg)
    }
}
