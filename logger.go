package nimbus

import "log"

// SetLogFlags will set flags for go's log. Set to no flags by default.
func (c *Client) SetLogFlags(flags int) {
	log.SetFlags(flags)
}

// Lprintln will log a message if level is less than or equal to debug.
func (c *Client) Lprintln(level int, message string) {
	c.Lprintf(level, message+"\n")
}

// Lprintf is the same as Lprintln, except you can print with a format string and arguments.
func (c *Client) Lprintf(level int, format string, v ...interface{}) {
	if level <= c.debug {
		log.Printf(format, v...)
	}
}
