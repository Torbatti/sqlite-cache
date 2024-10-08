package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Users struct which contains
// an array of users
// type Games struct {
// 	Games []Game `json:"games"`
// }

// User struct which contains a name
// a type and a list of social links
type JsonGame struct {
	Game          string `json:"Game"`
	Year          int    `json:"Year"`
	Dev           string `json:"Dev"`
	DevLink       string `json:"DevLink"`
	Publisher     string `json:"Publisher"`
	PublisherLink string `json:"PublisherLink"`
	GameLink      string `json:"GameLink"`
	Platform      string `json:"Platform"`
	PlatformLink  string `json:"PlatformLink"`
}

func OpenParseJson(file_path string) []JsonGame {
	// Open our jsonFile
	jsonFile, err := os.Open(file_path)

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened :", file_path)

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	var jsonGames []JsonGame
	jsonByte, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(jsonByte, &jsonGames)

	// var games []models.Game
	// for i := 0; i < len(jsonGames); i++ {
	// 	games = append(games, models.Game{
	// 		Game:      jsonGames[i].Game,
	// 		Year:      sql.NullString{String: string(jsonGames[i].Year)},
	// 		Dev:       sql.NullString{String: jsonGames[i].Dev},
	// 		Publisher: sql.NullString{String: jsonGames[i].Publisher},
	// 		Platform:  sql.NullString{String: jsonGames[i].Platform},
	// 	})
	// 	println(games[i].Game)
	// }

	return jsonGames
}
