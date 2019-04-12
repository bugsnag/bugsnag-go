package controllers

import (
	"fmt"
	"os"

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
	if _, err := os.Open("nonexistent_file.txt"); err != nil {
		if errClass := os.Getenv("ERROR_CLASS"); errClass != "" {
			bugsnag.Notify(err, c.Args["context"], bugsnag.ErrorClass{Name: errClass})
		} else {
			bugsnag.Notify(err, c.Args["context"])
		}
	}
	return c.Render()
}

func (c App) Unhandled() revel.Result {
	// Invalid type assertion, will panic
	func(a interface{}) string {
		return a.(string)
	}(struct{}{})
	return c.Render()
}

func (c App) Session() revel.Result {
	return c.Render()
}

func (c App) AutoNotify() revel.Result {
	go func(ctx interface{}) {
		defer func() { recover() }()
		defer bugsnag.AutoNotify(ctx)
		panic("Go routine killed with auto notify")
	}(c.Args["context"])
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
	defer bugsnag.Recover(c.Args["context"])
	panic("Request killed but recovered")
}

func (c App) User() revel.Result {
	bugsnag.Notify(fmt.Errorf("oops"), bugsnag.User{
		Id:    "test-user-id",
		Name:  "test-user-name",
		Email: "test-user-email",
	})
	return c.Render()
}
