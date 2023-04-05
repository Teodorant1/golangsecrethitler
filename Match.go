package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Match struct {
	lock                   sync.RWMutex
	id                     string
	scheduled_for_deletion bool
	password               string
	founder_password       string
	stage                  string // this will dictate what can be done
	substage               int
	president              *Player
	chancellor             *Player
	lastElected_pres       string
	lastElected_chanc      string
	hitler                 *Player
	players                map[string]*Player
	playernames            []string
	policies               []string
	discardedpolicies      []string
	FashPowers             []string
	libDecs                int
	fashDecs               int
	failedElections        int
	vetoEnabled            bool
	chancellorWantsVeto    bool
	game_stage_enum        game_stage_enum
	fash_stage_enum        fash_stage_enum
	election               election
	waitingfor             string
	currentaction          string
}

func (match *Match) central_method(request request) {
	match.lock.Lock()
	defer match.lock.Unlock()
	if match.contains_player(match.players, request.name) && match.players[request.name].isAlive == true {
		if match.players[request.name].password == request.playerpassword && match.password == request.gamepassword {
			match.central_methodv2(request)
			match.scheduled_for_deletion = false
		}
	}
}

func (match *Match) central_methodv2(request request) {
	if match.stage == match.game_stage_enum.election() && match.substage == 1 && request.action == "nominatechancellor" {
		if match.president.password == request.playerpassword && match.password == request.gamepassword {
			match.chancellor = match.players[request.target]
			match.waitingfor = "all"
			match.currentaction = "voting on the new government"
			match.substage = 2
		}
	}
	if match.substage == 2 && match.players[request.name].hasVoted == false && request.action == "vote" {
		if request.target == "ja" {
			match.election.ja++
			match.players[request.name].votedFor = "ja"
		}
		if request.target == "nein" {
			match.election.nein++
			match.players[request.name].votedFor = "nein"

		}
		match.players[request.name].hasVoted = true
		match.election.totalvotes++
		// if everyone has voted
		if match.getlivingplayers() == match.election.totalvotes {
			for _, player := range match.players {
				player.hasVoted = false
			}
			// if the election succeeds
			if match.election.ja > match.election.nein {

				for i := 0; i < 3; i++ {
					match.president.policies = append(match.president.policies, match.policies[0])
					match.policies = append(match.policies[:0], match.policies[1:]...)
				}

				match.waitingfor = match.president.name
				match.currentaction = "president is looking at top 3 cards of the deck"
				match.stage = match.game_stage_enum.policy()
				match.substage = 1

				for _, player := range match.players {
					player.isTermLimited = false
				}

				match.chancellor.isTermLimited = true
				if match.getlivingplayers() > 5 {
					match.president.isTermLimited = true
				}
				match.election.failedelections = 0
				if match.election.specialelection == false {
					match.election.lastnormalpresident = match.president.name
				}
			}
			// if the election fails
			if match.election.ja < match.election.nein {
				match.failedElections++
				_, nextpresidentIndex, playernames2 := match.getPresidents(match.election.lastnormalpresident)

				match.substage = 1
				match.stage = match.game_stage_enum.election()
				match.president = match.players[playernames2[nextpresidentIndex]]
				// if the government is collapsed
				if match.election.failedelections == 3 {
					match.collapsegovernment()
					match.election.failedelections = 0

				}
			}
			// cleanup happens here
			match.election.totalvotes = 0
			match.election.ja = 0
			match.election.nein = 0

			for _, player := range match.players {
				player.hasVoted = false
			}
		}
	}

	if match.stage == match.game_stage_enum.policy() && match.substage == 1 && match.president.name == request.name {

		if request.action == "pickpolicy" {
			match.discardedpolicies = append(match.discardedpolicies, match.president.policies[request.policyindex])
			match.president.policies = append(match.president.policies[:request.policyindex], match.president.policies[request.policyindex+1:]...)
			for _, policy := range match.president.policies {
				match.chancellor.policies = append(match.chancellor.policies, policy)
			}
			match.president.policies = nil
			match.waitingfor = match.chancellor.name
			match.currentaction = "chancellor is picking 1 from 2 policies passed to them from the president to pass as a law"
			match.substage = 2

		}

		if request.action == "veto" {
			if match.vetoEnabled == true && match.chancellorWantsVeto == true {
				match.chancellorWantsVeto = false
				for i := range match.chancellor.policies {
					match.discardedpolicies = append(match.discardedpolicies, match.policies[i])
				}
				match.chancellor.policies = nil

				_, nextpresidentIndex, playernames2 := match.getPresidents(match.election.lastnormalpresident)

				if match.election.failedelections == 3 {
					match.collapsegovernment()
					match.election.failedelections = 0

				}

				if match.election.failedelections < 3 {
					match.president = match.players[playernames2[nextpresidentIndex]]
					match.currentaction = "waiting for president to pick a chancellor"
					match.stage = match.game_stage_enum.election()
					match.substage = 1
				}
			}
		}
	}
	//	_, nextpresidentIndex, playernames2 := match.getPresidents(match.election.lastnormalpresident)
	//
	//	match.president = match.players[playernames2[nextpresidentIndex]]
	//	match.currentaction = "waiting for president to pick a chancellor"
	//	match.stage = match.game_stage_enum.election()
	//	match.substage = 1

	if match.stage == match.game_stage_enum.policy() && match.substage == 2 && match.chancellor.name == request.name {
		if request.action == "veto" {
			match.chancellorWantsVeto = true
			match.waitingfor = match.chancellor.name + " or " + match.president.name
		}
		if request.action == "pickpolicy" {
			match.chancellorWantsVeto = false

			if match.chancellor.policies[request.policyindex] == "fascist" {
				match.fashDecs++

				if match.FashPowers[match.fashDecs] != match.fash_stage_enum.none() {

					if match.FashPowers[match.fashDecs] == "spydeck" {

						intel := "Top 3 policies are : " + match.policies[0] + " , " + match.policies[1] + " , " + match.policies[2] + " . "

						match.president.intel = append(match.president.intel, intel)

						_, nextpresidentIndex, playernames2 := match.getPresidents(match.election.lastnormalpresident)
						match.president = match.players[playernames2[nextpresidentIndex]]
						match.currentaction = "waiting for president to pick a chancellor"
						match.stage = match.game_stage_enum.election()
						match.substage = 1

					} else {
						match.stage = match.game_stage_enum.fascistpower()
						match.currentaction = match.FashPowers[match.fashDecs]
						match.waitingfor = match.president.name
					}

				}
				if match.FashPowers[match.fashDecs] == match.fash_stage_enum.none() {
					_, nextpresidentIndex, playernames2 := match.getPresidents(match.election.lastnormalpresident)
					match.president = match.players[playernames2[nextpresidentIndex]]
					match.currentaction = "waiting for president to pick a chancellor"
					match.stage = match.game_stage_enum.election()
					match.substage = 1
				}

			}
			if match.chancellor.policies[request.policyindex] == "liberal" {
				match.libDecs++
				_, nextpresidentIndex, playernames2 := match.getPresidents(match.election.lastnormalpresident)
				match.president = match.players[playernames2[nextpresidentIndex]]
				match.currentaction = "waiting for president to pick a chancellor"
				match.stage = match.game_stage_enum.election()
				match.substage = 1
			}

			match.chancellor.policies = append(match.chancellor.policies[:request.policyindex], match.chancellor.policies[request.policyindex+1:]...)
			match.discardedpolicies = append(match.discardedpolicies, match.chancellor.policies[0])
			match.chancellor.policies = nil

			if len(match.policies) < 3 {
				match.policies = append(match.policies, match.discardedpolicies...)
				match.discardedpolicies = nil
				match.randomizepolicies()
			}

		}
	}
	if match.stage == match.game_stage_enum.fascistpower() && request.playerpassword == match.president.password && request.action == "fascistpower" {
		match.handleFascistPower(request)
		_, nextpresidentIndex, playernames2 := match.getPresidents(match.election.lastnormalpresident)
		match.president = match.players[playernames2[nextpresidentIndex]]
		match.currentaction = "waiting for president to pick a chancellor"
		match.stage = match.game_stage_enum.election()
		match.substage = 1
	}

}

