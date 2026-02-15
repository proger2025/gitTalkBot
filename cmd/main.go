package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"gittalk/internal/githubManager"
	
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

type Step int

const (
	StepNone Step = iota
	StepWaitRepoGo
)

type userSettingsBot struct {
	UID          int64
	UserLanguage string // "en" / "ru"
	Step         Step
}

var (
	mu             sync.Mutex
	users          = map[int64]*userSettingsBot{}
	userActiveLink = make(map[int64]bool)
)

func getUser(uid int64) *userSettingsBot {
	mu.Lock()
	defer mu.Unlock()

	u, ok := users[uid]
	if !ok {
		u = &userSettingsBot{
			UID:          uid,
			UserLanguage: "en", // default
			Step:         StepNone,
		}
		users[uid] = u
	}
	return u
}

func main() {
	log.Println("Start bot")

	// bot settings

	err := godotenv.Load(".env")

	b, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("tokenBot"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	// main menu
	mainMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
	btnGo := mainMenu.Text("Golang")
	btnSettings := mainMenu.Text("Settings")

	mainMenu.Reply(
		mainMenu.Row(btnGo),
		mainMenu.Row(btnSettings),
	)

	makeSettingsMarkup := func(u *userSettingsBot) *tele.ReplyMarkup {
		menu := &tele.ReplyMarkup{}
		caption := fmt.Sprintf("Language: %s", u.UserLanguage)
		btnLang := menu.Data(caption, "userlang")
		menu.Inline(menu.Row(btnLang))
		return menu
	}

	b.Handle("/start", func(c tele.Context) error {
		u := getUser(c.Sender().ID)

		if u.UserLanguage == "en" {
			return c.Send("Hello! Choose from supported programming languages.", mainMenu)
		}
		return c.Send("Привет! Выберите один из поддерживаемых языков программирования.", mainMenu)
	})

	b.Handle(&btnGo, func(c tele.Context) error {
		u := getUser(c.Sender().ID)
		userActiveLink[u.UID] = true
		u.Step = StepWaitRepoGo

		if u.UserLanguage == "en" {
			return c.Send("Send a link to the repository (GOLANG ONLY!)")
		}
		return c.Send("Отправьте ссылку на репозиторий (ТОЛЬКО GOLANG!)")
	}) // golang

	// get link

	b.Handle(tele.OnText, func(c tele.Context) error {
		u := getUser(c.Sender().ID)
		if userActiveLink[u.UID] == true {
			if u.Step != StepWaitRepoGo {
				return nil
			}

			ghLink := c.Text()
			if !strings.HasPrefix(ghLink, "https://github.com/") {
				userActiveLink[u.UID] = false

				if u.UserLanguage == "en" {
					return c.Send("Incorrect link. try again.\n")
				} else {
					return c.Send("Некорректная ссылка. Попробуйте еще раз.")
				}

			}

			log.Printf("User: %v . Link: %v", u.UID, ghLink)
			u.Step = StepNone

			c.Send("I accepted the repository, I'm starting the analysis...\n")

			filePathMd, err := githubmanager.GetInfoGit(ghLink)

			if err != nil {
				log.Println(err)
			}

			doc := &tele.Document{
				File:     tele.FromDisk(filePathMd),
				FileName: "api.md",
				Caption:  "Generated documentation in Markdown",
			}

			c.Send(doc)

			os.RemoveAll(filePathMd)

			
		}
		return nil
		
	})

	// settings

	b.Handle(&btnSettings, func(c tele.Context) error {
		u := getUser(c.Sender().ID)
		return c.Send(`
Settings

Version: 0.1v (beta)

`, makeSettingsMarkup(u))
})

	b.Handle(&tele.Btn{Unique: "userlang"}, func(c tele.Context) error {
		u := getUser(c.Sender().ID)

		if u.UserLanguage == "en" {
			u.UserLanguage = "ru"

			return c.Edit("Язык изменён на русский", makeSettingsMarkup(u))
		}

		u.UserLanguage = "en"
		return c.Edit("Language changed to English", makeSettingsMarkup(u))
	})

	b.Start()
}



