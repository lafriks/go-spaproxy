package spaproxy

import (
	"fmt"
	"regexp"
	"strconv"
)

// VueDevProxyOptions represents VueJS development proxy options.
type VueDevProxyOptions struct {
	RunnerType RunnerType
	ScriptName string
	Dir        string
	Port       int
	Mode       string
}

// NewVueDevProxy returns new SpaDevProxy instance for VueJS development.
func NewVueDevProxy(options *VueDevProxyOptions) (SpaDevProxy, error) {
	port, err := GetFreePort(options.Port)
	if err != nil {
		return nil, err
	}

	env := []string{
		fmt.Sprintf("PORT=%d", port),
		fmt.Sprintf("DEV_SERVER_PORT=%d", port),
		"BROWSER=none",
	}

	scriptName := options.ScriptName
	if len(scriptName) == 0 {
		scriptName = "serve"
	}

	args := make([]string, 0)
	args = append(args, "--port", strconv.Itoa(port))
	args = append(args, "--host", "localhost")
	if len(options.Mode) > 0 {
		args = append(args, "--mode", options.Mode)
	}

	return NewSpaDevProxy(&SpaDevProxyOptions{
		RunnerType:  options.RunnerType,
		ScriptName:  scriptName,
		Dir:         options.Dir,
		Env:         env,
		Args:        args,
		Port:        port,
		StartRegexp: regexp.MustCompile("App running at"),
	})
}
