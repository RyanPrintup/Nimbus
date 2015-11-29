package nimbus

import "strings"

const (
	// 512 minus CR-LF
	MAX_SIZE = 510
)

type Message struct {
	Prefix string
	Command string
	Args []string
	Middle string
	Trailing string
}

func crlfCutsetFunc(r rune) bool {
	return r == '\r' || r == '\n'
}

/**
 * Prefix
 * 	Message starts with :
 *	Space indicates end of prefix
 *
 * Command
 * 	Follows prefix if there is one
 *	Space indicates end
 *
 *
 *	Args
 *		Middle
 *			Seperated by space
 *			: indicates end
 *		Trailing
 *			Anything after :
 *
 */

func ParseMessage(raw string) *Message {
	message := &Message{} 

	// check if message is empty
	if raw = strings.TrimFunc(raw, crlfCutsetFunc); len(raw) < 2 {
		return nil
	}

	if raw[0] == ':' {
		prefix := strings.Split(raw, " ")
		message.Prefix = prefix[0][1 : len(prefix[0])]
		raw = raw[strings.Index(raw, " ") + 1 : len(raw)]
	}

	message.Command = strings.Split(raw, " ")[0]
	raw = raw[strings.Index(raw, " ") + 1 : len(raw)]

	if strings.Contains(raw, ":") {
		s := strings.Split(raw, ":")
		message.Middle, message.Trailing = s[0][0 : len(s[0]) - 1], s[1]
	} else {
		message.Middle = raw
	}

	if len(message.Middle) > 0 {
		message.Args = strings.Split(message.Middle, " ")
	}

	if len(message.Trailing) > 0 {
		message.Args = append(message.Args, strings.Split(message.Trailing, " ")...)
	}

	return message
}