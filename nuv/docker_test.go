package main

import "fmt"

func ExampleDockerVersion() {
	//*DryRunFlag = false
	DryRunPush("19.03.5", "!no docker")
	out, err := dockerVersion(true)
	fmt.Println(out, err)
	// out, err = dockerVersion(true)
	// fmt.Println(out, err)
	// Output:
	// docker version --format {{.Server.Version}}
	// 19.03.5 <nil>
	// docker version --format {{.Server.Version}}
	//  no docker
}

func ExampleDockerInfo() {
	DryRunPush("!bad", "Info: hello")
	out, err := dockerInfo(true)
	fmt.Println(err, out+"*")
	out, err = dockerInfo(true)
	fmt.Println(err, out+"*")
	// Output:
	// docker info
	// docker is not running *
	// docker info
	// <nil> Info: hello*
}
