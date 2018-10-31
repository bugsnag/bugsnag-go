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
