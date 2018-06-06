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

type Link struct {
	ProfileName string
	Href        string
}

func main() {
	// Parse
	cfg := ParseFile("/.aws/config", "^profile (.*)")
	creds := ParseFile("/.aws/credentials", "(.*)")

	// Merge
	links := make([]*Link, 0)
	for _, e := range cfg {
		links = append(links, e)
	}
	for _, e := range creds {
		links = append(links, e)
	}

	// Find max
	max := 0
	for _, l := range links {
		if max < len(l.ProfileName) {
			max = len(l.ProfileName)
		}
	}

	// Print
	w := tabwriter.NewWriter(os.Stdout, max+2, 2, 0, ' ', tabwriter.TabIndent)
	for _, l := range links {
		fmt.Fprintf(w, "%v\t%v\n", l.ProfileName, l.Href)
		w.Flush()
	}
}

func ParseFile(path, sectionpattern string) []*Link {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	cfg, err := ini.Load(home + path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	links := make([]*Link, 0)
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

		l := "https://signin.aws.amazon.com/switchrole?roleName=" + roleName +
			"&account=" + accountId +
			"&displayName=" + profileName

		links = append(links, &Link{ProfileName: profileName, Href: l})
	}

	return links
}
