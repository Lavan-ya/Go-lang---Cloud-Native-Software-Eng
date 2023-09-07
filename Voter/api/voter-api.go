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
	voterList *voter.VoterList
}

func NewVoterApi() *VoterApi {
	dbHandler, err := voter.NewVoter()
	if err != nil {
		return nil
	}

	return &VoterApi{voterList: dbHandler}
}

//done
func (v *VoterApi) DeleteVoter(ctx *gin.Context) {
	VoterId, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := v.voterList.DeleteItem(VoterId); err != nil {
		log.Println("Error deleting item: ", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
	
}

//done
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

//done
func (v *VoterApi) InsertPoll(ctx *gin.Context) {
	VoterID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	
	var requestBody struct {
		PollID   uint64 `json:"poll_id"`
		VoteDate time.Time `json:"vote_date"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
/*
	voteDate, err := time.Parse("2006-01-02T15:04:05Z", requestBody.VoteDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}*/
    v.voterList.AddTopoll(VoterID,requestBody.PollID,requestBody.VoteDate)

}

//done
func (v *VoterApi) PostVoter(c *gin.Context) {
	var list voter.Voter
	if err := c.ShouldBindJSON(&list); err != nil {
		fmt.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
fmt.Println("Value is -------------------------",list)
	if err := v.voterList.AddVoterlist(list); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

//done
func (v *VoterApi) GetVoterJson(ctx *gin.Context) {
	voterID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	
	voter,err := v.voterList.GetItem(voterID)
	if err != nil {
		log.Println("Item not found: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, voter)
}

func (v *VoterApi) GetVoterListJson(ctx *gin.Context) {
	voter,err := v.voterList.GetFullItem()
	if err != nil {
		log.Println("Item not found: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(voter)
	ctx.Header("Content-Type", "application/json")
	ctx.Data(http.StatusOK, "application/json", b)
}

//done
func (v *VoterApi) GetVoterHistory(ctx *gin.Context) {
	voterID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	voter,err := v.voterList.GetVoterHistoryItem(voterID)
	if err != nil {
		log.Println("Item not found: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, voter)
}

//done
func (v *VoterApi) GetVoterPoolid(ctx *gin.Context) {
	voterID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	PollID, _ := strconv.ParseUint(ctx.Param("pollid"), 10, 64)
	voter,err := v.voterList.GetVoterPoolidItem(voterID,PollID)
	if err != nil {
		log.Println("Item not found: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, voter)
}

func (v *VoterApi) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"Status":  "healthy",
		"Message": "API is working properly",
		"Version": 1.1,
		"port":    1080,
	})
}

