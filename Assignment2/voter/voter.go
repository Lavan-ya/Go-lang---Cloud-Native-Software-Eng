package voter

import (
	"errors"
	"time"
)

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
type VoterList struct {
	Voters map[uint64]Voter //A map of VoterIDs as keys and Voter structs as values
}

// constructor for VoterList struct
func NewVoter(id uint64, fn, ln string) *Voter {
	return &Voter{
		VoterID:     id,
		FirstName:   fn,
		LastName:    ln,
		VoteHistory: []voterPoll{},
	}
}

func (v *VoterList) PutVoter(VoterID uint64, item Voter) error {
	list, ok := v.Voters[item.VoterID]
	if !ok {
		return errors.New("item doesn't exists")
	}
	list.FirstName = item.FirstName
	list.LastName = item.LastName
	v.Voters[item.VoterID] = list
	return nil
}

func (v *VoterList) AddVoterlist(item Voter) error {
	_, ok := v.Voters[item.VoterID]
	if ok {
		return errors.New("item already exists")
	}
	v.Voters[item.VoterID] = item

	return nil
}

func (v *Voter) AddPoll(pollID uint64, voteDate time.Time) {

	v.VoteHistory = append(v.VoteHistory, voterPoll{PollID: pollID, VoteDate: voteDate})
}

/*
func (v *VoterList) GetVoterDetails() ([]Voter, error) {
	var details []Voter
	for _, item := range v.Voters {
		details = append(details, item)
	}
	return details, nil
}

func (v *Voter) AddPoll(pollID uint64) {
	v.VoteHistory = append(v.VoteHistory, voterPoll{PollID: pollID, VoteDate: time.Now()})

}

func (v *Voter) ToJson() string {
	b, _ := json.Marshal(v)
	return string(b)
}
*/
