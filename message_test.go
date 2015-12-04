package nimbus

import (
	"testing"
	"reflect"
)

var testMessages = [...]*Message {
	&Message{
		Raw: "JOIN #channel",
		Prefix: "",
		Command: "JOIN",
		Args: []string{ "#channel" },
		Middle: "#channel",
		Trailing: "",
	},
	&Message{
		Raw: ":Nimbus PRIVMSG #test :Test response",
		Prefix: "Nimbus",
		Command: "PRIVMSG",
		Args: []string{ "#test", "Test response" },
		Middle: "#test",
		Trailing: "Test response",
	},
	&Message{
		Raw: ":Nimbus!Nimbus@hostname COMMAND arg1 arg2 arg3 arg4 arg5 :Long random message for testing",
		Prefix: "Nimbus!Nimbus@hostname",
		Command: "COMMAND",
		Args: []string{ "arg1", "arg2", "arg3", "arg4", "arg5", "Long random message for testing" },
		Middle: "arg1 arg2 arg3 arg4 arg5",
		Trailing: "Long random message for testing",
	},
}

func TestParseMessage(t *testing.T) {
	for _, test := range testMessages {
		m, err := ParseMessage(test.Raw)

		if err != nil {
			t.Error("Failed to parse message ", err)
		}

		if !reflect.DeepEqual(m, test) {
			t.Error("Failed to parse message ", test.Raw)
			t.Logf("Output: %#v", m)
			t.Logf("Expected: %#v", test)
		}
	}
}