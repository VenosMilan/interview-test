package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	if err := cmd.Start(); err != nil {
		assert.Equal(t, err, nil)
	}

	time.Sleep(5 * time.Second)

	resp, err := http.Get("http://localhost:8080/readyz")
	assert.Equal(t, err, nil)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Equal(t, err, nil)

	assert.Equal(t, "OK", string(body))

	time.Sleep(5 * time.Second)

	err = cmd.Process.Signal(os.Interrupt)
	assert.Equal(t, err, nil)

	state, err := cmd.Process.Wait()
	assert.Equal(t, err, nil)

	assert.NotNil(t, state)
	assert.True(t, state.Exited())
}
