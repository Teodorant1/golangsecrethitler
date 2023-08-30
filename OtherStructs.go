package main

import (
	"encoding/json"
	"fmt"
)

type exportJson struct {
	Playernames       []string
	Libdecs           int
	Fashdecs          int
	Policies          int
	Discardedpolicies int
	President         string
	Chancellor        string
	Fashpowers        []string
	Intel             []string
	Votes             []string
	RolePolicies      []string
	WaitingFor        string
	CurrentAction     string
	Stage             string
	Substage          int
	VetoEnabled       bool
}
type exporter struct {
}

func (exporter *exporter) getExportJsonStruct(request request, match *Match) exportJson {

	_, _, playernames2 := match.getPresidents(match.election.lastnormalpresident)

	var rolepolicies []string = nil
	var votes []string = nil

	if !(match.stage == match.game_stage_enum.election()) {
		for i := range votes {
			votes[i] = votes[i] + match.players[votes[i]].votedFor
		}
	}

	if request.name == match.president.name {
		rolepolicies = match.president.policies
	}
	if request.name == match.chancellor.name {
		rolepolicies = match.chancellor.policies
	}
	return exportJson{
		Playernames:       playernames2,
		Libdecs:           match.libDecs,
		Fashdecs:          match.fashDecs,
		Policies:          len(match.policies),
		Discardedpolicies: len(match.discardedpolicies),
		President:         match.president.name,
		Chancellor:        match.chancellor.name,
		Fashpowers:        match.FashPowers,
		Intel:             match.players[request.name].intel,
		Votes:             votes,
		RolePolicies:      rolepolicies,

		WaitingFor:    match.waitingfor,
		CurrentAction: match.currentaction,

		Stage:       match.stage,
		Substage:    match.substage,
		VetoEnabled: match.vetoEnabled,
	}
}

func (exporter *exporter) transformExportIntoJson(exportJson exportJson) []uint8 {
	reply, err := json.Marshal(exportJson)

	if err != nil {
		fmt.Println(err.Error())
	}

	return reply

}

type election struct {
	totalvotes          int
	ja                  int
	nein                int
	failedelections     int
	specialelection     bool
	lastnormalpresident string
}

type Player struct {
	password      string
	name          string
	party         string
	isHitler      bool
	isAlive       bool
	hasVoted      bool
	votedFor      string
	Vetoing       bool
	isTermLimited bool
	intel         []string
	policies      []string
}

type request struct {
	playerpassword string
	gamepassword   string
	name           string
	gameid         string
	action         string
	target         string
	policyindex    int
	category       string
}
