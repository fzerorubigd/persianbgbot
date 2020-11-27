package main

import (
	"log"
	"os"
	"syscall"

	"github.com/fzerorubigd/clictx"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ogier/pflag"

	_ "github.com/fzerorubigd/persianbgbot/internal/bloodrage"
	"github.com/fzerorubigd/persianbgbot/pkg/menu"
	"github.com/fzerorubigd/persianbgbot/pkg/telegram"
)

func main() {
	ctx := clictx.Context(syscall.SIGABRT, syscall.SIGTERM, syscall.SIGKILL)

	var (
		token    string
		menuSize int
		debug    bool
	)
	pflag.StringVar(&token, "token", "", "Telegram bot token, if it's empty, it tries the env TELEGRAM_BOT_TOKEN")
	pflag.IntVar(&menuSize, "menu-size", 7, "Items in menu (other than navigation items)")
	pflag.BoolVar(&debug, "debug", false, "Show debug log")

	pflag.Parse()
	if token == "" {
		token = os.Getenv("TELEGRAM_BOT_TOKEN")
	}
	telegram.InitLibrary(func() telegram.Menu {
		m, err := menu.CreateMemoryMenu(menuSize, menu.AllGames()...)
		if err != nil {
			log.Fatal(err)
		}

		return m
	})

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case update := <-updates:
			msg, err := telegram.Update(update)
			if err != nil {
				log.Println(err)
			}

			for i := range msg {
				if _, err := bot.Send(msg[i]); err != nil {
					log.Println(err)
				}
			}
		case <-ctx.Done():
			log.Print("Done")
			return
		}
	}
}
