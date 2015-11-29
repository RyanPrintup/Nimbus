package nimbus

import (
    "github.com/sorcix/irc"
)

type Listener func(*irc.Message)

func (c *Client) AddListener(e string, l Listener) {
    c.listeners[e] = append(c.listeners[e], l)
}

func (c *Client) emit(e string, msg *irc.Message) {
    for _, listener := range c.listeners[e] {
        go listener(msg)
    }
}
