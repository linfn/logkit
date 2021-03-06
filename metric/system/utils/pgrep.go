package utils

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type PID int32

var ExecCommand = exec.Command

type PIDFinder interface {
	PidFile(path string) ([]PID, error)
	Pattern(pattern string) ([]PID, error)
	Uid(user string) ([]PID, error)
	FullPattern(path string) ([]PID, error)
}

// Implemention of PIDGatherer that execs pgrep to find processes
type Pgrep struct {
	path string
}

func NewPgrep() (PIDFinder, error) {
	path, err := exec.LookPath("pgrep")
	if err != nil {
		return nil, fmt.Errorf("could not find pgrep binary: %s", err)
	}
	return &Pgrep{path: path}, nil
}

func (pg *Pgrep) PidFile(path string) ([]PID, error) {
	var pids []PID
	pidString, err := ioutil.ReadFile(path)
	if err != nil {
		return pids, fmt.Errorf("failed to read pidfile '%s'. Error: '%s'",
			path, err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(pidString)))
	if err != nil {
		return pids, err
	}
	pids = append(pids, PID(pid))
	return pids, nil
}

func (pg *Pgrep) Pattern(pattern string) ([]PID, error) {
	args := []string{pattern}
	return find(pg.path, args)
}

func (pg *Pgrep) Uid(user string) ([]PID, error) {
	args := []string{"-u", user}
	return find(pg.path, args)
}

func (pg *Pgrep) FullPattern(pattern string) ([]PID, error) {
	args := []string{"-f", pattern}
	return find(pg.path, args)
}

func find(path string, args []string) ([]PID, error) {
	out, err := run(path, args)
	if err != nil {
		return nil, err
	}
	return ParseOutput(out)
}

func run(path string, args []string) (string, error) {
	out, err := exec.Command(path, args...).Output()
	if err != nil {
		return "", fmt.Errorf("error running %s: %s", path, err)
	}
	return string(out), err
}

func ParseOutput(out string) ([]PID, error) {
	pids := []PID{}
	fields := strings.Fields(out)
	for _, field := range fields {
		pid, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		if err == nil {
			pids = append(pids, PID(pid))
		}
	}
	return pids, nil
}
