To compile the proto file
``
protoc -I=. --gofast_out=plugins=grpc:. ./example.proto
``