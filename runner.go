package spaproxy

// RunnerType runner type.
type RunnerType int

const (
	// RunnerTypeNpm runner type NPM.
	RunnerTypeNpm RunnerType = iota

	// RunnerTypeNpx runner type NPX.
	RunnerTypeNpx

	//RunnerTypeYarn runner type Yarn.
	RunnerTypeYarn

	// RunnerTypeCustom custom runner type.
	RunnerTypeCustom
)

func prepareRunner(runnerType RunnerType, scriptName string, args ...string) (string, []string) {
	if runnerType == RunnerTypeCustom {
		return scriptName, args
	}

	var path string
	switch runnerType {
	case RunnerTypeNpm:
		path = "npm"
	case RunnerTypeNpx:
		path = "npx"
	case RunnerTypeYarn:
		path = "yarn"
	}

	a := make([]string, 0, len(args)+1)

	if runnerType == RunnerTypeNpm {
		a = append(a, "run")
	}

	a = append(a, scriptName)

	if runnerType == RunnerTypeNpm {
		a = append(a, "--")
	}

	a = append(a, args...)

	return path, a
}
