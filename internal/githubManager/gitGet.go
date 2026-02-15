package githubmanager

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v76/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func githubData(url string) (string, string, string, string, error) {
	err := godotenv.Load(".env")
	if err != nil  {
		log.Println(err)
	}

	apiKeyGit := os.Getenv("apiKeyGit") // api key


	ownerGit, repoGit := normalizeUrl(url)

	token := apiKeyGit

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repo, _, err := client.Repositories.Get(ctx, ownerGit, repoGit)
	if err != nil {
		log.Println(err)
	}

	repoDisc := repo.GetDescription()

	repo.GetCommitsURL()

	fmt.Println("Connect ✅")
	repoPath, err := CloneRepo(url, ownerGit, repoGit)
	if err != nil {
		log.Println(err)
	} else {
		return repoPath, ownerGit, repoGit, repoDisc , nil
	}

	return "", "", "", "", err
	  
}


func normalizeUrl(url string) (string, string) {
	var owner string 
	url = strings.TrimPrefix(url, "https://github.com/")

	for _, v := range url {
		if string(v) == "/" {
			break
		}
		owner += string(v)
	}
	
	return owner, strings.TrimPrefix(url, owner + "/")
}
