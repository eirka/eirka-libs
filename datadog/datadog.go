package datadog

import (
	"github.com/DataDog/datadog-go/statsd"
)

var (
	// Client holds our DataDog client
	Client        *statsd.Client
	clientAddress = "127.0.0.1:8125"
	clientBuffer  = 10
)

// New initializes a new DataDog client and sets it to the Client variable
func New() (err error) {
	// Create our new DataDog statsd client
	Client, err = statsd.NewBuffered(clientAddress, clientBuffer)
	if err != nil {
		return
	}
	return
}
