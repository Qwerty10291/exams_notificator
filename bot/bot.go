package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler func(*Bot, tgbotapi.Update)

type BotConfig struct{
	Token string
	NoCommandMessage string
	Debug bool
}

type Bot struct {
	*tgbotapi.BotAPI
	config BotConfig
	updateChan tgbotapi.UpdatesChannel
	handlers map[string]Handler
}

func NewBot(config BotConfig) (*Bot, error){
	botApi, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil{
		return nil, err
	}
	botApi.Debug = config.Debug
	
	u := tgbotapi.NewUpdate(0)
	updateChan := botApi.GetUpdatesChan(u)
	return &Bot{
		BotAPI: botApi,
		config: config,
		updateChan: updateChan,
	}, nil
}

func (b *Bot) StartPolling() {
	for update := range b.updateChan{
		b.updateHandler(update)
	}
}

func (b *Bot) SetCommandHandler(command string, handler Handler) {
	b.handlers[command] = handler
}

func (b *Bot) updateHandler(update tgbotapi.Update) {
	if update.Message == nil{
		return
	}

	if update.Message.IsCommand(){
		if handler, ok := b.handlers[update.Message.Command()]; ok{
			handler(b, update)
		} else {
			b.Send(tgbotapi.NewMessage(update.Message.Chat.ID, b.config.NoCommandMessage))
		}
	}
}