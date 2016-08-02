package nimbus

// Config for Nimbus
type Config struct {
	Channels []string

	RealName string
	UserName string
	Password string
	Modes    string

	Debug int
}
