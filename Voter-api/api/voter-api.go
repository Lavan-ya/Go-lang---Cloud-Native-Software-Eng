package api

import (
	"Assignment2/voter"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VoterApi struct {
	voterList voter.VoterList
}

func NewVoterApi() *VoterApi {
	return &VoterApi{
		voterList: voter.VoterList{
			Voters: make(map[uint64]voter.Voter),
		},
	}
}

func (v *VoterApi) DeleteVoter(ctx *gin.Context) {
	VoterId, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	delete(v.voterList.Voters, VoterId)
	fmt.Println("Item Deleted successfully")
}

func (v *VoterApi) UpdateVoter(ctx *gin.Context) {
	VoterId, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	var item voter.Voter

	if err := ctx.BindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	if err := v.voterList.PutVoter(VoterId, item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the voter"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Voter ID":     VoterId,
		"FirstName":    item.FirstName,
		"LastName":     item.LastName,
		"VoterHistory": item.VoteHistory,
	})
}

func (v *VoterApi) InsertPoll(ctx *gin.Context) {
	VoterID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	//PollID, _ := strconv.ParseUint(ctx.Param("pollid"), 10, 64)

	var requestBody struct {
		PollID   uint64 `json:"poll_id"`
		VoteDate string `json:"vote_date"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	voteDate, err := time.Parse("2006-01-02T15:04:05Z", requestBody.VoteDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	_, ok := v.voterList.Voters[VoterID]
	if !ok {
		fmt.Println("item doesnot exist")
		return
	}
	voter := v.voterList.Voters[VoterID]
	voter.AddPoll(requestBody.PollID, voteDate)
	v.voterList.Voters[VoterID] = voter
}

func (v *VoterApi) PostVoter(c *gin.Context) {
	var list voter.Voter
	if err := c.ShouldBindJSON(&list); err != nil {
		fmt.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.voterList.AddVoterlist(list); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (v *VoterApi) GetVoterJson(ctx *gin.Context) {
	voterID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	voter := v.voterList.Voters[voterID]
	ctx.JSON(http.StatusOK, gin.H{
		"Voter ID":     voter.VoterID,
		"First Name":   voter.FirstName,
		"Last Name":    voter.LastName,
		"Vote History": voter.VoteHistory,
	})
}

func (v *VoterApi) GetVoterListJson(ctx *gin.Context) {
	b, _ := json.Marshal(v.voterList)
	ctx.Header("Content-Type", "application/json")
	ctx.Data(http.StatusOK, "application/json", b)
}

func (v *VoterApi) GetVoterHistory(ctx *gin.Context) {
	VoterID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	voter := v.voterList.Voters[VoterID]
	ctx.JSON(http.StatusOK, gin.H{
		"Voter ID":      voter.VoterID,
		"Voter History": voter.VoteHistory,
	})
}

func (v *VoterApi) GetVoterPoolid(ctx *gin.Context) {
	VoterID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	PollID, _ := strconv.ParseUint(ctx.Param("pollid"), 10, 64)
	voter := v.voterList.Voters[VoterID]
	VoteHistory := voter.VoteHistory
	for _, poll := range VoteHistory {
		if poll.PollID == PollID {
			ctx.JSON(http.StatusOK, gin.H{
				"Poll ID":   poll.PollID,
				"Poll Date": poll.VoteDate,
			})
		}
	}

}
func (v *VoterApi) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"Status":  "healthy",
		"Message": "API is working properly",
		"Version": 1.1,
		"port":    1080,
	})
}

/*func (v *VoterApi) AddVoter(voterID uint64, firstName, lastName string) {
	v.voterList.Voters[voterID] = *voter.NewVoter(voterID, firstName, lastName)
}

func (v *VoterApi) AddPoll(voterID uint64, pollID uint64) {
	voter := v.voterList.Voters[voterID]
	voter.AddPoll(pollID)
	v.voterList.Voters[voterID] = voter
}

func (v *VoterApi) GetVoter(voterID uint) voter.Voter {
	voter := v.voterList.Voters[voterID]
	return voter
}

func (v *VoterApi) GetVoterList() voter.VoterList {
	return v.voterList
}

func (v *VoterApi) ListAllVoter(c *gin.Context) {
	voterList, err := v.voterList.GetVoterDetails()
	if err != nil {
		log.Println("Error Getting All Items: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if voterList == nil {
		voterList = make([]voter.Voter, 0)
	}
	c.JSON(http.StatusOK, voterList)
}*/
