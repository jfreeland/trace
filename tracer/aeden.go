package tracer

// This is named aiden because github.com/jfreeland/traceroute is based off of
// https://github.com/aeden/traceroute

import (
	"context"
	"fmt"
	"sync"

	"github.com/jfreeland/traceroute"
	tr "github.com/jfreeland/traceroute"
	log "github.com/sirupsen/logrus"

	"github.com/jfreeland/trace/storage"
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
	go run(trace)
	for {
		select {
		case <-trace.waitCh:
			go run(trace)
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

func run(trace *Trace) {
	// defer trace.wg.Done()
	// trace.wg.Add(1)
	options := tr.Options{}
	c := make(chan traceroute.Hop)
	go func() {
		for {
			hop, ok := <-c
			if !ok {
				return
			}
			printHop(hop)
		}
	}()
	_, err := tr.Traceroute(trace.host, &options, c)
	if err != nil {
		log.Printf("that sucks: %v", err)
	}
	select {
	case <-trace.ctx.Done():
		return
	default:
		trace.waitCh <- 1
	}
}

func printHop(hop traceroute.Hop) {
	addr := fmt.Sprintf("%v.%v.%v.%v", hop.Address[0], hop.Address[1], hop.Address[2], hop.Address[3])
	hostOrAddr := addr
	if hop.Host != "" {
		hostOrAddr = hop.Host
	}
	if hop.Success {
		fmt.Printf("%-3d %v (%v)  %v\n", hop.TTL, hostOrAddr, addr, hop.ElapsedTime)
	} else {
		fmt.Printf("%-3d *\n", hop.TTL)
	}
}
