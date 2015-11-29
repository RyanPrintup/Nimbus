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

    conn     net.Conn
    writer   *irc.Encoder
    reader   *irc.Decoder
    listeners map[string][]Listener

    Nick     string
    RealName string
    UserName string
    Password string

    // AutoRejoin  bool
    // MaxRetries  int

    // debug bool
}

func New(server, nick string, config Config) *Client {
    c := &Client{ Server:   server,
               Port:     config.Port,
               Channels: config.Channels,
               listeners: make(map[string][]Listener),

               Nick:   nick,
               RealName: config.RealName,
               UserName: config.UserName,
               Password: config.Password,
             }
    return c
}

func (c *Client) Connect(callback func(error)) error {
    var conn net.Conn
    var err error

    conn, err = net.Dial("tcp", c.Server + ":" + c.Port)

    if err != nil {
        return err
    }

    c.conn = conn
    c.reader = irc.NewDecoder(conn)
    c.writer = irc.NewEncoder(conn)

    if c.Password != "" {
        c.Send(irc.PASS, c.Password)
    }

    c.Send(irc.NICK, c.Nick)
    c.Send(irc.USER, c.Nick, c.UserName, "0", "*", ":" + c.Nick)

    callback(c.register())
    return nil
}

func (c *Client) register() error {
    for {
        message, err := c.reader.Decode()

        if err != nil {
            return err
        }

        fmt.Println(message.String())

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
        message, err := c.reader.Decode()

        if err != nil {
            fmt.Println(err)
            return err
        }

        if message.Command == irc.PING {
            c.Send(irc.PONG, message.Trailing)
        }

        fmt.Println(message.String())
        c.emit(message.Command, message)
    }
}
