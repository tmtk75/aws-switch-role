package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/Atrox/homedir"
	"gopkg.in/ini.v1"
)

type ProfileName string
type URL string
type Links map[ProfileName]URL

func main() {
	// Parse
	cfg := ParseFile("/.aws/config", "^profile (.*)")
	creds := ParseFile("/.aws/credentials", "(.*)")

	// Merge
	links := make(Links)
	for k, e := range *cfg {
		links[k] = e
	}
	for k, e := range *creds {
		links[k] = e
	}

	// Find max length
	max := 0
	for k := range links {
		if max < len(k) {
			max = len(k)
		}
	}

	// Print
	w := tabwriter.NewWriter(os.Stdout, max+2, 2, 0, ' ', tabwriter.TabIndent)
	for k, l := range links {
		fmt.Fprintf(w, "%v\t%v\n", k, l)
		w.Flush()
	}
}

func ParseFile(path, sectionpattern string) *Links {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	cfg, err := ini.Load(home + path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	links := make(Links)
	re := regexp.MustCompile(sectionpattern)
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

		url := "https://signin.aws.amazon.com/switchrole?roleName=" + roleName +
			"&account=" + accountId +
			"&displayName=" + profileName

		links[ProfileName(profileName)] = URL(url)
	}

	return &links
}
