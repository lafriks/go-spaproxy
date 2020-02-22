package spaproxy

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync/atomic"
	"syscall"
)

// SpaDevProxy instance.
type SpaDevProxy interface {
	// Start backgroud process and wait when it is ready to accept connections.
	Start(ctx context.Context) error

	// Stop background process by killing it.
	Stop() error

	// HandleFunc handles reverse proxy function.
	HandleFunc(w http.ResponseWriter, r *http.Request)
}

// SpaDevProxyOptions options to use for SpaDevProxy.
type SpaDevProxyOptions struct {
	// RunnerType to use. Defaults to NPM.
	RunnerType RunnerType

	// ScriptName is the name of the script to run.
	ScriptName string

	// Args holds command line arguments.
	Args []string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	// As a special case on Windows, SYSTEMROOT is always added if
	// missing and not explicitly set to the empty string.
	Env []string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Start runs the command in the
	// calling process's current directory.
	Dir string

	// StartRegexp specifies regular expression used to detect
	// when background process is ready to accept incomming requests.
	StartRegexp *regexp.Regexp

	// Port specifies localhost port that background service is
	// using to accept requests.
	Port int

	// ShowBuildInfo specifies option either to show webpack building
	// progress or not.
	ShowBuildInfo bool
}

type spaDevProxy struct {
	options *SpaDevProxyOptions
	cmd     *exec.Cmd
	proxy   *httputil.ReverseProxy
}

// NewSpaDevProxy creates new SpaDevProxy instance.
func NewSpaDevProxy(options *SpaDevProxyOptions) (SpaDevProxy, error) {
	remote, err := url.Parse(fmt.Sprintf("http://localhost:%d", options.Port))
	if err != nil {
		return nil, err
	}

	return &spaDevProxy{
		options: options,
		proxy:   httputil.NewSingleHostReverseProxy(remote),
	}, nil
}

// Start backgroud process and wait when it is ready to accept connections.
func (p *spaDevProxy) Start(ctx context.Context) error {
	if _, err := os.Stat(p.options.Dir); err != nil {
		return err
	}
	path, args := prepareRunner(p.options.RunnerType, p.options.ScriptName, p.options.Args...)
	p.cmd = newCommand(ctx, path, args...)
	p.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	p.cmd.Env = append(p.options.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	p.cmd.Dir = p.options.Dir

	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}

	done := p.forwardOutput(stdout, stderr)

	err = p.cmd.Start()

	<-done

	return err
}

var ansiRegexp = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

func (p *spaDevProxy) cleanOutput(line string) string {
	line = ansiRegexp.ReplaceAllString(line, "")
	if !p.options.ShowBuildInfo && strings.HasPrefix(line, "<s> [webpack.Progress]") {
		return ""
	}
	return line
}

func (p *spaDevProxy) forwardOutput(stdout io.ReadCloser, stderr io.ReadCloser) chan struct{} {
	done := make(chan struct{})

	var c int32

	stdoutScanner := bufio.NewScanner(stdout)
	go func() {
		for stdoutScanner.Scan() {
			line := p.cleanOutput(stdoutScanner.Text())
			if len(line) == 0 {
				continue
			}
			if c == 0 && p.options.StartRegexp != nil && p.options.StartRegexp.MatchString(line) {
				if atomic.CompareAndSwapInt32(&c, 0, 1) {
					done <- struct{}{}
					close(done)
				}
			}
			fmt.Printf("%s\n", line)
		}

		if atomic.CompareAndSwapInt32(&c, 0, 1) {
			done <- struct{}{}
			close(done)
		}
	}()

	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			line := p.cleanOutput(stderrScanner.Text())
			if len(line) == 0 {
				continue
			}
			fmt.Printf("ERR: %s\n", line)
		}

		if atomic.CompareAndSwapInt32(&c, 0, 1) {
			done <- struct{}{}
			close(done)
		}
	}()

	return done
}

// Stop background process by killing it.
func (p *spaDevProxy) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	// Using syscall as Process.Kill does not kill child processes
	if err := syscall.Kill(-p.cmd.Process.Pid, syscall.SIGKILL); err != nil {
		return err
	}
	return p.cmd.Wait()
}

// HandleFunc handles reverse proxy function.
func (p *spaDevProxy) HandleFunc(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
