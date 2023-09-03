package main

import (
	"net/http"
	"time"
)

func main() {
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
