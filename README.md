# RESTful BookServer HTTP API server using [Go](https://github.com/golang), [Cobra CLI](https://github.com/spf13/cobra), [Go-chi](https://github.com/go-chi/chi)
A simple api server for learning purpose.

## Running the server ##
### Running the server from direct source code ##
- `git clone https://github.com/souravbiswassanto/api-bookserver.git`
- `cd api-bookserver`
- `go build .`
- `./api-bookserver start . or ./api-bookserver start -p 3000`

### Running the server from docker image ###
- `docker pull souravbiswassanto/bookserver`
- `docker run -dp <choosen port>:8081 souravbiswassanto/bookserver`
---------------
