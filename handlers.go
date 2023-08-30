package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type matchHandler struct {
	matches map[string]*Match
}

// SCENARIOS that this handler will handle
// create game, join game, start game, mutate game , get game state
//
//	w.Write([]byte("The time is: "))
func (matchHandler matchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body1, err1 := io.ReadAll(r.Body)

	if err1 != nil {
		fmt.Println(err1.Error())

	} else {
		var request request

		err := json.Unmarshal(body1, &request)

		if err != nil {
			fmt.Println(err.Error())
		}

		if request.category == "create" {

			matchHandler.matches[request.gameid] = &Match{
				id:                     request.gameid,
				password:               request.gamepassword,
				founder_password:       request.playerpassword,
				stage:                  "Waiting for game to start",
				substage:               0,
				players:                make(map[string]*Player),
				waitingfor:             request.name,
				currentaction:          "Players Assembling",
				playernames:            []string{request.name},
				scheduled_for_deletion: false,
			}
		}
		if matchHandler.matches[request.gameid].password == request.gamepassword {
			if request.category == "join" {

				matchHandler.matches[request.gameid].addplayer(Player{
					password: request.playerpassword,
					name:     request.name,
				})
				matchHandler.matches[request.gameid].playernames =
					append(matchHandler.matches[request.gameid].playernames, request.name)
			}
			if request.category == "start" && request.playerpassword == matchHandler.matches[request.gameid].founder_password {
				matchHandler.matches[request.gameid].LaunchGame()
			}
			if request.category == "mutate" {
				matchHandler.matches[request.gameid].central_method(request)
			}
			if request.category == "getgamestate" {
				//	match := matchHandler.matches[request.gameid]

				var e exporter

				exportJsonstruct := e.getExportJsonStruct(request, matchHandler.matches[request.gameid])
				// getExportJsonStruct(  request ,  matchHandler.matches[request.gameid] )

				reply := e.transformExportIntoJson(exportJsonstruct)

				w.Write((reply))

			}

		}
	}

}
