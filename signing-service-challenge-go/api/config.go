package api

import "time"

// Config collects all the configuration values, both needed or optional, for
// the server
type Config struct {
	ListenAddress   string
	RequestsTimeout time.Duration
	// TODO: more stuff like db connection strings etc...
}
