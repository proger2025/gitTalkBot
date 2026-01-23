package analyze

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CheckSanity(repoPath string) error {
	checkGo, err := findFiles(repoPath)

	if err != nil {
		log.Println(err)
		return err
	}

	if checkGo {
		log.Println("OK")
		return nil
	} 

	return errors.New("No .go files found")

}


func findFiles(repoPath string) (bool, error) {
	var errStopWalk error = errors.New("stopWalk")
	// trash = .git, vendor, testdata, .github
	var checkGo bool
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} 

		if info.IsDir() && (info.Name() == ".git" || info.Name() == "vendor" || info.Name() == ".github" || info.Name() == "testdata") {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			checkGo = true
			return errStopWalk
		}

		return nil

	})

	if err == errStopWalk {
		err = nil
	}

	return checkGo, err
}

