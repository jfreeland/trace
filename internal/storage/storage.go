package storage

import (
	"github.com/jfreeland/trace/internal/data"
)

// Storage defines the interface that should be implemented
// if it is needed to add a new storage type.
type Storage interface {
	// Get(key string) (interface{}, bool)
	// GetAll() map[string]interface{}
	// For troubleshooting.
	PrettyPrint()

	GetAllHosts() map[string]*data.Host
	GetHost(ip string) (*data.Host, bool)
	StoreHost(host *data.Host)

	GetAllResults() map[string][]*data.TracerouteResult
	GetResults(test string) ([]*data.TracerouteResult, bool)
	StoreResult(test string, result *data.TracerouteResult)
}
