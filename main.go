package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	TelegramToken string `yaml:"telegram_token"`
	GameShortName string `yaml:"game_short_name"`
	Port          string `yaml:"port"`
	GameURL       string `yaml:"game_url"`
}

// loadConfig reads and parses the configuration from config.yaml
func loadConfig() (Config, error) {
	var cfg Config

	// Check if config file exists
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		return cfg, fmt.Errorf("config.yaml not found")
	}

	// Read config file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return cfg, fmt.Errorf("error reading config: %v", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("error parsing YAML: %v", err)
	}

	return cfg, nil
}

func main() {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		log.Println("Falling back to environment variables")

		config = Config{
			TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
			GameShortName: os.Getenv("GAME_SHORT_NAME"),
			Port:          os.Getenv("PORT"),
			GameURL:       os.Getenv("GAME_URL"),
		}
		if config.Port == "" {
			config.Port = "8080"
		}
		if config.GameURL == "" {
			config.GameURL = "https://kuvaev.me/telegame/"
		}
	}

	var bot *tgbotapi.BotAPI

	if config.TelegramToken != "" {
		bot, err = tgbotapi.NewBotAPI(config.TelegramToken)
		if err != nil {
			log.Printf("Error initializing Telegram bot: %v", err)
		} else {
			log.Printf("Authorized on account %s", bot.Self.UserName)

			// Start polling for updates
			u := tgbotapi.NewUpdate(0)
			u.Timeout = 60
			updates := bot.GetUpdatesChan(u)

			// Handle updates in a goroutine
			go func() {
				for update := range updates {
					if update.Message != nil && update.Message.IsCommand() {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
						switch update.Message.Command() {
						case "start":
							msg.Text = "Welcome to the Telegram game bot!"
							msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonURL("Play now", config.GameURL),
								),
							)
						default:
							msg.Text = "Unknown command"
						}
						bot.Send(msg)
					}
				}
			}()
		}
	} else {
		log.Println("TELEGRAM_TOKEN not set, bot functionality disabled")
	}

	// Set up HTTP server
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", handleRoot)

	// Create server
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s...\n", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

// handleRoot handles the root endpoint
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Telegram Game Backend is running!")
}
