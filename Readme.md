# Project:

This project is a connection between the different components.

1. Vots - Vote API
2. Votrs - Voter API
3.  Pol - Poll API

The different endpoints in Vote API - Vots:

1. GET    /votes    (get all votes)
2. GET    /vote/:id  (get vote by id)
3. GET    /votelists/:id/:idx     (get the poll -idx associated with a vote - id)

The different endpoints in Voter API - Votrs:

1. GET    /voters      (get all voters)
2. GET    /voter/:id   (get voter by id)
3. DELETE /voter/:id   (delete voter by id)

The different endpoints in Voter API - Votrs:

1. GET    /polls      (get all polls)
2. GET    /poll/:id   (get poll by id)
3. DELETE /delpoll/:id   (delete poll by id)


### To Run:

1. Move into the docker directory
2. Run `docker compose up`
2. The IP for the different endpoints will show up in the terminal

