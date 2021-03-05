package tracer

// This is named aiden because github.com/jfreeland/traceroute is based off of
// https://github.com/aeden/traceroute

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jfreeland/traceroute"
	tr "github.com/jfreeland/traceroute"
	log "github.com/sirupsen/logrus"

	"github.com/jfreeland/trace/internal/data"
	"github.com/jfreeland/trace/internal/storage"
)

var running = make(map[string]*Trace)

type Aeden struct {
	db storage.Storage
}

type Trace struct {
	ctx    context.Context
	host   string
	wg     *sync.WaitGroup
	waitCh chan int
	quitCh chan bool
}

// NewAeden returns an aeden tracerouter
func NewAeden(db storage.Storage) *Aeden {
	return &Aeden{
		db: db,
	}
}

// NewTrace returns a traceroute instance
func NewTrace(ctx context.Context, host string, wg *sync.WaitGroup, waitCh chan int, quitCh chan bool) *Trace {
	return &Trace{
		ctx:    ctx,
		host:   host,
		wg:     wg,
		waitCh: waitCh,
		quitCh: quitCh,
	}
}

// Run runs a traceroute to host
func (a *Aeden) Run(host string) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	waitCh := make(chan int)
	quitCh := make(chan bool)
	defer close(waitCh)
	trace := NewTrace(ctx, host, &wg, waitCh, quitCh)
	running[host] = trace
	go a.run(trace)
	for {
		select {
		case <-trace.waitCh:
			go a.run(trace)
		case <-trace.quitCh:
			cancel()
			return
		}
	}
}

// Stop stops a running traceroute
func (a *Aeden) Stop(host string) {
	trace, ok := running[host]
	if !ok {
		log.Error("trace did not exist")
		return
	}
	delete(running, host)
	trace.quitCh <- true
}

func (a *Aeden) run(trace *Trace) {
	// defer trace.wg.Done()
	// trace.wg.Add(1)
	options := tr.Options{}
	c := make(chan traceroute.Hop)
	results := &data.TracerouteResult{}
	go func() {
		for {
			hop, ok := <-c
			if !ok {
				return
			}
			_, ok = a.db.GetHost(hop.AddressString())
			if !ok {
				a.db.StoreHost(&data.Host{
					IP: hop.AddressString(),
					Meta: &data.HostMeta{
						Address: hop.HostOrAddressString(),
					},
				})
			}
			host, _ := a.db.GetHost(hop.AddressString())
			hopResult := &data.Hop{
				Host:     host,
				Duration: hop.ElapsedTime,
			}
			results.Hops = append(results.Hops, hopResult)
			printHop(hop)
		}
	}()
	_, err := tr.Traceroute(trace.host, &options, c)
	if err != nil {
		log.Printf("that sucks: %v", err)
	}
	results.Time = time.Now()
	a.db.StoreResult(trace.host, results)
	select {
	case <-trace.ctx.Done():
		return
	default:
		trace.waitCh <- 1
	}
}

func printHop(hop traceroute.Hop) {
	if hop.Success {
		fmt.Printf("%-3d %v (%v)  %v\n", hop.TTL, hop.AddressString(), hop.HostOrAddressString(), hop.ElapsedTime)
	} else {
		fmt.Printf("%-3d *\n", hop.TTL)
	}
}
