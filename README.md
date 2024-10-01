# A test task from Effective Mobile

### API

To view the server API, click on the link: http://localhost:8080/swagger/index.html.

### Launching the application in Goland:
- `docker run --name my-postgres -e POSTGRES_PASSWORD=12345 -p 5432:5432 -d postgres`
- In the cmd/main.go file, line 38 should be uncommented out, and line 39 should be commented out;
- `go build cmd/main.go`
- `./main`