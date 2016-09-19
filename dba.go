package blueprint

// Database describes a methods to connect and close database connections.
type Database interface {
	Dial(name string) error
	Name() string
	Close()
}
