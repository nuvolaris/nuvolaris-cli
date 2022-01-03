package main

import "fmt"

func ExamplePreflightEnsureDockerVersion() {
	DryRunPush("19.03.5", "10.03.5", MinDockerVersion, "!no docker")
	fmt.Println(ensureDockerVersion(true))
	fmt.Println(ensureDockerVersion(true))
	fmt.Println(ensureDockerVersion(true))
	fmt.Println(ensureDockerVersion(true))
	// Output:
	// docker version --format {{.Server.Version}}
	// <nil>
	// docker version --format {{.Server.Version}}
	// installed docker version 10.3.5 is no longer supported
	// docker version --format {{.Server.Version}}
	// <nil>
	// docker version --format {{.Server.Version}}
	// no docker
}

func ExampleInHomePath() {
	fmt.Println(isInHomePath("/home/nuvolaris"))
	fmt.Println(isInHomePath("/var/run"))
	fmt.Println(isInHomePath(""))
	// Output:
	// <nil>
	// work directory /var/run should be below your home directory /home/nuvolaris;
	// this is required to be accessible by Docker
	// <nil>
}

func ExamplePreflightDockerMemory() {
	fmt.Println(checkDockerMemory("\nTotal Memory: 11GiB\n"))
	fmt.Println(checkDockerMemory("\nTotal Memory: 3GiB\n"))
	// Output:
	// <nil>
	// nuv needs 4GB memory allocatable on docker
}
