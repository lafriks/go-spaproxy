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

func TestReactDevProxy_Valid(t *testing.T) {
	stdout := catchStdout(t, func() {
		proxy, err := NewReactDevProxy(&ReactDevProxyOptions{
			Dir: "examples/webapps/reactjs/",
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

		assert.Equal(t, "Express", res.Header.Get("X-Powered-By"))

		content, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)

		assert.Contains(t, string(content), "/static/js/main.chunk.js")
		assert.Contains(t, string(content), "You need to enable JavaScript to run this app.")
	})
	assert.Contains(t, stdout, "Starting the development server")
	assert.Contains(t, stdout, "You can now view reactjs in the browser")
}

func TestReactDevProxy_SpecificPort(t *testing.T) {
	stdout := catchStdout(t, func() {
		proxy, err := NewReactDevProxy(&ReactDevProxyOptions{
			Dir:  "examples/webapps/reactjs/",
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

func TestReactDevProxy_InvalidDir(t *testing.T) {
	proxy, err := NewReactDevProxy(&ReactDevProxyOptions{
		Dir: "does_not_exist",
	})
	assert.NoError(t, err)

	err = proxy.Start(context.Background())
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
