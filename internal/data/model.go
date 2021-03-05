package data

import "time"

// HostMeta is metadata about a host.
// TODO: ASN, ... there's probably a library that has all the things I care
// about.
type HostMeta struct {
	Address string
}

// Host is a host along a network path.
type Host struct {
	IP   string
	Meta *HostMeta
}

// Hop is a traceroute hop.
type Hop struct {
	Host     *Host
	Duration time.Duration
}

// TracerouteResult is a slice of Hops.
type TracerouteResult struct {
	Time time.Time
	Hops []*Hop
}
