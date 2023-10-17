package healthcheck

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeGetReadyEndpoint(t *testing.T) {
	server := httptest.NewServer(MakeGetReadyEndpoint())

	resp, err := http.Get(server.URL + "/readyz")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	expected := "OK"
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("expected %s, got %s", expected, string(b))
	}
}
