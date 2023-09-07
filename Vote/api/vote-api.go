package api

import (
	"Assignment2/vote"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VoteApi struct {
	voteList *vote.VoteList
}

func NewVoteApi() *VoteApi {
	dbHandler, err := vote.NewVote()
	if err != nil {
		return nil
	}

	return &VoteApi{voteList: dbHandler}
}

func (v *VoteApi) PostVoter(c *gin.Context) {
	var list vote.Vote
	if err := c.ShouldBindJSON(&list); err != nil {
		fmt.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println("Value is -------------------------", list)
	if err := v.voteList.AddItem(list); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	list.Links = []vote.Link{
		{Rel: "self", Href: fmt.Sprintf("/voters/%d", list.VoteID)},
		{Rel: "delete", Href: fmt.Sprintf("/voters/%d/delete", list.VoteID)},
		// ... Add more links as needed
	}
	c.JSON(http.StatusOK, list)
}

func (v *VoteApi) GetVoterListJson(ctx *gin.Context) {
	voter, err := v.voteList.GetFullItem()
	if err != nil {
		log.Println("Item not found: ", err)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	for i := range voter {
		voter[i].Links = []vote.Link{
			{Rel: "self", Href: fmt.Sprintf("/voters/%d", voter[i].VoteID)},
			{Rel: "delete", Href: fmt.Sprintf("/voters/%d/delete", voter[i].VoteID)},
			// ... Add more links as needed
		}
	}

	//b, _ := json.Marshal(voter)
	//ctx.Header("Content-Type", "application/json")
	//ctx.Data(http.StatusOK, "application/json", b)
	if ctx.DefaultQuery("format", "json") == "html" {
		htmlResp := "<ul>"
		for _, v := range voter {
			htmlResp += "<li>"
			htmlResp += fmt.Sprintf("VoterID: %d, PollID: %d", v.VoterID, v.PollID)
			for _, link := range v.Links {
				htmlResp += fmt.Sprintf(" [<a href='%s'>%s</a>]", link.Href, link.Rel)
			}
			htmlResp += "</li>"
		}
		htmlResp += "</ul>"
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlResp))
	} else {
		b, _ := json.Marshal(voter)
		ctx.Header("Content-Type", "application/json")
		ctx.Data(http.StatusOK, "application/json", b)
	}
}

func (v *VoteApi) DeleteVoter(ctx *gin.Context) {
	VoterId, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := v.voteList.DeleteItem(VoterId); err != nil {
		log.Println("Error deleting item: ", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)

}