// func (match *Match) checkIfPowerEnabled() []string {}
func (match *Match) getCandidatesForChancellor() []string {

	var playernames2 []string

	for _, player := range match.players {
		if player.isTermLimited == false {
			playernames2 = append(playernames2, player.name)
		}
	}
	return playernames2
}
func (match *Match) collapsegovernment() {
	//also gotta do the random policy here

	if match.policies[0] == "fascist" {
		match.fashDecs++
	}
	if match.policies[0] == "liberal" {
		match.libDecs++
	}
	match.policies = append(match.policies[:0], match.policies[1:]...)

	if len(match.policies) < 3 {
		match.policies = append(match.policies, match.discardedpolicies...)
		match.discardedpolicies = nil
		match.randomizepolicies()
	}

	match.waitingfor = match.president.name
	match.currentaction = "waiting for president to pick a chancellor"
	match.stage = match.game_stage_enum.election()
	match.substage = 1
}
func (match *Match) addplayer(player Player) {

	match.lock.Lock()
	defer match.lock.Unlock()
	if len(match.players) < 10 {
		match.players[player.name] = &player
	}
	fmt.Println("bazinga")
}
func (match *Match) getPresidents(president string) (int, int, []string) {
	var playernames2 []string
	for i := 0; i < len(match.playernames); i++ {
		if match.players[match.playernames[i]].isAlive == true {
			playernames2 = append(playernames2, match.playernames[i])
		}

	}
	presidentIndex := 0
	nextpresidentIndex := 1
	for i := 0; i < len(playernames2); i++ {
		if playernames2[i] == president {
			presidentIndex = i
			break
		}
	}
	if presidentIndex == len(playernames2)-1 {
		nextpresidentIndex = 1
	} else {
		nextpresidentIndex = presidentIndex + 1
	}

	return presidentIndex, nextpresidentIndex, playernames2
}

