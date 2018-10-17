package controllers

import (
	"fmt"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Handled() revel.Result {
	bugsnag.Notify(fmt.Errorf("oopsie"), c.Args["context"])
	return c.Render()
}

func (c App) Unhandled() revel.Result {
	crash(struct{}{})
	return c.Render()
}

func crash(a interface{}) string {
	return a.(string)
}
