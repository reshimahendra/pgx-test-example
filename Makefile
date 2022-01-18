test:
	go test ./... -v -cover -short
test-cover:
	go test ./... -v -short -coverprofile=proof.out && go tool cover -html=proof.out
