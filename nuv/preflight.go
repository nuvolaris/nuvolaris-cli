package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/units"
	"github.com/coreos/go-semver/semver"
	"github.com/mitchellh/go-homedir"

	log "github.com/sirupsen/logrus"
)

// Preflight perform preflight checks
func Preflight(skipDockerVersion bool, dir string) (string, error) {
	info, err := dockerInfo(false)
	if err != nil {
		return "", err
	}
	err = checkDockerMemory(info)
	if err != nil {
		return "", err
	}
	if !skipDockerVersion {
		err = ensureDockerVersion(false)
		if err != nil {
			return "", err
		}
	}
	err = isInHomePath(dir)
	if err != nil {
		return "", err
	}
	return info, nil
}

func ensureDockerVersion(dryRun bool) error {
	version, err := dockerVersion(dryRun)
	if err != nil {
		return err
	}
	vA := semver.New(MinDockerVersion)
	vB := semver.New(strings.TrimSpace(version))
	if vB.Compare(*vA) == -1 {
		return fmt.Errorf("installed docker version %s is no longer supported", vB)
	}
	return nil
}

func isInHomePath(dir string) error {
	// do not check if the directory is empty
	if dir == "" {
		return nil
	}
	homePath, err := homedir.Dir()
	if err != nil {
		return err
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(dir, homePath) {
		return fmt.Errorf("work directory %s should be below your home directory %s;\nthis is required to be accessible by Docker", dir, homePath)
	}
	return nil
}

// checkDockerMemory checks docker memory
func checkDockerMemory(info string) error {
	var search = regexp.MustCompile(`Total Memory: (.*)`)
	result := search.FindString(string(info))
	if result == "" {
		return fmt.Errorf("docker is not running")
	}
	mem := strings.Split(result, ":")
	memory := strings.TrimSpace(mem[1])
	n, err := units.ParseStrictBytes(memory)
	if err != nil {
		return err
	}
	log.Debug("mem:", n)
	//fmt.Println(n)
	if n <= int64(MinDockerMem) {
		return fmt.Errorf("nuv needs 4GB memory allocatable on docker")
	}
	return nil

}
