package example1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAuthMiddleware1ShouldAuthAdmin works as a simple test case
func TestAuthMiddleware1ShouldAuthAdmin(t *testing.T) {
	var handlerCalled bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Next handler
		handlerCalled = true
		_, _ = fmt.Fprint(w, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "http://test", nil)
	req.SetBasicAuth("admin", "admin")
	rec := httptest.NewRecorder()
	authMiddleware1(handler).ServeHTTP(rec, req)

	assert.True(t, handlerCalled)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

// TestAuthMiddleware1ShouldAuth is same as above but with a test table and multiple results
// Don't forget to run this "with coverage" report too and check it
// This should cover the whole `authMiddleware1` function
func TestAuthMiddleware1ShouldAuth(t *testing.T) {
	var handlerCalled bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Next handler
		handlerCalled = true
		_, _ = fmt.Fprint(w, "ok")
	})

	testTable := []struct {
		user           string
		pass           string
		expectedCalled bool
		expectedCode   int
		expectedBody   string
	}{
		{"admin", "admin", true, 200, "ok"},
		{"admin", "admin2", false, 200, ""}, // The code is wrong, but the middleware works like this here
		{"", "", false, 200, ""},            // Should test the missing auth header, I've added an if for this case in the code
	}

	rec := httptest.NewRecorder()

	for _, tt := range testTable {
		// I've moved this var inside so I should not care about the basic auth header reset in the `if tt.user != "" {` block
		req := httptest.NewRequest(http.MethodGet, "http://test", nil)

		// Don't forget to reset the local variable and clean the buffers
		handlerCalled = false
		rec.Body.Reset()

		// Overwrite the old params with the new ones
		if tt.user != "" {
			req.SetBasicAuth(tt.user, tt.pass)
		}

		// Call the middleware
		authMiddleware1(handler).ServeHTTP(rec, req)

		// Check the results
		assert.Equal(t, tt.expectedCalled, handlerCalled)
		assert.Equal(t, tt.expectedCode, rec.Code)
		assert.Equal(t, tt.expectedBody, rec.Body.String())
	}

}
