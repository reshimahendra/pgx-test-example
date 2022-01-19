run:
	go run main.go
build:
	go build -o server -ldflags '-s -w' main.go
test:
	go test ./... -v -cover -short
test-cover:
	go test ./... -v -short -coverprofile=proof.out && go tool cover -html=proof.out
show-test-cover:
	go tool cover -html=proof.out

