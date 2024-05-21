package bugsnag

import "fmt"

type reportPublisher interface {
	publishReport(*payload) error
	publishReportPool(p *payload) error
}

type defaultReportPublisher struct{}

func (*defaultReportPublisher) publishReport(p *payload) error {
	//p.logf("notifying bugsnag: %s", p.Message)
	if !p.notifyInReleaseStage() {
		return fmt.Errorf("not notifying in %s", p.ReleaseStage)
	}
	if p.Synchronous {
		return p.deliver()
	}

	go func(p *payload) {
		if err := p.deliver(); err != nil {
			// Ensure that any errors are logged if they occur in a goroutine.
			p.logf("bugsnag/defaultReportPublisher.publishReport: %v", err)
		}
	}(p)
	return nil
}

func (*defaultReportPublisher) publishReportPool(p *payload) error {
	if !p.notifyInReleaseStage() {
		return fmt.Errorf("not notifying in %s", p.ReleaseStage)
	}

	if !asyncPool.TrySubmit(func() {
		if err := p.deliver(); err != nil {
			// Ensure that any errors are logged if they occur in a goroutine.
			p.logf("bugsnag/defaultReportPublisher.publishReport: %v", err)
		}
	}) {
		return fmt.Errorf("failed to submit report to async pool")
	}
	return nil
}
