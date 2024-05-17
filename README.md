# go-gallery-hub

## Overview
Go App Dev 2024 Spring semseter course project. CRUD API application built with Go programming language.
It offers a set of endpoints to manage data about Sanrio cahracters and related TV shows. Below are the supported operations

```
Create Character: POST /characters
Get Character by ID: GET /characters/:id
Get Character List: GET /characters
Update Image: PUT /characters/:id
Delete Image: DELETE characters/:id

Create Show: POST /shows
Get Show by ID: GET /shows/:id
Update Show: PUT /shows/:id
Delete Show: DELETE shows/:id
List Shows: GET /shows

List Characters of Show: GET /shows/:id/characters

Register User: POST /users
Activate User: PUT /users/activated
Create Authentication Token: POST /tokens/authentication
```



## Database structure

```
Table characters {
  id bigserial [primary key]
  name text
  personality text
  description text
  hobbies text
  affiliations []text
}

Table shows {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  title text
  description text
  realeasedate timestamp
}

Table show_characters {
  id bigserial [primary key]
  show_id bigserial
  character_id bigserial
}

Ref: show_characters.character < characters.id
Ref: show_characters.shows < shows.id
```