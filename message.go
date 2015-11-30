package nimbus

import (
	"strings"
	"errors"
	"fmt"
)

const (
	ENDLINE = "\r\n"

	// 512 minus CR-LF
	MAX_SIZE = 512 - len(ENDLINE)
)

type Message struct {
	Raw string
	Prefix string
	Command string
	Args []string
	Middle string
	Trailing string
}

func crlfCutsetFunc(r rune) bool {
	return r == '\r' || r == '\n'
}

func ParseMessage(raw string) (*Message, error) {
	message := &Message{} 

	message.Raw = raw
	fmt.Println(raw)

	// Check if message is empty
	if raw = strings.TrimFunc(raw, crlfCutsetFunc); len(raw) < 2 {
		return nil, errors.New("empty message")
	}

	// If prefix is present extract
	if raw[0] == ':' {
		prefix := strings.Split(raw, " ")
		message.Prefix = prefix[0][1 : len(prefix[0])]
		raw = raw[strings.Index(raw, " ") + 1 : len(raw)]
	}

	// Grab command
	message.Command = strings.Split(raw, " ")[0]
	raw = raw[strings.Index(raw, " ") : len(raw)]

	// Check if there is trailing data or just a middle
	fmt.Println(message.Raw)
	if strings.Contains(raw, " :") {
		s := strings.Split(raw, ":")
		message.Middle, message.Trailing = s[0][0 : len(s[0]) - 1], s[1]
	} else {
		message.Middle = raw
	}

	// If middle, extract data to args
	if len(message.Middle) > 0 {
		message.Args = strings.Split(message.Middle, " ")
	}

	// If trailing, append data to args
	if len(message.Trailing) > 0 {
		message.Args = append(message.Args, strings.Split(message.Trailing, " ")...)
	}

	return message, nil
}

func (m *Message) Bytes() []byte {
	return []byte(m.Raw)
}