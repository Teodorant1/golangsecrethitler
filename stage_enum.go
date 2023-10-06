package main

type fash_stage_enum struct{}
type game_stage_enum struct{}
type request_dot_action_enum struct{}

func (reqActEnum request_dot_action_enum) nominate_chancellor() string {
	return "nominatechancellor"
}

func (reqActEnum request_dot_action_enum) pickpolicy() string {
	return "pickpolicy"
}
func (reqActEnum request_dot_action_enum) veto() string {
	return "veto"
}
func (reqActEnum request_dot_action_enum) fascistpower() string {
	return "fascistpower"
}
func (reqActEnum request_dot_action_enum) vote() string {
	return "vote"
}

// ///////////////////////////////////////////////////////////////
func (fsg fash_stage_enum) murder() string {
	return "murder"
}
func (fsg fash_stage_enum) spyguy() string {
	return "spyguy"

}
func (fsg fash_stage_enum) spydeck() string {
	return "spydeck"

}
func (fsg fash_stage_enum) special_election() string {
	return "special_election"

}
func (fsg fash_stage_enum) veto() string {
	return "veto"
}

func (fsg fash_stage_enum) none() string {
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
