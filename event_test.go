package bugsnag

import (
	"fmt"
	"testing"
)

func BenchmarkNewEvent(b *testing.B) {
	data := []interface{}{1223, 123.1234, "f234t24t34g", fmt.Errorf("Oopsie")}
	notifier := New()
	for i := 0; i < b.N; i++ {
		newEvent(data, notifier)
	}
}
