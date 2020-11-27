package telegram

import (
	"errors"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	sessions        = sync.Map{}
	newMenuFunction func() Menu
)

type sessionData struct {
	m  Menu
	ts time.Time
}

func updateSession(chatID int64) (Menu, bool) {
	sess, ok := sessions.Load(chatID)
	if !ok {
		data := &sessionData{
			m:  newMenuFunction(),
			ts: time.Now(),
		}

		sessions.Store(chatID, data)
		return data.m, true
	}

	sess.(*sessionData).ts = time.Now()

	sessions.Store(chatID, sess)
	return sess.(*sessionData).m, false
}

func cleanUp() {
	sessions.Range(func(k, v interface{}) bool {
		s := v.(*sessionData)
		if time.Since(s.ts) > time.Hour*24 {
			sessions.Delete(k)
		}

		return true
	})
}

func sendMessage(chatID int64, resp Response) ([]tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(chatID, resp.buttonText())
	msg.ParseMode = "html"
	switch t := resp.(type) {
	case *message:
		return []tgbotapi.MessageConfig{msg}, nil
	case *button:
		msg.ReplyMarkup = createKeyboard(t.keys, 3)
		return []tgbotapi.MessageConfig{msg}, nil
	default:
		return nil, errors.New("invalid message type")
	}
}

// Update handle updates from telegram
func Update(update tgbotapi.Update) ([]tgbotapi.MessageConfig, error) {
	if update.Message == nil {
		return nil, nil
	}

	chatID := update.Message.Chat.ID
	m, newSession := updateSession(chatID)
	if newSession {
		return sendMessage(chatID, m.Reset())
	}

	return sendMessage(chatID, m.Process(update.Message.Text))
}

// InitLibrary should be called before calling update
func InitLibrary(fn func() Menu) {
	newMenuFunction = fn

	go func() {
		t := time.NewTicker(time.Hour)
		for range t.C {
			cleanUp()
		}
	}()
}
