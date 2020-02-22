package spaproxy

import (
	"regexp"
	"strconv"
)

// AngularDevProxyOptions represents Angular development proxy options.
type AngularDevProxyOptions struct {
	RunnerType RunnerType
	ScriptName string
	Dir        string
	Port       int
	Production bool
}

// NewAngularDevProxy returns new SpaDevProxy instance for Angular development.
func NewAngularDevProxy(options *AngularDevProxyOptions) (SpaDevProxy, error) {
	port, err := GetFreePort(options.Port)
	if err != nil {
		return nil, err
	}

	scriptName := options.ScriptName
	if len(scriptName) == 0 {
		scriptName = "start"
	}

	args := make([]string, 0)
	args = append(args, "--port", strconv.Itoa(port))
	args = append(args, "--host", "localhost")
	args = append(args, "--open", "false")
	args = append(args, "--prod")
	if options.Production {
		args = append(args, "true")
	} else {
		args = append(args, "false")
	}

	return NewSpaDevProxy(&SpaDevProxyOptions{
		RunnerType:  options.RunnerType,
		ScriptName:  scriptName,
		Dir:         options.Dir,
		Env:         []string{},
		Args:        args,
		Port:        port,
		StartRegexp: regexp.MustCompile("is listening on"),
	})
}
