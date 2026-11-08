package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func runTelegram() error {
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("не удалось прочитать настройки из .env")
	}
	token := strings.TrimSpace(os.Getenv("JOBSCOUT_BOT_TOKEN"))
	if token == "" {
		return fmt.Errorf("не задана переменная JOBSCOUT_BOT_TOKEN")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("не удалось подключить бота: проверьте токен и доступ к Telegram")
	}

	log.Printf("Бот подключён: @%s", bot.Self.UserName)

	config := tgbotapi.NewUpdate(0)
	config.Timeout = 30

	updates := bot.GetUpdatesChan(config)
	defer bot.StopReceivingUpdates()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.Type != "private" {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		var text string

		switch update.Message.Command() {
		case "start":
			text = "Привет! Я JobScout — помощник для поиска вакансий.\n\n" +
				"Для поиска отправь:\n/search Golang\n\n" +
				"Поиск выполняется через API «Работы России»."

		case "search":
			query := strings.TrimSpace(update.Message.CommandArguments())

			if query == "" {
				text = "Укажи поисковый запрос. Например:\n/search Golang"
				break
			}

			result, err := searchTrudvsemVacancies(query)
			if err != nil {
				log.Println("Ошибка поиска через API Работы России")
				text = "Не удалось загрузить вакансии. Попробуй позже."
				break
			}

			text = result

		default:
			text = "Доступные команды:\n/start\n/search Golang"
		}

		message := tgbotapi.NewMessage(update.Message.Chat.ID, text)

		if _, err := bot.Send(message); err != nil {
			log.Println("Не удалось отправить ответ пользователю")
		}
	}

	return nil
}

func searchTrudvsemVacancies(query string) (string, error) {
	items, err := fetchTrudvsemVacancies(query)
	if err != nil {
		return "", err
	}

	vacancies := filterVacancies(items, query)

	if len(vacancies) == 0 {
		return "На полученной странице вакансий совпадений не найдено.", nil
	}

	var result strings.Builder

	result.WriteString(
		"Источник: Работа России — https://trudvsem.ru\n" +
			"Поиск по первой странице, до 20 вакансий.\n",
	)

	shown := 0

	for _, vacancy := range vacancies {
		if shown >= 5 {
			break
		}

		salary := strings.TrimSpace(vacancy.Salary)
		if salary == "" {
			salary = "не указана"
		}

		experience := "не указан"
		if vacancy.ExperienceYears != nil {
			experience = fmt.Sprintf("%d", *vacancy.ExperienceYears)
		}

		card := fmt.Sprintf(
			"\n%s\nКомпания: %s\nРегион: %s\n"+
				"Зарплата: %s\nОпыт, лет: %s\n%s\n",
			vacancy.Name,
			vacancy.Company,
			vacancy.Location,
			salary,
			experience,
			vacancy.URL,
		)

		if result.Len()+len(card) > 3500 {
			break
		}

		result.WriteString(card)
		shown++
	}

	fmt.Fprintf(
		&result,
		"\nСовпадений на странице: %d. Показано: %d.\n"+
			"Опыт и дополнительные условия не отфильтрованы — проверь карточку источника.",
		len(vacancies),
		shown,
	)

	return result.String(), nil
}
