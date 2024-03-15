// Package golangspectester provides a test engine to check
// the compliance of several Go compilers and interpreters
// against the Go Specification: https://go.dev/ref/spec.
// The tests are self-contained and implemented in individual
// files.
package golangspectester

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TestFunc is a type describing the user provided test function.
type TestFunc func(path, expected string, isError bool, code io.Reader) bool

// Engine describes a test engine.
type Engine struct {
	tfp string   // test folder path
	tf  TestFunc // test function
}

// NewEngine returns a new test engine.
func NewEngine(testFolderPath string, testFunc TestFunc) *Engine {
	return &Engine{
		tfp: testFolderPath,
		tf:  testFunc,
	}
}

// Start executes the tests located in the engine test folder path,
// using the provided engine test function. It returns the error encountered.
func (e *Engine) Start() error {
	return filepath.WalkDir(e.tfp, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		expected, isErr := expectedFromComment(path)

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		if !e.tf(path, expected, isErr, f) {
			return fmt.Errorf("testsuite stopped by TestFunc: %q", path)
		}

		return nil
	})
}

func expectedFromComment(p string) (out string, isErr bool) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, p, nil, parser.ParseComments)
	if len(f.Comments) == 0 {
		return "", false
	}
	text := f.Comments[len(f.Comments)-1].Text()

	// Sometimes the comment ends with a space.
	// We need to use \s to avoid the IDE trimming the comment.
	text = strings.Replace(text, "\\s", " ", -1)
	if strings.HasPrefix(text, "Output:\n") {
		return strings.TrimPrefix(text, "Output:\n"), false
	}
	if strings.HasPrefix(text, "Error:\n") {
		return strings.TrimPrefix(text, "Error:\n"), true
	}
	return "", false
}
