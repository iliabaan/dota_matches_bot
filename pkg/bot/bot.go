package bot

import (
	"context"
	"dota_matches_bot/pkg/model"
	"dota_matches_bot/pkg/repository/redis"
	"dota_matches_bot/pkg/service/stratz"
	"dota_matches_bot/pkg/util"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type Bot struct {
	Bot          *tgbotapi.BotAPI
	Redis        *redis.Client
	Players      []int
	PlayersAlias map[int64]string
	ChatID       int64
}

func NewBot(api *tgbotapi.BotAPI, redis *redis.Client) *Bot {
	ids := []int{
		345007736,
		56514065,
		1110385317,
		936686976,
		334002719,
		202008401,
		182414749,
		89782335,
		93712171,
	}

	playersMap := map[int64]string{
		56514065:   "Игоряка Краполь",
		345007736:  "Илюха Фишкин",
		1110385317: "Тохинс",
		936686976:  "Михаиликус",
		334002719:  "Бустер",
		202008401:  "Дмитрий Гадюка Убеев",
		182414749:  "Мишанус",
		89782335:   "Дрюс",
		93712171:   "ЯроSlave",
	}

	chatID, _ := strconv.Atoi(os.Getenv("TG_CHAT_ID"))

	return &Bot{
		Bot:          api,
		Redis:        redis,
		Players:      ids,
		PlayersAlias: playersMap,
		ChatID:       int64(chatID),
	}
}

func (b *Bot) Handle(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	matches := make(chan model.Match)

	logrus.Println("Started listening incoming matches")

	for {
		select {
		case <-ticker.C:
			for _, id := range b.Players {
				go func(id int) {
					lm, err := stratz.GetLastMatch(int64(id))
					if err != nil {
						logrus.Errorf("error while getting last match from Stratz: %v", err)
						return
					}
					if len(lm.Data.Player.Matches) == 0 {
						return
					}
					currentLast, _ := b.Redis.GetLastMatch(ctx, id)

					if currentLast != lm.Data.Player.Matches[0].ID && lm.Data.Player.Matches[0].ParsedDateTime != 0 {
						err := b.Redis.SetLastMatch(ctx, id, lm.Data.Player.Matches[0].ID)
						if err != nil {
							logrus.Errorf("error while setting last match into Redis: %v", err)
							return
						}
						matches <- lm
					} else {
						logrus.Printf("No new matches detected for %s", lm.Data.Player.SteamAccount.Name)
					}
				}(id)
			}

		case match := <-matches:
			msg := b.createMessage(match)
			if _, err := b.Bot.Send(msg); err != nil {
				logrus.Errorf("Error while sending message via Telegram")
			}

		}
	}
}

func (b *Bot) createMessage(match model.Match) tgbotapi.MessageConfig {
	m := match.Data.Player.Matches[0]
	p := m.Players[0]
	var action string
	kda := fmt.Sprintf("%d/%d/%d", p.Kills, p.Deaths, p.Assists)

	var kdaCalculated float32
	if p.Deaths == 0 {
		kdaCalculated = float32(p.Kills + p.Assists)
	} else {
		kdaCalculated = float32(p.Kills+p.Assists) / float32(p.Deaths)
	}

	if p.IsVictory {
		action = "победил" + stratz.GetOutcome(m.AnalysisOutcome)
	} else {
		action = "проиграл" + stratz.GetOutcome(m.AnalysisOutcome)
	}

	duration := util.HumanReadableTime(m.DurationSeconds)

	msgText := fmt.Sprintf("%s сыграл <a href=\"%s\">%s</a> на %s и %s\n"+
		"KDA: %s (%.2f)\n"+
		"Нетворс: %d\n"+
		"Импакт: %d\n"+
		"Длительность матча: %s\n",
		b.PlayersAlias[match.PlayerID],
		fmt.Sprintf("https://stratz.com/matches/%d", m.ID),
		stratz.GetGameMode(m.GameMode),
		p.Hero.DisplayName,
		action,
		kda,
		kdaCalculated,
		p.Networth,
		p.Imp,
		duration)

	msg := tgbotapi.NewMessage(b.ChatID, msgText)

	msg.ParseMode = "html"

	return msg
}
