# Gin Server

Sample web server built using [Gin](https://github.com/gin-gonic/gin) (a framework for building web applications in [Go](https://golang.org))

## Getting Started

### Prerequisites
---
* Docker (v18.03+)

### Start up
---
To build and run this project locally, run the following commands at the root directory of the project
```bash
make build start-db run
```
This will build the docker image, start up a [MongoDB](https://www.mongodb.com) container for the project and then start the project at port `http://localhost:55099`

## Built With
---
* [Golang](https://golang.org) - Language used
* [Gin](https://github.com/gin-gonic/gin) - Web Framework used
* [MongoDB](https://www.mongodb.com) - Database used
  
## Documentation
---
Documentation for the Project can be found at [gin-server-doc](https://documenter.getpostman.com/view/8916756/SztG3mKJ?version=latest)