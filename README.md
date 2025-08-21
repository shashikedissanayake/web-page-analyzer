# Simple web page analyzer application

## Client application

This application consists of a React application that is packaged with Docker.

### Installation guidance

#### Prerequisites

This service uses: NodeJS.

#### Installation

- go to the client directory using `cd client`
- Install application dependencies using  `npm install` or `yarn`.
- Create an `.env` file with keys described in Config Keys section

#### Running

To run application we can use one of the following commands:
- `npm run start` - will run application in development mode.
- `docker-compose up client` - will run the application in Docker container.

#### Config keys

|Key| Description | Sample | 
|--|--|--|
| REACT_APP_SERVER_URL | URL of the backend server | http://localhost:3001|

## Server application

This application consists of a Golang backend that is packaged with Docker.

### Installation guidance

#### Prerequisites

This service uses: GoLang (V1.24.6).

#### Installation

- go to client directory using `cd server`
- Install application dependencies using `go get` or  `go mod download`.
- Create an `.env` file with keys described in Config Keys section

#### Running

To run application we can use one of the following commands:
- `go run main.go` - will run application.
- `docker-compose up server` - will run the application in Docker container.

#### Config keys

|Key| Description | Sample | 
|--|--|--|
| PORT | Port of the backend server | 3001 (Docker image expose 3001 port) |

## Running applications with Docker

Docker-compose file is included in the project base directory to run both services. To use Docker-Compose file, use the following steps

#### Prerequisites

- Make sure Docker is installed in the system before running any of the below.
- Make sure to create `.env` files in both directories according to the guidance of the individual project 

  
#### Running

- Make sure to run command below on the project base directory
- `docker-compose up --build` will spin up both the client and server applications
