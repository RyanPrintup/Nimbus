package nimbus

import (
	"strings"
	"errors"
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

func splitInTwo(s, sep string) (string, string) {
    x := strings.SplitN(s, sep, 2)
    return x[0], x[1]
}

func crlfCutsetFunc(r rune) bool {
	return r == '\r' || r == '\n'
}

func ParseMessage(raw string) (*Message, error) {
	message := &Message{}

	message.Raw = raw

	// Check if message is empty
	if raw = strings.TrimFunc(raw, crlfCutsetFunc); len(raw) < 2 {
		return nil, errors.New("empty message")
	}

	// If prefix is present extract
	if raw[0] == ':' {
		prefix := strings.SplitN(raw, " ", 2)
		message.Prefix, raw = prefix[0][1:], prefix[1]
	}

	// Grab command
	message.Command, _ = splitInTwo(raw, " ")
	raw = strings.SplitAfterN(raw, message.Command, 2)[1]

	// Check if there is trailing data or just a middle
	if strings.Contains(raw, " :") {
		message.Middle, message.Trailing = splitInTwo(raw, " :")
	} else {
		message.Middle = raw
	}

	// Middle will have some leading whitespace due to command extraction
	message.Middle = strings.TrimLeft(message.Middle, " ")

	// If middle, split into args
	if message.Middle != "" {
		message.Args = strings.Split(message.Middle, " ")
	}

	// If trailing, append to args
	if message.Trailing != "" {
		message.Args = append(message.Args, message.Trailing)
	}

	return message, nil
}

func (m *Message) Bytes() []byte {
	return []byte(m.Raw)
}
