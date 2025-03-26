# test-songs-api

This is a RESTful API service on Go for managing a music database. It allows you to add, receive, update and delete information about songs. Main functions:

Storing information about songs in a PostgreSQL database, including title, artist, release date, text, and link.

Filtering and pagination of the song list by various parameters (band, name, release date, etc.).

Integration with an external API to get additional information about a song when added.

Logging using logrus to track errors and events.

Using environment variables via godotenv for database and API configuration.

The API is developed using Gin and provides convenient endpoints for working with the music library.

Here are the possibilities of the project:

```
[GIN-debug] GET    /songs

[GIN-debug] POST   /songs

[GIN-debug] GET    /songs/:id
   
[GIN-debug] PUT    /songs/:id
   
[GIN-debug] DELETE /songs/:id
  
[GIN-debug] GET    /info
    
[GIN-debug] GET    /swagger/*any
 ```
In order to test this implementation, you need to run the `go run main.go` command and follow the link `http://localhost:8080/swagger/index.html #/`

the test task itself looks like this:

Implementing an online song library

It is necessary to implement the following

1. Set up rest methods

- Getting library data with filtering by all fields and
pagination

- Getting the lyrics of a song with pagination by verses

- Deleting a song

- Changing song data

- Adding a new song in the format

```JSON
{
"group": "Muse",
"song": "Supermassive Black Hole"
}
```

2. When adding, make a request to the API described by swagger. The API
described by swagger will be raised when checking the test assignment.
You don't need to implement it separately.

```
openapi: 3.0.3
info:
   title: Music info
   version: 0.0.1
paths:
   /info:
   get:
      parameters:
         - name: group
         in: query
         required: true
         schema:
            type: string
         - name: song
         in: query
         required: true
         schema:
            type: string
      responses:
      '200':
         description: Ok
         content:
            application/json:
               schema:
                  $ref: '#/components/schemas/SongDetail'
            '400':
               description: Bad request
            '500':
               description: Internal server error
components:
   schemas:
      SongDetail:
         required:
            - releaseDate
            - text
            - link
         type: object
         properties:
            releaseDate:
               type: string
               example: 16.07.2006
            text:
               type: string
               example: Ooh baby, don't you know I suffer?\nOoh baby, can
you hear me moan?\nYou caught me under false pretenses\nHow long
before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set
my soul alight
            link:
               type: string
               example: https://www.youtube.com/watch?v=Xsp3_a-PMTw
```

3. Put the enriched information in the postgres database (the database structure should
be created by migrations at the start of the service)

4. Cover the code with debug and info logs

5. Export the configuration data to the .env file

6. Generate a swagger for the implemented API
