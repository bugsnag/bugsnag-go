package bugsnag

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

type (
	// A breadcrumb callback that returns if the breadcrumb should be added.
	OnBreadcrumbCallback func(*Breadcrumb) bool

	BreadcrumbState struct {
		// These callbacks are run in reverse order and determine if the breadcrumb should be added.
		OnBreadcrumbCallbacks []OnBreadcrumbCallback
		// Currently added breadcrumbs in order from newest to oldest
		Breadcrumbs []Breadcrumb
	}
)

// OnBreadcrumb adds a callback to be run before a breadcrumb is added.
// If false is returned, the breadcrumb will be discarded.
func (breadcrumbs *BreadcrumbState) OnBreadcrumb(callback OnBreadcrumbCallback) {
	if breadcrumbs.OnBreadcrumbCallbacks == nil {
		breadcrumbs.OnBreadcrumbCallbacks = []OnBreadcrumbCallback{}
	}

	breadcrumbs.OnBreadcrumbCallbacks = append(breadcrumbs.OnBreadcrumbCallbacks, callback)
}

// Runs all the OnBreadcrumb callbacks, returning true if the breadcrumb should be added.
func (breadcrumbs *BreadcrumbState) runBreadcrumbCallbacks(breadcrumb *Breadcrumb) bool {
	if breadcrumbs.OnBreadcrumbCallbacks == nil {
		return true
	}

	// run in reverse order
	for i := range breadcrumbs.OnBreadcrumbCallbacks {
		callback := breadcrumbs.OnBreadcrumbCallbacks[len(breadcrumbs.OnBreadcrumbCallbacks)-i-1]
		if !callback(breadcrumb) {
			return false
		}
	}
	return true
}

// Add the breadcrumb onto the list of breadcrumbs, ensuring that the number of breadcrumbs remains below maximumBreadcrumbs.
func (breadcrumbs *BreadcrumbState) appendBreadcrumb(breadcrumb Breadcrumb, maximumBreadcrumbs int) error {
	if breadcrumbs.runBreadcrumbCallbacks(&breadcrumb) {
		if breadcrumbs.Breadcrumbs == nil {
			breadcrumbs.Breadcrumbs = []Breadcrumb{}
		}
		breadcrumbs.Breadcrumbs = append([]Breadcrumb{breadcrumb}, breadcrumbs.Breadcrumbs...)
		if len(breadcrumbs.Breadcrumbs) > 0 && len(breadcrumbs.Breadcrumbs) > maximumBreadcrumbs {
			breadcrumbs.Breadcrumbs = breadcrumbs.Breadcrumbs[:len(breadcrumbs.Breadcrumbs)-1]
		}
	}
	return nil
}
