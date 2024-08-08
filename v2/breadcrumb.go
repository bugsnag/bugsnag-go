package bugsnag

import "time"

type BreadcrumbType = string

const (
	// Changing screens or content being displayed, with a defined destination and optionally a previous location.
	BreadcrumbTypeNavigation BreadcrumbType = "navigation"
	// Sending and receiving requests and responses.
	BreadcrumbTypeRequest BreadcrumbType = "request"
	// Performing an intensive task or query.
	BreadcrumbTypeProcess BreadcrumbType = "process"
	// Messages and severity sent to a logging platform.
	BreadcrumbTypeLog BreadcrumbType = "log"
	// Actions performed by the user, like text input, button presses, or confirming/cancelling an alert dialog.
	BreadcrumbTypeUser BreadcrumbType = "user"
	// Changing the overall state of an app, such as closing, pausing, or being moved to the background, as well as device state changes like memory or battery warnings and network connectivity changes.
	BreadcrumbTypeState BreadcrumbType = "state"
	// An error which was reported to Bugsnag encountered in the same session.
	BreadcrumbTypeError BreadcrumbType = "error"
	// User-defined, manually added breadcrumbs.
	BreadcrumbTypeManual BreadcrumbType = "manual"
)

// Key value metadata that is displayed with the breadcrumb.
type BreadcrumbMetaData map[string]interface{}

// Remove any values from meta-data that have keys matching the filters,
// and any that are recursive data-structures.
func (meta BreadcrumbMetaData) sanitize(filters []string) interface{} {
	return sanitizer{
		Filters: filters,
		Seen:    make([]interface{}, 0),
	}.Sanitize(meta)
}

type Breadcrumb struct {
	// The time at which the event occurred, in ISO 8601 format.
	Timestamp string
	// A short summary describing the event, such as the user action taken or a new application state.
	Name string
	// A category which describes the breadcrumb.
	Type BreadcrumbType
	// Additional information about the event, as key/value pairs.
	MetaData BreadcrumbMetaData
}

type maximumBreadcrumbsValue interface {
	isValid() bool
	trimBreadcrumbs(breadcrumbs []Breadcrumb) []Breadcrumb
}

type MaximumBreadcrumbs int

func (length MaximumBreadcrumbs) isValid() bool {
	return length >= 0 && length <= 100
}

func (length MaximumBreadcrumbs) trimBreadcrumbs(breadcrumbs []Breadcrumb) []Breadcrumb {
	if int(length) >= 0 && len(breadcrumbs) > int(length) {
		return breadcrumbs[:int(length)]
	}
	return breadcrumbs
}

type (
	// A breadcrumb callback that returns if the breadcrumb should be added.
	onBreadcrumbCallback func(*Breadcrumb) bool

	breadcrumbState struct {
		// These callbacks are run in reverse order and determine if the breadcrumb should be added.
		onBreadcrumbCallbacks []onBreadcrumbCallback
		// Currently added breadcrumbs in order from newest to oldest
		breadcrumbs []Breadcrumb
	}
)

// onBreadcrumb adds a callback to be run before a breadcrumb is added.
// If false is returned, the breadcrumb will be discarded.
func (breadcrumbs *breadcrumbState) onBreadcrumb(callback onBreadcrumbCallback) {
	if breadcrumbs.onBreadcrumbCallbacks == nil {
		breadcrumbs.onBreadcrumbCallbacks = []onBreadcrumbCallback{}
	}

	breadcrumbs.onBreadcrumbCallbacks = append(breadcrumbs.onBreadcrumbCallbacks, callback)
}

// Runs all the OnBreadcrumb callbacks, returning true if the breadcrumb should be added.
func (breadcrumbs *breadcrumbState) runBreadcrumbCallbacks(breadcrumb *Breadcrumb) bool {
	if breadcrumbs.onBreadcrumbCallbacks == nil {
		return true
	}

	// run in reverse order
	for i := range breadcrumbs.onBreadcrumbCallbacks {
		callback := breadcrumbs.onBreadcrumbCallbacks[len(breadcrumbs.onBreadcrumbCallbacks)-i-1]
		if !callback(breadcrumb) {
			return false
		}
	}
	return true
}

// Add the breadcrumb onto the list of breadcrumbs, ensuring that the number of breadcrumbs remains below maximumBreadcrumbs.
func (breadcrumbs *breadcrumbState) leaveBreadcrumb(message string, configuration *Configuration, rawData ...interface{}) {
	breadcrumb := Breadcrumb{
		Timestamp: time.Now().Format(time.RFC3339),
		Name:      message,
		Type:      BreadcrumbTypeManual,
		MetaData:  BreadcrumbMetaData{},
	}
	for _, datum := range rawData {
		switch datum := datum.(type) {
		case BreadcrumbMetaData:
			breadcrumb.MetaData = datum
		case BreadcrumbType:
			breadcrumb.Type = datum
		default:
			panic("Unexpected type")
		}
	}

	if breadcrumbs.runBreadcrumbCallbacks(&breadcrumb) {
		if breadcrumbs.breadcrumbs == nil {
			breadcrumbs.breadcrumbs = []Breadcrumb{}
		}
		breadcrumbs.breadcrumbs = append([]Breadcrumb{breadcrumb}, breadcrumbs.breadcrumbs...)
		if configuration.MaximumBreadcrumbs != nil {
			breadcrumbs.breadcrumbs = configuration.MaximumBreadcrumbs.trimBreadcrumbs(breadcrumbs.breadcrumbs)
		}
	}
}

func (configuration *Configuration) breadcrumbEnabled(breadcrumbType BreadcrumbType) bool {
	if configuration.EnabledBreadcrumbTypes == nil {
		return true
	}
	for _, enabled := range configuration.EnabledBreadcrumbTypes {
		if enabled == breadcrumbType {
			return true
		}
	}
	return false
}

func (breadcrumbs *breadcrumbState) leaveBugsnagStartBreadcrumb(configuration *Configuration) {
	if configuration.breadcrumbEnabled(BreadcrumbTypeState) {
		breadcrumbs.leaveBreadcrumb("Bugsnag loaded", configuration, BreadcrumbTypeState)
	}
}

func (breadcrumbs *breadcrumbState) leaveEventBreadcrumb(event *Event, configuration *Configuration) {
	if event == nil {
		return
	}
	if !configuration.breadcrumbEnabled(BreadcrumbTypeError) {
		return
	}
	metadata := BreadcrumbMetaData{
		"errorClass": event.ErrorClass,
		"message":    event.Message,
		"unhandled":  event.Unhandled,
		"severity":   event.Severity.String,
	}
	breadcrumbs.leaveBreadcrumb(event.Error.Error(), configuration, BreadcrumbTypeError, metadata)
}
