package model

type Match struct {
	PlayerID int64
	Data     struct {
		Player struct {
			SteamAccount struct {
				Name string `json:"name"`
			} `json:"steamAccount"`
			Matches []struct {
				ID              int64       `json:"id"`
				DurationSeconds int         `json:"durationSeconds"`
				AnalysisOutcome string      `json:"analysisOutcome"`
				GameMode        string      `json:"gameMode"`
				StartDateTime   int         `json:"startDateTime"`
				ParsedDateTime  int         `json:"parsedDateTime"`
				AverageRank     interface{} `json:"averageRank"`
				Players         []struct {
					IsVictory bool `json:"isVictory"`
					Kills     int  `json:"kills"`
					Hero      struct {
						DisplayName string `json:"displayName"`
					} `json:"hero"`
					Deaths   int `json:"deaths"`
					Assists  int `json:"assists"`
					Imp      int `json:"imp"`
					Networth int `json:"networth"`
				} `json:"players"`
			} `json:"matches"`
		} `json:"player"`
	} `json:"data"`
}

// id
// match_id
// radiant_won
// game_mode
// duration
//
