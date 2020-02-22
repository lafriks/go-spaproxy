package spaproxy

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSvelteDevProxy_Valid(t *testing.T) {
	stdout := catchStdout(t, func() {
		proxy, err := NewSvelteDevProxy(&SvelteDevProxyOptions{
			Dir: "examples/webapps/svelte/",
		})
		assert.NoError(t, err)

		err = proxy.Start(context.Background())
		assert.NoError(t, err)

		//nolint:errcheck
		defer proxy.Stop()

		ts := httptest.NewServer(http.HandlerFunc(proxy.HandleFunc))
		defer ts.Close()

		res, err := http.Get(ts.URL)
		assert.NoError(t, err)

		content, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)

		assert.Contains(t, string(content), "/build/bundle.js")
		assert.Contains(t, string(content), "Svelte app")
	})
	assert.Contains(t, stdout, "Your application is ready")
}

func TestSvelteDevProxy_SpecificPort(t *testing.T) {
	stdout := catchStdout(t, func() {
		proxy, err := NewSvelteDevProxy(&SvelteDevProxyOptions{
			Dir:  "examples/webapps/svelte/",
			Port: 12345,
		})
		assert.NoError(t, err)

		err = proxy.Start(context.Background())
		assert.NoError(t, err)

		// To not kill too fast
		time.Sleep(100 * time.Millisecond)

		//nolint:errcheck
		proxy.Stop()
	})
	assert.Contains(t, stdout, "http://localhost:12345")
}

func TestSvelteDevProxy_InvalidDir(t *testing.T) {
	proxy, err := NewSvelteDevProxy(&SvelteDevProxyOptions{
		Dir: "does_not_exist",
	})
	assert.NoError(t, err)

	err = proxy.Start(context.Background())
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
