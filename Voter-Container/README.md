Build and Run Using Single Docker Command:
------------------------------------------
For API: http://localhost:1080
Redis: http://localhost:8001

Build and Run Using Docker Commands:
====================================
Build the docker image:
************************
docker-compose build

Run the docker container:
**************************
docker-compose up

Redis Database : 
******************

We can use the script provided by you to make it available


Routes :
=========

get-by-id           Get the Voter list info by id
get-all				Get all voters information
add-voter          	Add the voter details into redis database
get-voter-history  	Fetch voter history for the voter id which we are giving
get-poll     		Fetch single voter poll data with id and pollId
add-poll-by-id     	Add a voter poll record for the voter based on id  
health-check       	Check about Voter's api health which is hardcoded
update-2:           Updating firstname and lastname and have to provide all the voters information with the new firstname, lastname value
delete-by-id        Deleting the voters history by using id
