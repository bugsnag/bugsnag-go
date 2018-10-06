package bugsnag

// Shipper is the expected interface for sending payload data to a backend
type Shipper interface {
	Deliver(p *Payload) error
}
