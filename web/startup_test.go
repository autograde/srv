package web

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestStartup(t *testing.T) {
	w := httptest.NewRecorder()
	// startup(w, r, nil)
	startup(w, nil, nil)
	fmt.Printf("%d - %s", w.Code, w.Body.String())
}
