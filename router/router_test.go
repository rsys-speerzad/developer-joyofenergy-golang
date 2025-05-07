package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"joi-energy-golang/domain"
)

func TestServer(t *testing.T) {
	port := os.Getenv("PORT")
	if err := os.Setenv("PORT", "8081"); err != nil {
		t.Fatalf("Failed to set PORT environment variable: %v", err)
	}
	defer func() {
		if err := os.Setenv("PORT", port); err != nil {
			t.Fatalf("Failed to restore PORT environment variable: %v", err)
		}
	}()

	server := NewServer()
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()
	defer func() {
		if err := server.Close(); err != nil {
			t.Errorf("Failed to close server: %v", err)
		}
	}()

	// Wait 50 milliseconds for server to start listening to requests
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://localhost:8081/readings/read/smartMeterId")

	assert.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("Failed to close response body: %v", err)
		}
	}()

	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	assert.Equalf(t, expectedContentType, actualContentType, "handler returned wrong content-type: got %v want %v", actualContentType, expectedContentType)
}

func TestEndpointEndpointSuccess(t *testing.T) {
	testHandler := newHandler()

	rr := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/readings/read/smartMeterId", nil)
	req.Header.Add("Content-Type", "application/json")

	testHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	expected := domain.StoreReadings{
		SmartMeterId:        "smartMeterId",
		ElectricityReadings: nil,
	}

	var actual domain.StoreReadings
	err := json.Unmarshal(rr.Body.Bytes(), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
