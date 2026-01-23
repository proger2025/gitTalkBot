package githubmanager

import (
	"log"
	"os"
	"testnewsaas/analyze"
	"testnewsaas/llm"
	"testnewsaas/parsing"
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

	err = llm.GenerateMarkdownFromTxt("D:/Go Proj/newSaas/tmpForLlm/"+repo+"_"+owner+".txt", "D:/Go Proj/newSaas/tmpForLlmMd/"+repo+"_"+owner+".md")
	// D:\Go Proj\newSaas\tmpForLlm\color_fatih.txt
	if err != nil {
		log.Println(err)
	}

	os.RemoveAll("D:/Go Proj/newSaas/tmp/" + owner + repo)
	os.RemoveAll("D:/Go Proj/newSaas/tmpForLlm/" + repo + "_" + owner + ".txt")
	// os.RemoveAll("D:/Go Proj/newSaas/tmpForLlmMd" + repo + "_" + owner + ".md")

	log.Println("END")

	return "D:/Go Proj/newSaas/tmpForLlmMd/"+repo+"_"+owner+".md", err

}

