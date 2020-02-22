package spaproxy

import (
	"fmt"
	"regexp"
	"strconv"
)

// ReactDevProxyOptions represents React development proxy options.
type ReactDevProxyOptions struct {
	RunnerType    RunnerType
	ShowBuildInfo bool
	ScriptName    string
	Dir           string
	Port          int
}

// NewReactDevProxy returns new SpaDevProxy instance for React development.
func NewReactDevProxy(options *ReactDevProxyOptions) (SpaDevProxy, error) {
	port, err := GetFreePort(options.Port)
	if err != nil {
		return nil, err
	}

	env := []string{
		fmt.Sprintf("PORT=%d", port),
		"BROWSER=none",
		"HOST=localhost",
	}

	scriptName := options.ScriptName
	if len(scriptName) == 0 {
		scriptName = "start"
	}

	args := make([]string, 0)
	args = append(args, "--port", strconv.Itoa(port))

	return NewSpaDevProxy(&SpaDevProxyOptions{
		RunnerType:  options.RunnerType,
		ScriptName:  scriptName,
		Dir:         options.Dir,
		Env:         env,
		Args:        args,
		Port:        port,
		StartRegexp: regexp.MustCompile("You can now view reactjs in the browser"),
	})
}
