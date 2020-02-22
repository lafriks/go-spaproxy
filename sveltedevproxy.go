package spaproxy

import (
	"fmt"
	"regexp"
)

// SvelteDevProxyOptions represents Svelte development proxy options.
type SvelteDevProxyOptions struct {
	RunnerType RunnerType
	ScriptName string
	Dir        string
	Port       int
}

// NewSvelteDevProxy returns new SpaDevProxy instance for SvelteJS development.
func NewSvelteDevProxy(options *SvelteDevProxyOptions) (SpaDevProxy, error) {
	port, err := GetFreePort(options.Port)
	if err != nil {
		return nil, err
	}

	env := []string{
		fmt.Sprintf("PORT=%d", port),
		"HOST=localhost",
	}

	scriptName := options.ScriptName
	if len(scriptName) == 0 {
		scriptName = "dev"
	}

	return NewSpaDevProxy(&SpaDevProxyOptions{
		RunnerType:  options.RunnerType,
		ScriptName:  scriptName,
		Dir:         options.Dir,
		Env:         env,
		Args:        []string{},
		Port:        port,
		StartRegexp: regexp.MustCompile("Your application is ready"),
	})
}
