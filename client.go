package nimbus

import (
	"fmt"
	"net"
	"strings"

	"gopkg.in/sorcix/irc.v1"
)

type Client struct {
	Server   string
	Port     string
	channels []string

	conn   net.Conn
	writer *irc.Encoder
	reader *irc.Decoder

	listeners map[string][]Listener

	nick     string
	RealName string
	UserName string
	Password string
	Modes    string

	quit chan error

	// AutoRejoin  bool
	// MaxRetries  int

	// debug bool
}

func NewClient(server string, nick string, config Config) *Client {
	client := &Client{
		Server:   server,
		Port:     config.Port,
		channels: config.Channels,

		listeners: make(map[string][]Listener),

		nick:     nick,
		RealName: config.RealName,
		UserName: config.UserName,
		Password: config.Password,
		Modes:    config.Modes,

		quit: make(chan error),
	}
	return client
}

func (c *Client) GetNick() string {
	return c.nick
}

func (c *Client) GetChannels() []string {
	return c.channels
}

func (c *Client) Connect() (err error) {
	c.conn, err = net.Dial("tcp", c.Server+":"+c.Port)

	if err != nil {
		return err
	}

	c.reader = irc.NewDecoder(c.conn)
	c.writer = irc.NewEncoder(c.conn)

	if c.Password != "" {
		c.Send(irc.PASS, c.Password)
	}

	c.Send(irc.USER, c.nick, c.UserName, "0", "*", ":"+c.nick)
	c.Send(irc.NICK, c.nick)

	return nil
}

func (c *Client) Quit() chan error {
	return c.quit
}

func (c *Client) Listen() {
	for {
		message, err := c.reader.Decode()

		if err != nil {
			fmt.Println(err)
			c.quit <- err
			return
		}

		switch message.Command {
		case irc.PING:
			c.Send(irc.PONG, message.Trailing)
		case irc.RPL_WELCOME:
			for _, channel := range c.channels {
				c.Send(irc.MODE, c.nick, c.Modes)
				c.Send(irc.JOIN, channel)
			}
		}

		fmt.Println(message.String())
		fmt.Println(message.Trailing)
		fmt.Println(message.Params)
		c.Emit(message.Command, message)
	}
}

func (c *Client) Send(raw ...string) {
	message := irc.ParseMessage(strings.Join(raw, " "))
	c.writer.Encode(message)
}

func (c *Client) Say(channel string, text string) {
	c.Send(irc.PRIVMSG, channel, ":"+text)
}
