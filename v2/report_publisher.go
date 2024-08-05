package bugsnag

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type reportPublisher interface {
	publishReport(*payload) error
	setMainProgramContext(context.Context)
}

func (defPub *defaultReportPublisher) delivery() {
	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh, syscall.SIGINT, syscall.SIGTERM)

waitForEnd:
	for {
		select {
		case <-signalsCh:
			defPub.isClosing = true
			break waitForEnd
		case <-defPub.mainProgramCtx.Done():
			defPub.isClosing = true
			break waitForEnd
		case p, ok := <-defPub.eventsChan:
			if ok {
				if err := p.deliver(); err != nil {
					// Ensure that any errors are logged if they occur in a goroutine.
					p.logf("bugsnag/defaultReportPublisher.publishReport: %v", err)
				}
			} else {
				p.logf("Event channel closed")
				return
			}
		}
	}

	// Send remaining elements from the queue
	close(defPub.eventsChan)
	for p := range defPub.eventsChan {
			if err := p.deliver(); err != nil {
				// Ensure that any errors are logged if they occur in a goroutine.
				p.logf("bugsnag/defaultReportPublisher.publishReport: %v", err)
			}
	}
}

type defaultReportPublisher struct {
	eventsChan     chan *payload
	mainProgramCtx context.Context
	isClosing      bool
}

func newPublisher() reportPublisher {
	defPub := defaultReportPublisher{isClosing: false, mainProgramCtx: context.TODO()}
	defPub.eventsChan = make(chan *payload, 100)

	go defPub.delivery()

	return &defPub
}

func (defPub *defaultReportPublisher) setMainProgramContext(ctx context.Context) {
	defPub.mainProgramCtx = ctx
}

func (defPub *defaultReportPublisher) publishReport(p *payload) error {
	p.logf("notifying bugsnag: %s", p.Message)
	if !p.notifyInReleaseStage() {
		return fmt.Errorf("not notifying in %s", p.ReleaseStage)
	}
	if p.Synchronous {
		return p.deliver()
	}

	if defPub.isClosing {
		return fmt.Errorf("main program is stopping, new events won't be sent")
	}

	select {
	case defPub.eventsChan <- p:
	default:
		p.logf("Events channel full. Discarding value")
	}

	return nil
}
