package nimbus

import (
    "net"
    "fmt"
)

type Client struct {
    Server   string
    Port     string
    Channels []string

    conn      net.Conn
    writer   *IRCWriter
    reader   *IRCReader

    listeners map[string][]Listener

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

        listeners: make(map[string][]Listener),

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

    c.conn = conn
    c.reader = NewIRCReader(conn)
    c.writer = NewIRCWriter(conn)

    if c.Password != "" {
        c.Send(PASS, c.Password)
    }

    c.Send(USER, c.Nick, c.UserName, "0", "*", ":" + c.Nick)
    c.Send(NICK, c.Nick)

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
            case PING:
                c.Send(PONG, message.Trailing)

            case RPL_WELCOME:
                for _, channel := range c.Channels {
                    c.Send(JOIN, channel)
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

        if message.Command == PING {
            c.Send(PONG, message.Trailing)
        }

        fmt.Println(message.Raw)
        c.Emit(message.Command, message)
    }
}
