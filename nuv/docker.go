package main

import "fmt"

func dockerInfo(dryRun bool) (string, error) {
	var out string
	var err error
	if dryRun {
		out, err = DryRunSysErr("@docker info")
	} else {
		out, err = SysErr("@docker info")
	}
	if err != nil {
		return "", fmt.Errorf("docker is not running")
	}
	return out, nil
}

func dockerVersion(dryRun bool) (string, error) {
	if dryRun {
		return DryRunSysErr("@docker version --format {{.Server.Version}}")
	}
	return SysErr("@docker version --format {{.Server.Version}}")
}
