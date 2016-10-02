package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func Parse() {
	file, err := os.Open("config.ini")
	if err != nil {
		fmt.Println("Warning: no config.ini file found.")
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	reComment := regexp.MustCompile("^$|^\\s*#.*$")
	reKeyValue := regexp.MustCompile("^([^=\\s]+)\\s*=\\s*(|[^ ]|[^ ].*[^ ])\\s*$")

	overridenFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		overridenFlags[f.Name] = true
	})

	existingFlags := make(map[string]bool)
	flag.VisitAll(func(f *flag.Flag) {
		existingFlags[f.Name] = true
	})

	for scanner.Scan() {
		line := scanner.Text()
		if reComment.MatchString(line) {
			continue
		}
		matches := reKeyValue.FindStringSubmatch(line)
		if matches == nil {
			fmt.Printf("Error: failed to parse config.ini line \"%s\"\r\n", line)
			os.Exit(1)
		}
		if matches[0] != line {
			fmt.Printf("Warning: full-string-matching regexp didn't match full string. This is decidedly odd.\r\n")
		}
		if _, isOverriden := overridenFlags[matches[1]]; isOverriden {
			continue
		}
		if _, isExisting := existingFlags[matches[1]]; !isExisting {
			fmt.Printf("Error: invalid key in config.ini: \"%s\"\r\n", matches[1])
			os.Exit(1)
		}
		flag.Set(matches[1], matches[2])
	}
}
