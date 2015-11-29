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

    conn   net.Conn
    writer *irc.Encoder
    reader *irc.Decoder

    Nick     string
    RealName string
    UserName string
    Password string

    // AutoConnect bool
    // AutoRejoin  bool
    // MaxRetries  int

    // debug bool
}

func New(server, nick string, config Config) *Client {
    c := &Client{ Server:   server,
               Port:     config.Port,
               Channels: config.Channels,

               Nick:   nick,
               RealName: config.RealName,
               UserName: config.UserName,
               Password: config.Password,
             }
    return c
}

func (c *Client) Connect() error {
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
        c.writer.Encode(&irc.Message{
            Command: irc.PASS,
            Params:  []string{c.Password},
        })
    }

    c.writer.Encode(&irc.Message{
        Command: irc.NICK,
        Params:  []string{c.Nick},
    })

    c.writer.Encode(&irc.Message{
        Command: irc.USER,
        Params:  []string{c.Nick, c.UserName, "0", "*", ":" + c.Nick},
    })

    ch := make(chan error)
    go c.listen(ch)

    return <- ch
}

func (c *Client) listen(ch chan<- error) {
    for {
        message, err := c.reader.Decode()

        if err != nil {
            ch <- err
            return
        }

        switch message.Command {
        case irc.PING:
            c.writer.Encode(&irc.Message{
                Command: irc.PONG,
                Params: []string{message.Trailing},
            })

        case irc.RPL_WELCOME:
            ch <- nil
            for _, channel := range c.Channels {
                c.writer.Encode(&irc.Message{
                    Command: irc.JOIN,
                    Params: []string{channel},
                })
            }
        }

        fmt.Println(message.String())
    }
}
