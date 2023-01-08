# Link Shortener

This is a link shortener service that converts a provided URL into a fictional, shortened URL. The service exposes production-ready APIs and uses MongoDB Atlas cloud DB as the data storage.

## Features

* Shortens a URL into a unique, fictional URL
* Stores shortened URLs for 30 days
* Refreshes the cache for a shortened URL if it has been accessed within its 30 day limit
* Requires authentication
* Implements OAuth 2.0 for authentication

## Local Installation

* Clone the repository:
`git clone https://github.com/your-username/link-shortener.git`
* Build the Go application:
`go build`
* Run the Go application:
`./link-shortener`

## Usage

The link shortener service exposes the following APIs:

* Shorten a URL `POST /shorten`
```
Authorization: Bearer <access_token>

{
  "long_url": "https://www.example.com"
}
```
This API shortens the provided URL and returns the shortened URL.

* Get a shortened URL
`GET /{short_url}`
```
Authorization: Bearer <access_token>
```
This API retrieves the original URL for the provided shortened URL.

## Running a Docker container

To build and run a Docker container, follow these steps:

* Build the Docker image:
```
docker build -t link-shortener .
```
* Run the Docker container:
```
docker run -p 8000:8000 link-shortener
```
This will start the Docker container and bind the container's port 8080 to the host's port 8080.
