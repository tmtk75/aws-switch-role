package main

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("/Users/tsakuma/.aws/config")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	re := regexp.MustCompile("^profile (.*)")
	for _, e := range cfg.Sections() {
		fmt.Printf("----\n")
		f := re.FindAllStringSubmatch(e.Name(), -1)
		if len(f) == 0 {
			continue
		}

		name := f[0][1] // [profile name]
		fmt.Printf("%v\n", name)
		//for _, k := range e.Keys() {
		//	fmt.Printf("%v\n", k.Name())
		//}
		k, err := e.GetKey("role_arn")
		if err != nil {
			//log.Fatalf("%v", err)
			continue
		}
		fmt.Printf("%v\n", k)
	}
}
