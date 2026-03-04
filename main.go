package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vericore/tetragon-sentinel/discord"
	"github.com/vericore/tetragon-sentinel/parser"
)

func main() {
	var alertBinaries multiFlag
	var alertFiles multiFlag
	flag.Var(&alertBinaries, "alert-binary", "Binary path to alert on (e.g. /bin/bash). Can be repeated or comma-separated.")
	flag.Var(&alertFiles, "alert-file", "File path to alert on when read (e.g. /etc/shadow, .env). Can be repeated or comma-separated.")
	flag.Parse()

	// Flatten comma-separated values for binaries
	var targetBinaries []string
	for _, b := range alertBinaries {
		for _, s := range strings.Split(b, ",") {
			if t := strings.TrimSpace(s); t != "" {
				targetBinaries = append(targetBinaries, t)
			}
		}
	}
	// Flatten comma-separated values for files
	var targetFiles []string
	for _, f := range alertFiles {
		for _, s := range strings.Split(f, ",") {
			if t := strings.TrimSpace(s); t != "" {
				targetFiles = append(targetFiles, t)
			}
		}
	}

	if len(targetBinaries) == 0 && len(targetFiles) == 0 {
		targetBinaries = []string{"/bin/bash", "/usr/bin/curl"} // default for demo
	}

	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		webhookURL = os.Getenv("DISCORD_WEBHOOK_SENTINEL")
	}
	if webhookURL == "" {
		fmt.Fprintln(os.Stderr, "tetragon-sentinel: set DISCORD_WEBHOOK_URL or DISCORD_WEBHOOK_SENTINEL (see docs/DISCORD-SETUP.md)")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	// Support long lines (Tetragon JSON can be large)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		result, match := parser.ProcessLine(line, targetBinaries, targetFiles)
		if !match || result == nil {
			continue
		}

		alertTypeStr := string(result.Type)
		if err := discord.SendAlert(webhookURL, alertTypeStr, result.Binary, result.Args, result.Pid, result.FilePath); err != nil {
			fmt.Fprintf(os.Stderr, "tetragon-sentinel: discord send failed: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "tetragon-sentinel: read error: %v\n", err)
		os.Exit(1)
	}
}

// multiFlag allows repeated flags and comma-separated values.
type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}
