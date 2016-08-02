package nimbus

import (
	"net"
	"strings"

	"gopkg.in/sorcix/irc.v1"
)

// Client contains all the information needed to properly run it's IRC connection
type Client struct {
	Server   string
	Port     string
	channels []string

	conn   net.Conn
	writer *irc.Encoder
	reader *irc.Decoder

	listeners map[string][]Listener

	Nick     string
	RealName string
	UserName string
	Password string
	Modes    string

	quit chan error

	debug int
}

// NewClient returns a new client. Only server, port and nick are required. The remaiming options
// may be set to their appropriate zero value in the config.
func NewClient(server string, port string, nick string, config Config) *Client {
	client := &Client{
		Server:   server,
		Port:     port,
		channels: config.Channels,

		listeners: make(map[string][]Listener),

		Nick:     nick,
		RealName: config.RealName,
		UserName: config.UserName,
		Password: config.Password,
		Modes:    config.Modes,

		quit:  make(chan error, 1),
		debug: config.Debug,
	}
	return client
}

// GetNick returns the client's currrent nick
func (c *Client) GetNick() string {
	return c.Nick
}

// SetNick will send an IRC Nick command to the IRC server. The client's nick will be changed
// regardless if the nick is not valid or in use by someone else on the server.
func (c *Client) SetNick(nick string) {
	c.Send(irc.NICK, nick)
	c.Nick = nick
}

// GetChannels will return the channels specified in the config.
func (c *Client) GetChannels() []string {
	return c.channels
}

// Connect will create a connection to the specified IRC server in the config. It will first send
// the password in the configuration, followed by the client's nick and user. If no username was
// specified in the config, the bot's nick will be used instead. If no realname was specified in the
// the config, the realname will default to "Nimbus IRC Client".
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

	if c.UserName == "" {
		c.UserName = c.Nick
	}

	if c.RealName == "" {
		c.RealName = "Nimbus IRC Client"
	}

	c.Send(irc.USER, c.UserName, "0", "*", ":"+c.RealName)
	c.Send(irc.NICK, c.Nick)

	return nil
}

// Quit returns the client's quit chan, which returns an error after disconnecting from the IRC
// server.
func (c *Client) Quit() chan error {
	return c.quit
}

// Listen will start the listen loop. The listen loop first reads a message from the IRC server, it
// then checks one of two things:
//
// A. Is a PING from the IRC server - Then the client will pong.
//
// B. Is a RPL_WELCOME - We'll send our mode and join all channels found in the config.
//
// Also, if the debug level is 1, it will print the received IRC message. If debug level is 2, it
// will also print the parsed params and trailing (mostly used for debugging purposes).
func (c *Client) Listen() {
	go c.listenLoop()
}

func (c *Client) listenLoop() {
	for {
		message, err := c.reader.Decode()

		if err != nil {
			c.quit <- err
			break
		}

		switch message.Command {
		case irc.PING:
			c.Send(irc.PONG, message.Trailing)
		case irc.RPL_WELCOME:
			for _, channel := range c.channels {
				c.Send(irc.MODE, c.Nick, c.Modes)
				c.Send(irc.JOIN, channel)
			}
		}

		c.Lprintln(1, message.String())

		c.Lprintf(2, "Trailing: %s\n", message.Trailing)
		c.Lprintf(2, "Params: %s\n", message.Params)

		c.Emit(message.Command, message)
	}
}

// Send sends a raw IRC message. Give your command as an argument separated list
func (c *Client) Send(raw ...string) {
	message := irc.ParseMessage(strings.Join(raw, " "))
	c.writer.Encode(message)
}

// Say is a convenience command. Specify a receiver (nick or channel) as the first argument and
// message as the second.
func (c *Client) Say(to string, text string) {
	c.Send(irc.PRIVMSG, to, ":"+text)
}
