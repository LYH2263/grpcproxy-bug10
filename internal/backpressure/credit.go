package backpressure

// CreditPair binds upstream/downstream windows for a session.
type CreditPair struct {
	Upstream   *Window
	Downstream *Window
}

func NewCreditPair(max int) *CreditPair {
	return &CreditPair{
		Upstream:   NewWindow(max, "upstream"),
		Downstream: NewWindow(max, "downstream"),
	}
}

func (p *CreditPair) ReleaseUpstream(n int) { p.Upstream.Release(n) }
func (p *CreditPair) ReleaseDownstream(n int) { p.Downstream.Release(n) }
