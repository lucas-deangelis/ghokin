package ghokin

import (
	"bytes"
	"fmt"
	"os"
	mpath "path"
	"path/filepath"
	"sync"
)

type commands map[string]string

type indent struct {
	backgroundAndScenario int
	step                  int
	tableAndDocString     int
}

// FileManager handles transformation on feature files
type FileManager struct {
	indentConf indent
	commands   commands
}

// NewFileManager creates a brand new FileManager, it requires indentation values and commands defined
// as a shell commands in comments
func NewFileManager(backgroundAndScenarioIndent int, stepIndent int, tableAndDocStringIndent int, commands map[string]string) FileManager {
	return FileManager{
		indent{
			backgroundAndScenarioIndent,
			stepIndent,
			tableAndDocStringIndent,
		},
		commands,
	}
}

// Transform formats and applies shell commands on feature file
func (f FileManager) Transform(filename string) (bytes.Buffer, error) {
	section, err := extractSections(filename)

	if err != nil {
		return bytes.Buffer{}, err
	}

	return transform(section, f.indentConf, f.commands)
}

// TransformAndReplace formats and applies shell commands on file or folders
// and replace the content of files
func (f FileManager) TransformAndReplace(path string, extensions []string) []error {
	errors := []error{}

	fi, err := os.Stat(path)
	if err != nil {
		return append(errors, err)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		errors = append(errors, f.replaceFolderWithContent(path, extensions)...)
	case mode.IsRegular():
		b, err := f.Transform(path)

		if err != nil {
			return append(errors, err)
		}

		if err := replaceFileWithContent(path, b.Bytes()); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func (f FileManager) replaceFolderWithContent(path string, extensions []string) []error {
	errors := []error{}
	fc := make(chan string)
	wg := sync.WaitGroup{}
	var mu sync.Mutex
	defer mu.Unlock()

	files, err := findFeatureFiles(path, extensions)

	if err != nil {
		return []error{err}
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			for file := range fc {
				b, err := f.Transform(file)

				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf(`An error occurred with file "%s" : %s`, file, err.Error()))
					mu.Unlock()
					continue
				}

				if err := replaceFileWithContent(file, b.Bytes()); err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf(`An error occurred with file "%s" : %s`, file, err.Error()))
					mu.Unlock()
				}
			}

			wg.Done()
		}()
	}

	for _, file := range files {
		fc <- file
	}

	close(fc)
	wg.Wait()
	mu.Lock()

	return errors
}

func replaceFileWithContent(filename string, content []byte) error {
	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	_, err = file.Write(content)

	return err
}

func findFeatureFiles(rootPath string, extensions []string) ([]string, error) {
	files := []string{}

	if err := filepath.Walk(rootPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, extension := range extensions {
			if !info.IsDir() && mpath.Ext(p) == "."+extension {
				files = append(files, p)
				break
			}
		}

		return nil
	}); err != nil {
		return []string{}, err
	}

	return files, nil
}
