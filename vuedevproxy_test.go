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

func catchStdout(t *testing.T, run func()) string {
	realStdout := os.Stdout
	defer func() { os.Stdout = realStdout }()

	r, fakeStdout, err := os.Pipe()
	assert.NoError(t, err)

	os.Stdout = fakeStdout

	run()

	assert.NoError(t, fakeStdout.Close())

	stdout, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.NoError(t, r.Close())

	return string(stdout)
}

func testVueDevProxyValid(t *testing.T, runnerType RunnerType) {
	stdout := catchStdout(t, func() {
		proxy, err := NewVueDevProxy(&VueDevProxyOptions{
			RunnerType: runnerType,
			Dir:        "examples/webapps/vue/",
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

		assert.Contains(t, string(content), "/js/app.js")
		assert.Contains(t, string(content), "We're sorry but vue doesn't work properly without JavaScript enabled. Please enable it to continue.")
	})
	assert.Contains(t, stdout, "Starting development server")
	assert.Contains(t, stdout, "App running at")
}

func TestVueDevProxy_Valid(t *testing.T) {
	testVueDevProxyValid(t, RunnerTypeNpm)
}

func TestVueDevProxyYarn_Valid(t *testing.T) {
	testVueDevProxyValid(t, RunnerTypeYarn)
}

func TestVueDevProxy_SpecificPort(t *testing.T) {
	stdout := catchStdout(t, func() {
		proxy, err := NewVueDevProxy(&VueDevProxyOptions{
			Dir:  "examples/webapps/vue/",
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

func TestVueDevProxy_InvalidDir(t *testing.T) {
	proxy, err := NewVueDevProxy(&VueDevProxyOptions{
		Dir: "does_not_exist",
	})
	assert.NoError(t, err)

	err = proxy.Start(context.Background())
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}
