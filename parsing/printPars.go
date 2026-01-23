package parsing

import (
	"log"
	"os"
	"strings"
)

func PrintFromParser(symbols []Symbol, owner, repo, repoDesc string) {
	nameFile := "D:/Go Proj/newSaas/tmpForLlm/" + repo + "_" + owner + ".txt"

	file, err := os.Create(nameFile)
	if err != nil {
		log.Println(err)
	}
	
	defer file.Close()

	var builder strings.Builder

	builder.WriteString("Ты анализируешь Golang библиотеку.\n")
	builder.WriteString("Объясни её так: \n")
	builder.WriteString("1) Напиши описание библиотеки: Зачем нужна в кратце, Когда использовать, С чем можно сравнить эту библиотеку в экосистеме Go или других языках, Важные нюансы (На что нужно обратить внимание) \n")
	builder.WriteString("2) Приведи 2-3 примера что можно делать с помощью этой библиотеки\n")
	builder.WriteString("3.1) Напиши заголовок ОСНОВНЫЕ ФУНКЦИИ и объясни их \n")
	builder.WriteString("3.2) Напиши заголовок ВТОРОСТЕПЕННЫЕ и объясни их \n")
	builder.WriteString("4) Покажи минимальный работоспособный пример (код) для самой типичной задачи этой библиотеки с комментариями. (2 примера) \n")
	builder.WriteString("5) Назови 1-2 смежные или часто используемые с ней библиотеки  \n")


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



