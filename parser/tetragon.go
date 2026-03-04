package parser

import (
	"encoding/json"
	"path"
	"strings"
)

// AlertType indicates whether the alert is from a binary exec or a sensitive file read.
type AlertType string

const (
	AlertBinaryExec AlertType = "BINARY_EXEC"
	AlertFileRead   AlertType = "FILE_READ"
)

// AlertResult holds the data needed to build a Discord embed for either alert type.
type AlertResult struct {
	Type     AlertType
	Binary   string
	Args     string
	Pid      uint32
	FilePath string // set for FILE_READ
}

// TetragonEvent represents a single line of Tetragon JSON (process_exec or process_kprobe).
type TetragonEvent struct {
	ProcessExec   *ProcessExecPayload   `json:"process_exec"`
	ProcessKprobe *ProcessKprobePayload `json:"process_kprobe"`
}

// ProcessExecPayload is the payload under the "process_exec" key.
type ProcessExecPayload struct {
	Process *ProcessInfo `json:"process"`
}

// ProcessKprobePayload is the payload under the "process_kprobe" key (file access, etc.).
type ProcessKprobePayload struct {
	Process *ProcessInfo `json:"process"`
	Args    []KprobeArg  `json:"args"`
}

// KprobeArg represents one argument; we care about file_arg.path for file reads.
type KprobeArg struct {
	FileArg *struct {
		Path string `json:"path"`
	} `json:"file_arg"`
}

// ProcessInfo holds binary, arguments, and pid.
type ProcessInfo struct {
	Binary    string `json:"binary"`
	Arguments string `json:"arguments"`
	Pid       uint32 `json:"pid"`
}

// pathMatchesTarget normalizes and checks if path matches any target (exact or suffix).
func pathMatchesTarget(filePath string, targetFiles []string) bool {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return false
	}
	for _, target := range targetFiles {
		target = strings.TrimSpace(target)
		if target == "" {
			continue
		}
		if filePath == target {
			return true
		}
		// Normalize ~ so "/root/.aws/credentials" matches "~/.aws/credentials"
		normTarget := strings.TrimPrefix(target, "~")
		if normTarget != "" && (strings.HasSuffix(filePath, normTarget) || filePath == normTarget) {
			return true
		}
		if strings.HasSuffix(filePath, target) {
			return true
		}
	}
	return false
}

// ProcessLine unmarshals a JSON line from Tetragon. It returns an AlertResult and true if:
// - it's a process_exec and the binary matches one of targetBinaries, or
// - it's a process_kprobe and one of the file_arg paths matches one of targetFiles.
func ProcessLine(line []byte, targetBinaries, targetFiles []string) (*AlertResult, bool) {
	if len(line) == 0 {
		return nil, false
	}

	var ev TetragonEvent
	if err := json.Unmarshal(line, &ev); err != nil {
		return nil, false
	}

	// 1) process_exec: binary match
	if ev.ProcessExec != nil && ev.ProcessExec.Process != nil {
		binary := strings.TrimSpace(ev.ProcessExec.Process.Binary)
		if binary != "" {
			for _, target := range targetBinaries {
				target = strings.TrimSpace(target)
				if target == "" {
					continue
				}
				if binary == target || path.Base(binary) == path.Base(target) {
					return &AlertResult{
						Type:   AlertBinaryExec,
						Binary: ev.ProcessExec.Process.Binary,
						Args:   ev.ProcessExec.Process.Arguments,
						Pid:    ev.ProcessExec.Process.Pid,
					}, true
				}
			}
		}
	}

	// 2) process_kprobe: file path match
	if ev.ProcessKprobe != nil && ev.ProcessKprobe.Process != nil && len(ev.ProcessKprobe.Args) > 0 {
		for _, arg := range ev.ProcessKprobe.Args {
			if arg.FileArg == nil || arg.FileArg.Path == "" {
				continue
			}
			p := arg.FileArg.Path
			if pathMatchesTarget(p, targetFiles) {
				return &AlertResult{
					Type:     AlertFileRead,
					Binary:   ev.ProcessKprobe.Process.Binary,
					Args:     ev.ProcessKprobe.Process.Arguments,
					Pid:      ev.ProcessKprobe.Process.Pid,
					FilePath: p,
				}, true
			}
		}
	}

	return nil, false
}
