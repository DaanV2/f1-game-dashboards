package game

type (
	packetProcessorOptions struct {
		host string
	}

	PacketOption = func(p *packetProcessorOptions)
)

func (p *packetProcessorOptions) _default() {
	p.host = ""
}

func (p *packetProcessorOptions) apply(opts ...PacketOption) {
	for _, o := range opts {
		o(p)
	}
}

func WithHost(host string) PacketOption {
	return func(p *packetProcessorOptions) {
		p.host = host
	}
}