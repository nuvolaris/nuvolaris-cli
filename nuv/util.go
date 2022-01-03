package main

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// buffer for dry run rusults
var dryRunBuf = []string{}

// DryRunPush saves dummy results for dry run execution
func DryRunPush(buf ...string) {
	dryRunBuf = buf
}

// DryRunPop returns a value from the buffer of dry run results
// returns an empty string if the  buffer is empty
func DryRunPop(buf ...string) string {
	res := ""
	if len(dryRunBuf) > 0 {
		res = dryRunBuf[0]
		dryRunBuf = dryRunBuf[1:]
	}
	return res
}

// SysErr executes a command in a convenient way:
// it splits the paramenter in arguments if separated by spaces,
// then accepts multiple arguments;
// logs errors in stderr and prints output in stdout;
// also returns output as a string, or an error if there is an error
// If the command starts with "@" do not print the output.
func SysErr(cli string, args ...string) (string, error) {
	return sysErr(false, cli, args...)
}

// DryRunSysErr performs a dry run of SysErr
// in this case it always prints the command
func DryRunSysErr(cli string, args ...string) (string, error) {
	return sysErr(true, cli, args...)
}

func sysErr(dryRun bool, cli string, args ...string) (string, error) {
	re := regexp.MustCompile(`[\r\t\n\f ]+`)
	a := strings.Split(re.ReplaceAllString(cli, " "), " ")
	params := args
	if len(a) > 1 {
		params = append(a[1:], args...)
	}
	exe := strings.TrimPrefix(a[0], "@")
	silent := strings.HasPrefix(a[0], "@")
	if dryRun {
		if len(params) > 0 {
			fmt.Printf("%s %s\n", exe, strings.Join(params, " "))
		} else {
			fmt.Println(exe)
		}
		res := DryRunPop()
		if strings.HasPrefix(res, "!") {
			return "", errors.New(res[1:])
		}
		return res, nil
	}

	log.Tracef("< %s %v\n", exe, params)
	cmd := exec.Command(exe, params...)
	out, err := cmd.CombinedOutput()
	res := string(out)
	if err != nil {
		log.Tracef("> ERROR: %s", err.Error())
		return "", err
	}
	log.Tracef("> %s", res)
	if !silent {
		fmt.Print(res)
	}
	return res, nil
}
