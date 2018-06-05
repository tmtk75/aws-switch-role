package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

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
		f := re.FindAllStringSubmatch(e.Name(), -1)
		if len(f) == 0 {
			continue
		}

		profileName := f[0][1] // [profile name]
		k, err := e.GetKey("role_arn")
		if err != nil {
			continue
		}

		a := strings.Split(k.String(), "/")
		roleName := a[1]
		a = strings.Split(a[0], ":")
		accountId := a[4]

		l := "https://signin.aws.amazon.com/switchrole?roleName=" + roleName +
			"&account=" + accountId +
			"&displayName=" + profileName

		fmt.Printf("# %v\n", profileName)
		fmt.Printf("%v\n", l)
		fmt.Printf("\n")
	}
}
