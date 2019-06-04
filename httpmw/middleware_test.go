package httpmw

import (
	"net/http"
	"testing"

	"github.com/sporkmonger/ecsevent"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerParent(t *testing.T) {
	assert := assert.New(t)
	monitor := ecsevent.New()
	mh := NewHandler(monitor)
	h := mh(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sm := FromRequest(r)
		assert.NotNil(sm)
		if sm != nil {
			assert.Equal(monitor, sm.Parent())
		}
	}))
	h.ServeHTTP(nil, &http.Request{})
}
