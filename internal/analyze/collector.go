package analyze

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)


func CountFiles(repoPath string) ([]string, error) {
	repoFlies, err := collectFiles(repoPath)

	if err != nil {
		log.Println(err)
		return repoFlies, err
	}

	if len(repoFlies) > 0 {	
		return repoFlies, nil
	} else {
		err = errors.New("No files of the required extension were found")
	}

	return repoFlies, err
}


func collectFiles(repoPath string) ([]string, error) {
	var repoFlies []string
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} 

		if info.IsDir() && (info.Name() == ".git" || info.Name() == "vendor" || info.Name() == ".github" || info.Name() == "testdata") {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			repoFlies = append(repoFlies, path)
		}

		return nil
	})

	return repoFlies, err
}