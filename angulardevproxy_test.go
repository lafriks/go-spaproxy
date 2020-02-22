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

func TestAngularDevProxy_Valid(t *testing.T) {
	stdout := catchStdout(t, func() {
		proxy, err := NewAngularDevProxy(&AngularDevProxyOptions{
			Dir: "examples/webapps/angular/",
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

		assert.Contains(t, string(content), "main.js")
		assert.Contains(t, string(content), "<app-root></app-root>")
	})
	assert.Contains(t, stdout, "Compiled successfully")
	assert.Contains(t, stdout, "is listening on")
}

func TestAngularDevProxy_SpecificPort(t *testing.T) {
	stdout := catchStdout(t, func() {
		proxy, err := NewAngularDevProxy(&AngularDevProxyOptions{
			Dir:  "examples/webapps/angular/",
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

func TestAngularDevProxy_InvalidDir(t *testing.T) {
	proxy, err := NewAngularDevProxy(&AngularDevProxyOptions{
		Dir: "does_not_exist",
	})
	assert.NoError(t, err)

	err = proxy.Start(context.Background())
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
