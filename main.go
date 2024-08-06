package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	Text     string `json:"text"`
	GroupID  string `json:"group_id"`
}

func main() {
	// Output file açma ve hata kontrolü
	outputFile, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer outputFile.Close()

	//TODO: TOKEN'I GÜVENLİK AÇISINDAN KODUN İÇİNE GÖMME DÜZELT
	var botToken string
	fmt.Print("Enter Bot Token: ")
	fmt.Scan(&botToken)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create new bot: %v", err)
	}

	bot.Debug = false

	// UpdateConfig yapısı ile güncellemeleri konfigüre etme
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	// Güncellemeleri almak için kanal oluşturma
	updates := bot.GetUpdatesChan(updateConfig)

	// Güncellemeleri işleme
	for update := range updates {

		if update.Message == nil {
			continue
		}

		// Mesaj yapısını oluşturma
		msgStruct := Message{
			Username: update.Message.From.UserName,
			UserID:   update.Message.From.ID,
			Text:     update.Message.Text,
			GroupID:  update.Message.MediaGroupID,
		}

		// JSON formatına dönüştürme
		msgJSON, err := json.Marshal(msgStruct)
		if err != nil {
			log.Printf("Failed to marshal message to JSON: %v", err)
			continue
		}

		// JSON formatında dosyaya yazma
		_, err = outputFile.Write(append(msgJSON, '\n'))
		if err != nil {
			log.Printf("Failed to write message to output file: %v", err)
		}

		/* //! BURADA MESAJA CEVAP OLUŞTURUYOR VE GÖNDERİYOR GEREK YOK
		// Cevap mesajı oluşturma
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		// Mesajı gönderme
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}

		*/

	}
}
