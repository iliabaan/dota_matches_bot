package stratz

import (
	"dota_matches_bot/pkg/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetLastMatch(playerID int64) (model.Match, error) {
	payload := strings.NewReader(
		fmt.Sprintf(
			"{\"query\":\"query LastsMatch($id: Long!) {\\n  player(steamAccountId: $id) {\\n    steamAccount {\\n      name\\n    }\\n    matches(request: {take: 1}) {\\n      id\\n      durationSeconds\\n      parsedDateTime\\n      analysisOutcome\\n      gameMode\\n      startDateTime\\n      averageRank\\n      players(steamAccountId: $id) {\\n        isVictory\\n        kills\\n        hero {\\n          displayName\\n        }\\n        deaths\\n        assists\\n        imp\\n        networth\\n      }\\n    }\\n  }\\n}\\n\",\"variables\":{\"id\":%d}}", playerID))

	var lastMatch model.Match

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://api.stratz.com/graphql?key="+os.Getenv("STRATZ_TOKEN_EXTRA"), payload)
	//req.Header.Set("Authorization", "Bearer "+os.Getenv("STRATZ_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return lastMatch, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return lastMatch, err
	}

	if res.StatusCode != 200 {
		if res.StatusCode == 429 {
			// need to change token for some time, then put it back
			//os.Setenv("STRATZ_TOKEN", os.Getenv("STRATZ_TOKEN_EXTRA"))
			//time.NewTimer(time.Hour)
			return lastMatch, errors.New("API rate limit exceeded")
		} else {
			return lastMatch, errors.New("error while getting data from Stratz")
		}
	}
	if err := json.Unmarshal(body, &lastMatch); err != nil {
		return lastMatch, err
	}

	lastMatch.PlayerID = playerID
	return lastMatch, nil
}

func GetGameMode(rawMode string) string {
	var modes = map[string]string{
		"TURBO":                  "турбо",
		"ALL_PICK_RANKED":        "ранкед",
		"ALL_PICK":               "all pick",
		"CAPTAINS_MODE":          "cm",
		"RANDOM_DRAFT":           "random draft",
		"SINGLE_DRAFT":           "single draft",
		"ALL_RANDOM":             "all random",
		"INTRO":                  "intro",
		"THE_DIRETIDE":           "Темный праздник",
		"REVERSE_CAPTAINS_MODE":  "reverse cm",
		"THE_GREEVILING":         "Гривиллинг",
		"TUTORIAL":               "tutorial",
		"MID_ONLY":               "соло мид",
		"LEAST_PLAYED":           "самый редкий",
		"NEW_PLAYER_POOL":        "новые игроки",
		"COMPENDIUM_MATCHMAKING": "компендиум",
		"CUSTOM":                 "кастомка",
		"CAPTAINS_DRAFT":         "отбор капитанов",
		"BALANCED_DRAFT":         "сбалансированный отбор",
		"ABILITY_DRAFT":          "ability draft",
		"EVENT":                  "событие",
		"ALL_RANDOM_DEATH_MATCH": "all random deathmatch",
		"SOLO_MID":               "соло мид",
		"MUTATION":               "мутация",
		"UNKNOWN":                "неизвестный",
	}

	mode, ok := modes[rawMode]
	if !ok {
		return rawMode
	}

	return mode
}

func GetOutcome(outcome string) string {
	var outcomes = map[string]string{
		"NONE":       "",
		"STOMPED":    " стомпом",
		"COMEBACK":   " с камбэком",
		"CLOSE_GAME": " на равных",
	}

	o, ok := outcomes[outcome]

	if !ok {
		return ""
	}

	return o
}
