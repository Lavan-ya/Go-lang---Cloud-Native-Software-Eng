package pollapi

import (
	"Assignment2/poll"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PollApi struct {
	pollList *poll.PollList
}

func NewPollApi() *PollApi {
	dbHandler, err := poll.NewPoll()
	if err != nil {
		return nil
	}
	return &PollApi{pollList: dbHandler}
}

func (v *PollApi) PostVoter(c *gin.Context) {
	var list poll.Poll
	if err := c.ShouldBindJSON(&list); err != nil {
		fmt.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
fmt.Println("Value is -------------------------",list)
	if err := v.pollList.AddItem(list); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (v *PollApi) GetVoterListJson(ctx *gin.Context) {
	voter,err := v.pollList.GetFullItem()
	if err != nil {
		log.Println("Item not found: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(voter)
	ctx.Header("Content-Type", "application/json")
	ctx.Data(http.StatusOK, "application/json", b)
}

func (v *PollApi) DeleteVoter(ctx *gin.Context) {
	VoterId, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := v.pollList.DeleteItem(VoterId); err != nil {
		log.Println("Error deleting item: ", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
	
}




