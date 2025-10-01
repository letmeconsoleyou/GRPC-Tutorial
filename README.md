**Execute the proto file:**
protoc --go_out=. --go-grpc_out=. path/to/yourfile.proto

**Install protobuf (to compile proto file):**
sudo apt install -y protobuf-compiler
protoc --version

**Install go plugins:**
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

**Export path:**
export PATH="$PATH:$(go env GOPATH)/bin"
