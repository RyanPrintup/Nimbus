package nimbus

import (
    "github.com/sorcix/irc"
)

type Listener func(*irc.Message)

func (c *Client) AddListener(e string, l Listener) {
    c.Listeners[e] = append(c.Listeners[e], l)
}

func (c *Client) emit(e string, msg *irc.Message) {
    for _, listener := range c.Listeners[e] {
        go listener(msg)
    }
}
