package main

type fash_power_enum struct{}
type game_stage_enum struct{}
type request_dot_action_enum struct{}

// THIS targets a player
func (reqActEnum request_dot_action_enum) nominate_chancellor() string {
	return "nominatechancellor"
}

// THIS targets a policy
func (reqActEnum request_dot_action_enum) pickpolicy() string {
	return "pickpolicy"
}

// THIS doesn't target
func (reqActEnum request_dot_action_enum) veto() string {
	return "veto"
}

// THIS can target a player under certain conditions
func (reqActEnum request_dot_action_enum) fascistpower() string {
	return "fascistpower"
}

// THIS targets a player
func (reqActEnum request_dot_action_enum) vote() string {
	return "vote"
}

// ///////////////////////////////////////////////////////////////

// this targets a player
func (fsg fash_power_enum) murder() string {
	return "murder"
}

// this targets a player
func (fsg fash_power_enum) spyguy() string {
	return "spyguy"

}

// doesn't target
func (fsg fash_power_enum) spydeck() string {
	return "spydeck"

}

// lets you target a player
func (fsg fash_power_enum) special_election() string {
	return "special_election"

}

func (fsg fash_power_enum) veto() string {
	return "veto"
}

func (fsg fash_power_enum) none() string {
	return "none"
}

///////////////////////////////////////////////////////////

func (gse game_stage_enum) election() string {
	return "election"
}

//	func (gse game_stage_enum) specialelection() string {
//		return "specialelection"
//	}
func (gse game_stage_enum) policy() string {
	return "policy"
}
func (gse game_stage_enum) fascistpower() string {
	return "fascist power"
}
