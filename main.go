package main

import (
	"errors"
	"log"
	"os"

	"github.com/sascha-andres/flag"
	"github.com/sascha-andres/gitignore/domain"
)

var (
	global, unique bool
)

func main() {
	flag.SetEnvPrefix("GIT_IGNORE")
	flag.BoolVar(&global, "global", false, "operate on global git ignore file")
	flag.BoolVar(&unique, "unique", false, "ensure patterns are unique")
	flag.Parse()

	if err := run(); err != nil {
		log.Printf("error running application: %q", err)
		os.Exit(1)
	}
}

func run() error {
	opts := make([]domain.ApplicationOption, 0)
	if global {
		opts = append(opts, domain.WithGlobal())
	}
	if unique {
		opts = append(opts, domain.WithUnique())
	}
	a, err := domain.NewApplication(opts...)
	if err != nil {
		return err
	}

	verbs := flag.GetVerbs()
	pattern := ""
	remainingArgs := flag.NArg()
	if remainingArgs != 1 {
		return errors.New("you need to specify exactly one additional argument (pattern to ignore)")
	}
	pattern = flag.Arg(0)

	if len(verbs) > 0 {
		switch verbs[0] {
		case "list":
			return a.List()
		case "remove":
			return a.Remove(pattern)
		}
	}

	return a.Add(pattern)
}
