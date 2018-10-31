package controllers

import (
	"fmt"

	"github.com/bugsnag/bugsnag-go"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Handled() revel.Result {
	bugsnag.Notify(fmt.Errorf("Oops"), c.Args["context"])
	return c.Render()
}

func (c App) Unhandled() revel.Result {
	crash := func(i interface{}) string {
		return i.(string)
	}
	crash(struct{}{})
	return c.Render()
}

func (c App) Configure() revel.Result {
	bugsnag.Notify(fmt.Errorf("Oops"), c.Args["context"])
	return c.Render()
}

func (c App) Metadata() revel.Result {
	bugsnag.Notify(fmt.Errorf("Oops"), bugsnag.MetaData{
		"Account": {
			"Name":           "Company XYZ",
			"Price(dollars)": "1 Million",
		},
	})
	return c.Render()
}

func (c App) OnBeforeNotify() revel.Result {
	bugsnag.OnBeforeNotify(
		func(event *bugsnag.Event, config *bugsnag.Configuration) error {
			if event.Message == "Ignore this error" {
				return fmt.Errorf("not sending errors to ignore")
			}
			// continue notifying as normal
			if event.Message == "Change error message" {
				event.Message = "Error message was changed"
			}
			return nil
		})

	notifier := bugsnag.New()
	notifier.NotifySync(fmt.Errorf("Don't ignore this error"), true)
	notifier.NotifySync(fmt.Errorf("Ignore this error"), true)
	notifier.NotifySync(fmt.Errorf("Change error message"), true)
	return c.Render()
}

func (c App) Recover() revel.Result {
	go func() {
		defer bugsnag.Recover()
		panic("Go routine killed")
	}()
	return c.Render()
}
