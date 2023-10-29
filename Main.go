package main

import (
	"net/http"
	"time"
)

func main() {
	//todo for this project
	// 1. we need to add dynamic rendering on the frontend and wire up all the REST api calls
	// 2. add tests as a form of sanity check to make sure everything is working
	// 3. setup some presets and secure art Assets
	// 4. setup MYSQL for the presets
	// 5. thorougly comment the code and add links to original board game documentation
	// 6. add more ATMOSPHERE
	// 7. ??
	// 8. Profit

	Testr := &tester{}
	Testr.testMatch()

	morseconverter := MorseConverter{}
	morseconverter.convertIntoMorseCode("Snape killed dumbledore")
	matches := make(map[string]*Match)
	server := http.NewServeMux()
	matchhandler := matchHandler{matches: matches}

	server.Handle("/", matchhandler)
	http.ListenAndServe(":8001", server)
	go func() {
		for {
			time.Sleep(2 * time.Hour)
			for key, match := range matches {
				if match.scheduled_for_deletion == true {
					delete(matches, key)
				} else {
					match.scheduled_for_deletion = true
				}
			}
		}
	}()

}
