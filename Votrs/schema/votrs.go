package schema

import "time"

type voterPoll struct {
	PollID   uint64
	VoteDate time.Time
}

type Voter struct {
	VoterID     uint64
	FirstName   string
	LastName    string
	VoteHistory []voterPoll
}
