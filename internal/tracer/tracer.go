package tracer

type Tracer interface {
	Run(host string)
	Stop(host string)
}
