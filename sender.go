package nimbus

import (
    "strings"
)

func (c *Client) Send(raw ...string) {
    message := ParseMessage(strings.Join(raw, " "))
    c.writer.Write(message.Bytes())
}

func (c *Client) Say(channel string, text string) {
    c.Send(PRIVMSG, channel, text)
}
