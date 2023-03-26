module github.com/mrsiano/grpc-go/security/advancedtls/examples

go 1.17

require (
	github.com/mrsiano/grpc-go v1.52.0
	github.com/mrsiano/grpc-go/examples v0.0.0-20230111003119-9b73c42daa31
	github.com/mrsiano/grpc-go/security/advancedtls v0.0.0-20230111003119-9b73c42daa31
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/mrsiano/grpc-go => ../../..

replace github.com/mrsiano/grpc-go/examples => ../../../examples

replace github.com/mrsiano/grpc-go/security/advancedtls => ../
