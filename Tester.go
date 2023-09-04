package main

import "fmt"

type tester struct {
}

func (tester *tester) testMatch() {

	match := &Match{
		id:               "12345",
		password:         "picard0",
		founder_password: "picard1",
		players:          make(map[string]*Player),
	}
	player1 := Player{
		password: "picard1",
		name:     "Picard1",
		//	intel:    make([]string, 0),
		//	policies: make([]string, 0),
	}
	player2 := Player{
		password: "picard2",
		name:     "Picard2",
	}
	player3 := Player{
		password: "picard3",
		name:     "Picard3"}
	player4 := Player{
		password: "picard4",
		name:     "Picard4",
	}
	player5 := Player{
		password: "picard5",
		name:     "Picard5",
	}
	player6 := Player{
		password: "picard6",
		name:     "Picard6",
	}

	//	match.lock.Unlock()

	match.addplayer(player1)
	match.addplayer(player2)
	match.addplayer(player3)
	match.addplayer(player4)
	match.addplayer(player5)
	match.addplayer(player6)

	match.LaunchGame()

	e := exporter{}
	request := request{
		playerpassword: "picard1",
		gamepassword:   "picard0",
		name:           "Picard1",
		gameid:         "12345",
		action:         "",
		target:         "",
		policyindex:    0,
		category:       "getgamestate",
	}

	exportJsonstruct := e.getExportJsonStruct(request, match)
	reply := e.transformExportIntoJson(exportJsonstruct)

	fmt.Println(string(reply))
	//	fmt.Println(match)
	//	fmt.Println(*match)

}
