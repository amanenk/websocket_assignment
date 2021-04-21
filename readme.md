## Usage

### Start Server
```
go run main.go
```
or run it in the container
```
docker build -t websocket_test_app .
docker run -p 5000:5000 websocket_test_app
```
### Start Client
```
go run ./client/client.go 
```
or
```
go run ./client/client.go  {ws url} {num of connections}
go run ./client/client.go  ws://localhost:5000/ws 100
```
