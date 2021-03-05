package storage

import (
	"fmt"
	"sync"

	"github.com/jfreeland/trace/internal/data"
)

// InMemory is an in memory storage.
type InMemory struct {
	mu      sync.RWMutex
	hosts   map[string]*data.Host
	results map[string][]*data.TracerouteResult
}

// NewInMemory creates a new in memory storage.
func NewInMemory() *InMemory {
	hosts := make(map[string]*data.Host)
	results := make(map[string][]*data.TracerouteResult)

	return &InMemory{
		hosts:   hosts,
		results: results,
	}
}

// // Get a single key.
// func (s *InMemory) Get(key string) (interface{}, bool) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	val, ok := s.data[key]

// 	return val, ok
// }

// // GetAll returns all values in storage.
// func (s *InMemory) GetAll() map[string]interface{} {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()

// 	return s.data
// }

// PrettyPrint prints all results.
func (s *InMemory) PrettyPrint() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fmt.Println("----------")
	fmt.Println("Hosts:")
	for _, host := range s.hosts {
		fmt.Printf("%v: %v\n", host.IP, host.Meta.Address)
	}
	fmt.Println("---")
	fmt.Println("Results:")
	for _, results := range s.results {
		for _, result := range results {
			fmt.Printf("%v\n", result.Time)
			for _, hop := range result.Hops {
				fmt.Printf("%v: %v\n", hop.Host.IP, hop.Duration.Milliseconds())
			}
		}
	}
	fmt.Println("----------")
}

// GetAllHosts returns all hosts.
func (s *InMemory) GetAllHosts() map[string]*data.Host {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.hosts
}

// GetHost returns a host.
func (s *InMemory) GetHost(ip string) (*data.Host, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.hosts[ip]

	return val, ok
}

// StoreHost a payload.
func (s *InMemory) StoreHost(host *data.Host) {
	s.mu.RLock()
	_, ok := s.hosts[host.IP]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		s.hosts[host.IP] = host
		s.mu.Unlock()
	}
}

// GetAllResults returns all results.
func (s *InMemory) GetAllResults() map[string][]*data.TracerouteResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.results
}

// GetResults returns all traceroute results for a host being tested.
func (s *InMemory) GetResults(test string) ([]*data.TracerouteResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.results[test]

	return val, ok
}

// StoreResult stores a traceroute test result.
func (s *InMemory) StoreResult(test string, result *data.TracerouteResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[test] = append(s.results[test], result)
}
