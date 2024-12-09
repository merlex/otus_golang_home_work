package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, file := range files {
		name := file.Name()
		if file.IsDir() || name == "" || strings.Contains(name, "=") {
			continue
		}
		content, e := func(dirname string, filename string) (string, error) {
			f, err := os.Open(filepath.Join(dirname, filename))
			if err != nil {
				return "", err
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			if !scanner.Scan() {
				return "", fmt.Errorf("error reading file %s", filename)
			}
			firstline := scanner.Text()
			if err := scanner.Err(); err != nil {
				return "", err
			}
			return strings.ReplaceAll(strings.TrimRight(firstline, " \t"), "\x00", "\n"), nil
		}(dir, name)

		if e != nil {
			content = ""
		}

		env[name] = EnvValue{
			Value:      content,
			NeedRemove: content == "",
		}
	}
	return env, nil
}
