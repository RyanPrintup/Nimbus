[![Build Status](https://travis-ci.org/RyanPrintup/Nimbus.svg)](https://travis-ci.org/RyanPrintup/Nimbus)
[![GoDoc](https://godoc.org/github.com/RyanPrintup/Nimbus?status.svg)](https://godoc.org/github.com/RyanPrintup/Nimbus)

# Nimbus
An event-driven IRC client written in Go.

### Overview

This client uses concurrent, sequential listeners to handle IRC events. It's been made to work in a
simple manner, and only handles initial connection to the server. We are in the process of adding
extra functionality, such as timeout detection, reconnect, and TLS configuration. Right now its
primary purpose is to be used as an IRC backbone for a more sophisticated client. 

As a matter of fact, this library was mainly made to serve as the backbone of the [Rain](https://github.com/raindevteam/rain) package.
Rain is somewhat of a framework for creating IRC bots with extra functionality, analogous to taking
a jackhammer to IRC, if Nimbus is a simple hammer.

If you're looking for something simple like Nimbus, but more standalone-ish (or well tested
and feature rich), we recommend using [ircx](https://github.com/nickvanw/ircx), especially if you need TLS.
Giving credit where its due, ircx was a heavy inspiration for Nimbus.

Also many thanks to sorcix for his [irc package](https://github.com/sorcix/irc) 

### Example Usage

This short example gives a good overall look at all the current functionality of Nimbus.

```go
package main

import (
	"fmt"
	"strings"

	"github.com/RyanPrintup/Nimbus"
	"gopkg.in/sorcix/irc.v1"
)

func main() {
	client := nimbus.NewClient("irc.canternet.org", "6667", "NimbusBot", nimbus.Config{
		Channels: []string{"#RainBot"},
		RealName: "Not Nimbus IRC Client",
		UserName: "Fabuloso",
		Debug:    1, // This is the debug level
	})

	client.AddListener(irc.PRIVMSG, func(msg *irc.Message) {
		to := msg.Params[0]

		if to == client.GetNick() {
			to = msg.Params[1]
		}

		if strings.ToLower(msg.Trailing) == "hello nimbusbot" {
			client.Say(to, "Well hello there!")
		}
	})

	if err := client.Connect(); err != nil {
		panic(err)
	}

	client.Listen()

	if err := <-client.Quit(); err != nil {
		fmt.Println(err)
	}
}
```

That should get you a simple bot up and running that responds to "hello nimbusbot."