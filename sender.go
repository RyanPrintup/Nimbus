package nimbus

import (
    "github.com/sorcix/irc"
    "strings"
)

func (c *Client) Send(raw ...string) {
    message := irc.ParseMessage(strings.Join(raw, " "))
    c.Writer.Write(message.Bytes())
}

func (c *Client) Say(channel string, text string) {
    c.Send(irc.PRIVMSG, channel, text)
}
