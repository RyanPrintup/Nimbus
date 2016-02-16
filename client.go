package nimbus

import (
	"fmt"
	"net"
	"strings"
)

type Client struct {
	Server   string
	Port     string
	Channels []string

	conn   net.Conn
	writer *IRCWriter
	reader *IRCReader

	listeners map[string][]Listener

	Nick     string
	RealName string
	UserName string
	Password string
	Modes    string

	Quit chan error

	// AutoRejoin  bool
	// MaxRetries  int

	// debug bool
}

func NewClient(server string, nick string, config Config) *Client {
	client := &Client{
		Server:   server,
		Port:     config.Port,
		Channels: config.Channels,

		listeners: make(map[string][]Listener),

		Nick:     nick,
		RealName: config.RealName,
		UserName: config.UserName,
		Password: config.Password,
		Modes: config.Modes,

		Quit: make(chan error),
	}
	return client
}

func (c *Client) Connect(callback func(error)) error {
	conn, err := net.Dial("tcp", c.Server+":"+c.Port)

	if err != nil {
		callback(err)
		return err
	}

	c.conn = conn
	c.reader = NewIRCReader(conn)
	c.writer = NewIRCWriter(conn)

	if c.Password != "" {
		c.Send(PASS, c.Password)
	}

	c.Send(USER, c.Nick, c.UserName, "0", "*", ":"+c.Nick)
	c.Send(NICK, c.Nick)

	callback(nil)

	return nil
}

func (c *Client) Listen() {
	for {
		message, err := c.reader.Read()

		if err != nil {
			fmt.Println(err)
			c.Quit <- err
			return
		}

		switch message.Command {
		case PING:
			c.Send(PONG, message.Trailing)
		case RPL_WELCOME:
			for _, channel := range c.Channels {
				c.Send(MODE, c.Nick, c.Modes)
				c.Send(JOIN, channel)
			}

		}

		fmt.Print(message.Raw)
		c.Emit(message.Command, message)
	}
}

func (c *Client) Send(raw ...string) {
	message, _ := ParseMessage(strings.Join(raw, " "))
	c.writer.Write(message.Bytes())
}

func (c *Client) Say(channel string, text string) {
	c.Send(PRIVMSG, channel, text)
}
