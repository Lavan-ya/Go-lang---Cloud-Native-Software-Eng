#!/bin/bash
VAR=${1:-localhost}    

#delete previous db if it exists and install jq
apt-get -y install jq

#delete the database
redis-cli -h $1 flushdb

#load polls
cat /data/poll.json | jq -c '.[]' |\
    while read json_object; do \
        pollid=$(jq -r '.PollID' <<< $json_object); \
        #echo $pubid  \
        rediscmd="redis-cli -h $1 JSON.set polls:$pollid . '$json_object'"; \
        echo $rediscmd; \
        eval $rediscmd; \
    done 

#load voter list
cat /data/voter.json | jq -c '.[]' |\
    while read json_object; do \
        voterid=$(jq -r '.VoterID' <<< $json_object); \
        #echo $pubid  \
        rediscmd="redis-cli -h $1 JSON.set voter:$voterid . '$json_object'"; \
        echo $rediscmd; \
        eval $rediscmd; \
    done 
#load votes list
cat /data/votes.json | jq -c '.[]' |\
    while read json_object; do \
        voteid=$(jq -r '.VoteID' <<< $json_object); \
        #echo $pubid  \
        rediscmd="redis-cli -h $1 JSON.set votes:$voteid . '$json_object'"; \
        echo $rediscmd; \
        eval $rediscmd; \
    done 

