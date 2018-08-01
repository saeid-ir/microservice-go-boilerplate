# Necessary commands to be execute while wotking on project

# Compile all proto files in go output
compile:
	protoc --proto_path=pb --go_out=communication pb/*.proto

# Build the service into single binary file
build:
	go build -v

# Build the service into single binary file for linux
# Becuse the containers OS are linux it is better to compile service with linux binary format
build_linux:
	GOOS=linux CGO_ENABLED=0 go build -v

