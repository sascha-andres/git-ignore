package domain

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"golang.org/x/exp/slices"
)

type (
	Application struct {
		// global is set to true when the global ignore file is about to be used
		global bool
		// unique filters duplicates in gitignore
		unique bool
	}

	// ApplicationOption is for functional params
	ApplicationOption func(app *Application) error
)

// WithUnique toggles uniqueness of patterns in .gitignore
func WithUnique() ApplicationOption {
	return func(app *Application) error {
		app.unique = true
		return nil
	}
}

// WithGlobal tells the application to work on ~/.config/git/ignore
func WithGlobal() ApplicationOption {
	return func(app *Application) error {
		app.global = true
		return nil
	}
}

// NewApplication constructs a new ignore handler
func NewApplication(opts ...ApplicationOption) (*Application, error) {
	a := &Application{}
	for i := range opts {
		err := opts[i](a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// userHomeDir returns the users home directory
func userHomeDir() string {
	home := os.Getenv("XDG_CONFIG_HOME")
	if home != "" {
		return home
	}
	home = os.Getenv("HOME")
	if home != "" {
		return home
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("unable to read home directory: %q", err)
	}
	return home
}

// getGitIgnoreFileData returns git ignore line by line
func (a *Application) getGitIgnoreFileData() ([]string, error) {
	file := a.getGitIgnorePath()
	_, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(file)
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// getGitIgnorePath builds a path to the expected git ignore
func (a *Application) getGitIgnorePath() string {
	file := ".gitignore"
	if a.global {
		homeDirectory := userHomeDir()
		file = path.Join(homeDirectory, ".config", "git", "ignore")
	}
	return file
}

// unique will return an array with unique values
func unique[T comparable](in []T) ([]T, error) {
	result := []T{}
	for _, t := range in {
		if !slices.Contains(result, t) {
			result = append(result, t)
		}
	}
	return result, nil
}

// remove will return an array with element removed
func remove[T comparable](in []T, element T) ([]T, error) {
	result := []T{}
	for _, t := range in {
		if t == element {
			continue
		}
		result = append(result, t)
	}
	return result, nil
}

// List will print out the git ignore file
func (a *Application) List() error {
	data, err := a.getGitIgnoreFileData()
	if err != nil {
		return err
	}
	for _, line := range data {
		fmt.Println(line)
	}
	return nil
}

// Add will add a pattern to gitignore file
func (a *Application) Add(pattern string) error {
	return a.work(pattern, func(strings []string, s string) ([]string, error) {
		if strings == nil {
			return nil, errors.New("nil line data provided")
		}
		strings = append(strings, pattern)
		return strings, nil
	})
}

// Remove will remove one pattern from git ignore
func (a *Application) Remove(pattern string) error {
	return a.work(pattern, func(strings []string, s string) ([]string, error) {
		if strings == nil {
			return nil, errors.New("nil line data provided")
		}
		return remove(strings, pattern)
	})
}

// work will perform operation, worker is a function passed doing the actual work while
// this is a wrapper doing the common stuff such as reading and writing the file
func (a *Application) work(pattern string, worker func([]string, string) ([]string, error)) error {
	var lines []string
	var err error

	if pattern == "" {
		return errors.New("no pattern provided")
	}

	file := a.getGitIgnorePath()

	_, err = os.Stat(file)
	if err == nil {
		lines, err = a.getGitIgnoreFileData()
		if err != nil {
			return err
		}
	} else {
		lines = []string{}
	}
	lines, err = worker(lines, pattern)
	if err != nil {
		return err
	}
	if a.unique {
		lines, err = unique(lines)
		if err != nil {
			return err
		}
	}
	_ = os.Remove(file)
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	for _, line := range lines {
		_, _ = f.WriteString(fmt.Sprintln(line))
	}
	return f.Close()
}
