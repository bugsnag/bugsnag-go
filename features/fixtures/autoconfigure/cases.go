package main

import (
	"fmt"

	"github.com/bugsnag/bugsnag-go/v2"
)

func explicitPanic() {
	panic("PANIQ!")
}

func handledEvent() {
	bugsnag.Notify(fmt.Errorf("gone awry!"))
}

func handledMetadata() {
	bugsnag.OnBeforeNotify(func(event *bugsnag.Event, config *bugsnag.Configuration) error {
		event.MetaData.Add("fruit", "Tomato", "beefsteak")
		event.MetaData.Add("snacks", "Carrot", "4")
		return nil
	})
	handledEvent()
}
