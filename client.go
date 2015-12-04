package nimbus

import (
    "github.com/sorcix/irc"
    "net"
    "fmt"

)

type Client struct {
    Server   string
    Port     string
    Channels []string

    net.Conn
    writer   *IRCWriter
    reader   *IRCReader

    Listeners map[string][]Listener

    Nick     string
    RealName string
    UserName string
    Password string

    // AutoRejoin  bool
    // MaxRetries  int

    // debug bool
}

func NewClient(server string, nick string, config Config) *Client {
    client := &Client{
        Server:   server,
        Port:     config.Port,
        Channels: config.Channels,

        Listeners: make(map[string][]Listener),

        Nick:   nick,
        RealName: config.RealName,
        UserName: config.UserName,
        Password: config.Password,
    }
    return client
}

func (c *Client) Connect(callback func(error)) error {
    conn, err := net.Dial("tcp", c.Server + ":" + c.Port)

    if err != nil {
        return err
    }

    c.Conn = conn
    c.reader = NewIRCReader(conn)
    c.writer = NewIRCWriter(conn)

    if c.Password != "" {
        c.Send(irc.PASS, c.Password)
    }

    c.Send(irc.USER, c.Nick, c.UserName, "0", "*", ":" + c.Nick)
    c.Send(irc.NICK, c.Nick)

    callback(c.register())
    return nil
}

func (c *Client) register() error {
    for {
        message, err := c.reader.Read()

        fmt.Println(message.Command)

        if err != nil {
            return err
        }

        switch message.Command {
            case irc.PING:
                c.Send(irc.PONG, message.Trailing)

            case irc.RPL_WELCOME:
                for _, channel := range c.Channels {
                    c.Send(irc.JOIN, channel)
                }
                return nil
        }
    }
}

func (c *Client) Listen(ch chan<- error) error {
    for {
        message, err := c.reader.Read()

        if err != nil {
            fmt.Println(err)
            return err
        }

        if message.Command == irc.PING {
            c.Send(irc.PONG, message.Trailing)
        }

        fmt.Println(message.Raw)
        c.Emit(message.Command, message)
    }
}