func (match *Match) getlivingplayers() int {
	var aliveNumbers int = 0
	for _, v := range match.players {
		if v.isAlive == true {
			aliveNumbers++
		}
	}
	return aliveNumbers
}

func (match *Match) handleFascistPower(request request) {

	if match.FashPowers[match.fashDecs] == match.fash_stage_enum.murder() {
		match.players[request.target].isAlive = false
	}

	if match.FashPowers[match.fashDecs] == match.fash_stage_enum.spyguy() {
		intel := "The party affiliation of " + request.target + " is " + match.players[request.target].party
		match.players[request.name].intel = append(match.players[request.name].intel, intel)
	}
	if match.FashPowers[match.fashDecs] == match.fash_stage_enum.special_election() {
		match.election.specialelection = true
		match.president = match.players[request.target]
	}

}
func (match *Match) generateFashPowers() {

	gamemodeNumber := match.calcFashNumbers()
	if gamemodeNumber == 1 {
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.none())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.none())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.none())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.spydeck())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.murder())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.murder())

	}
	if gamemodeNumber == 2 {
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.none())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.none())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.spyguy())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.special_election())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.murder())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.murder())

	}
	if gamemodeNumber == 3 {
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.none())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.spyguy())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.spyguy())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.special_election())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.murder())
		match.FashPowers = append(match.FashPowers, match.fash_stage_enum.murder())

	}

}

func (match *Match) calcFashNumbers() int {
	playersize := len(match.playernames)
	if playersize > 8 {
		return 3
	} else if playersize < 7 {
		return 1
	} else {
		return 2
	}
	return 2
}
func (match *Match) contains_player(m map[string]*Player, suspect string) bool {

	if _, ok := m[suspect]; ok {
		//	fmt.Printf("Value is : %s \n", v)
		return true
	} else {
		fmt.Println("Key not found")
		return false
	}

	return false
}

func (match *Match) randomize_players() {
	rand.Shuffle(len(match.playernames), func(i, j int) {
		match.playernames[i], match.playernames[j] = match.playernames[j], match.playernames[i]
	})
}
func (match *Match) randomizepolicies() {

	// rand.Seed()
	rand.Shuffle(len(match.policies), func(i, j int) {
		match.policies[i], match.policies[j] = match.policies[j], match.policies[i]
	})
}

func (match *Match) LaunchGame() {
	// this commented piece of code needs to be in the constructor
	//	match.players = make(map[string]*Player)
	for k := range match.players {
		match.playernames = append(match.playernames, k)
	}
	match.randomize_players()

	fashnumbers := match.calcFashNumbers()
	match.hitler = match.players[match.playernames[0]]
	match.players[match.playernames[0]].isHitler = true
	match.players[match.playernames[0]].party = "fascist"

	for helpers := 1; helpers < fashnumbers+1; helpers++ {
		match.players[match.playernames[helpers]].party = "fascist"
	}
	match.randomize_players()
	match.president = match.players[match.playernames[0]]

	for i := 0; i < 6; i++ {
		match.policies = append(match.policies, "liberal")
	}
	for i := 0; i < 11; i++ {
		match.policies = append(match.policies, "fascist")
	}

	match.randomizepolicies()
	match.waitingfor = match.president.name
	match.currentaction = "waiting for president to pick a chancellor"
	match.stage = match.game_stage_enum.election()
	match.substage = 1
	var fascists []string
	for _, player := range match.players {

		if len(match.players) > 6 {
			if player.party == "fascist" && player.isHitler == false {
				fascists = append(fascists, player.name)
			}
		}

		if len(match.players) < 7 {
			if player.party == "fascist" {
				fascists = append(fascists, player.name)
			}
		}

	}
	var intel = " The fascists are : "
	var hitlerintel = "Hitler is actually : " + match.hitler.name
	for i := range fascists {
		intel = intel + fascists[i] + " , "
	}
	for i := range fascists {
		match.players[fascists[i]].intel = append(match.players[fascists[i]].intel, intel)
		match.players[fascists[i]].intel = append(match.players[fascists[i]].intel, hitlerintel)
	}
	match.generateFashPowers()

}
