# test-songs-api

This is a RESTful API service on Go for managing a music database. It allows you to add, receive, update and delete information about songs. Main functions:

Storing information about songs in a PostgreSQL database, including title, artist, release date, text, and link.

    Filtering and pagination of the song list by various parameters (band, name, release date, etc.).

    Integration with an external API to get additional information about a song when added.

    Logging using logrus to track errors and events.

    Using environment variables via godotenv for database and API configuration.

The API is developed using Gin and provides convenient endpoints for working with the music library.

Here are the possibilities of the project:

`
[GIN-debug] GET    /songs                    
[GIN-debug] POST   /songs                    
[GIN-debug] GET    /songs/:id             
[GIN-debug] PUT    /songs/:id               
[GIN-debug] DELETE /songs/:id              
[GIN-debug] GET    /info                     
[GIN-debug] GET    /swagger/*any
 `
link to the testing of this application

`http://localhost:8080/swagger/index.html#/`
