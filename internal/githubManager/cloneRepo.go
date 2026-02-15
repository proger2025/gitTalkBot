package githubmanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)


func CloneRepo(url, owner, repoGit string) (string, error) {
	var err error
	timer := 60 * time.Second
	
	ctx, cancel := context.WithTimeout(context.Background(), timer)
	defer cancel()

	repoPath, err := cloneRepoComannd(ctx, url, owner ,repoGit)

	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("timeout")
	}
	
	return repoPath, err

}


func cloneRepoComannd(ctx context.Context, url, owner, repoGit string) (string, error) {

	path := fmt.Sprintf("/tmp/%v%v", owner, repoGit)

	args := []string{"clone", "--depth=1", url, path}

	cmd := exec.CommandContext(ctx ,"git", args...)

	output, err := cmd.CombinedOutput() 

	if err != nil {
		log.Println(err)
		log.Println("rep:", string(output))
		return "", err
	}

	return path, nil
}

