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

	// отделяем автора от репозитория https://github.com/author/repo

	ownerGit, repoGit := normalizeUrl(url)

	// 1. Берём токен
	token := apiKeyGit

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// 4. GitHub client
	client := github.NewClient(tc)

	// 5. Тестовый запрос (инфо о репозитории)
	repo, _, err := client.Repositories.Get(ctx, ownerGit, repoGit)
	if err != nil {
		log.Println(err)
	}

	repoDisc := repo.GetDescription()

	repo.GetCommitsURL()

	fmt.Println("Подключение успешно ✅")
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
