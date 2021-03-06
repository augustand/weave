// docker_tls_args: find the docker daemon's tls args
// This reimplements pgrep, because pgrep will only look in /proc, but not in
// $PROCFS.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func dockerTLSArgs(args []string) error {
	if len(args) > 0 {
		cmdUsage("docker-tls-args", "")
	}
	procRoot := os.Getenv("PROCFS")
	if procRoot == "" {
		procRoot = "/proc"
	}
	dirEntries, err := ioutil.ReadDir(procRoot)
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		dirName := dirEntry.Name()
		if _, err := strconv.Atoi(dirName); err != nil {
			continue
		}

		if comm, err := ioutil.ReadFile(filepath.Join(procRoot, dirName, "comm")); err != nil || string(comm) != "docker\n" {
			continue
		}

		cmdline, err := ioutil.ReadFile(filepath.Join(procRoot, dirName, "cmdline"))
		if err != nil {
			continue
		}

		isDaemon := false
		tlsArgs := []string{}
		args := bytes.Split(cmdline, []byte{'\000'})
		for i := 0; i < len(args); i++ {
			arg := string(args[i])
			switch {
			case arg == "-d" || arg == "daemon":
				isDaemon = true
				break
			case arg == "--tls", arg == "--tlsverify":
				tlsArgs = append(tlsArgs, arg)
			case strings.HasPrefix(arg, "--tls"):
				tlsArgs = append(tlsArgs, arg)
				if len(args) > i+1 &&
					!strings.Contains(arg, "=") &&
					!strings.HasPrefix(string(args[i+1]), "-") {
					tlsArgs = append(tlsArgs, string(args[i+1]))
					i++
				}
			}
		}
		if !isDaemon {
			continue
		}

		fmt.Println(strings.Join(tlsArgs, " "))
		return nil
	}

	return fmt.Errorf("cannot locate running docker daemon")
}
