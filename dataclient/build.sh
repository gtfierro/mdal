GOPATH=/home/gabe/go
python -m grpc_tools.protoc -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I. --python_out=. --grpc_python_out=. mdal.proto

