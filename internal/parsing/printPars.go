package parsing

import (
	"log"
	"os"
	"strings"
)

func PrintFromParser(symbols []Symbol, owner, repo, repoDesc string) {
	nameFile := "/tmpForLlm/" + repo + "_" + owner + ".txt"

	file, err := os.Create(nameFile)
	if err != nil {
		log.Println(err)
	}
	
	defer file.Close()

	var builder strings.Builder

	builder.WriteString("You are analyzing a Golang library.\n")
	builder.WriteString("Explain it as follows: \n")
	builder.WriteString("1) Write a description of the library: What is it for (briefly), When to use it, What can it be compared to in the Go ecosystem or other languages, Important nuances (What to pay attention to) \n")
	builder.WriteString("2) Provide 2-3 examples of what can be done using this library\n")
	builder.WriteString("3) Write the heading MAIN FUNCTIONS and explain them \n")
	builder.WriteString("4) Show a minimum viable example (code) for the most typical task of this library with comments (2 examples) \n")
	builder.WriteString("5) Name 1-2 related or frequently used libraries with it \n")


	builder.WriteString("Repository: " + owner + "/" + repo + "\n")
	builder.WriteString("Description: " + repoDesc + "\n\n")
	builder.WriteString("=================================\n\n")


	c := 0
	for _, v := range symbols {
		if c == 25  {
			break
		}

		builder.WriteString("\nSYMBOL:\n")
		builder.WriteString("ID: " + v.Id + "\n")
		builder.WriteString("KIND: " + v.Kind + "\n")
		builder.WriteString("PACKAGE: " + v.PackageName + "\n")
		builder.WriteString("SIGNATURE: " + v.Signature + "\n")
		
		

		if v.DocComment == "" {
			builder.WriteString("DOC: " + "(none) \n")
		} else {
			builder.WriteString("DOC: " + v.DocComment + "\n")
		}


		builder.WriteString("\n--- \n\n")

		c++
	}

	

	_, err = file.WriteString(builder.String())
	if err != nil {
		log.Println(err)
		return
	}


	log.Printf("file %v is input. Save...", nameFile)
	

}



