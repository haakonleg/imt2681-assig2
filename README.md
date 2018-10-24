# IMT2681 Assignment 2

This is my submission for assignment 2 in the course IMT2681 Cloud Technologies. The task was to build a RESTful API which allows users to browse information about IGC files, using the open-source library [goigc](https://github.com/marni/goigc).

Users can add URLs to IGC resources to a database on the server and query information about added tracks. There is also webhook functionality which allows to subscribe to recieve information about newly registered tracks.

The API is deployed on Heroku here:
https://haakoleg-imt2681-assig2.herokuapp.com/

## Usage
There are two "main.go" executable files in the folder "cmd". The "paragliding" one is the main executable for serving the API. It requires two environment variables to be set:

- PORT
  - specifies which port the API is served on
- PARAGLIDING_MONGO
  - URL to a mongoDB database which will be used by the API for storing data about tracks and webhooks

The other executable "clocktrigger" is an independent executable deployed elsewhere which runs an infinite loop which checks every 10 minutes whether new tracks have been registered. If this is the case, a Slack webhook is notified and users will be notified about this.
