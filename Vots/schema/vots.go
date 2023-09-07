package schema

type Vote struct {
	VoteID    uint
	VoterID   uint
	PollID    uint
	VoteValue uint
}

type pollOption struct {
	PollOptionID   uint
	PollOptionText string
}

type Poll struct {
	PollID       uint
	PollTitle    string
	PollQuestion string
	PollOptions  []pollOption
}
