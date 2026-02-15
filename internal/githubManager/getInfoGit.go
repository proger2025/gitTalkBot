package githubmanager

import (
	"log"
	"os"
	"gittalk/internal/analyze"
	"gittalk/internal/llm"
	"gittalk/internal/parsing"
)

func GetInfoGit(url string) (string, error) {
	repoPath, owner, repo, repoDisc, err := githubData(url)

	if err != nil {
		log.Println(err)
	}

	err = analyze.CheckSanity(repoPath)

	if err != nil {
		log.Println(err)
	}

	repoFiles, err := analyze.CountFiles(repoPath)

	if err != nil {
		log.Println(err)
	}

	res := parsing.Ast(repoFiles)
	parsing.PrintFromParser(parsing.BuildSymbol(res), owner, repo, repoDisc)

	err = llm.GenerateMarkdownFromTxt("/tmpForLlm/"+repo+"_"+owner+".txt", "/tmpForLlmMd/"+repo+"_"+owner+".md")
	
	if err != nil {
		log.Println(err)
	}

	os.RemoveAll("/tmp/" + owner + repo)
	os.RemoveAll("/tmpForLlm/" + repo + "_" + owner + ".txt")


	log.Println("END")

	return "/tmpForLlmMd/"+repo+"_"+owner+".md", err

}

