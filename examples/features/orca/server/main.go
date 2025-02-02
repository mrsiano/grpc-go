/*
 *
 * Copyright 2023 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Binary server is an example server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mrsiano/grpc-go/codes"
	"github.com/mrsiano/grpc-go/internal"
	"github.com/mrsiano/grpc-go/orca"
	"github.com/mrsiano/grpc-go/status"

	pb "github.com/mrsiano/grpc-go/examples/features/proto/echo"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	// Report a sample cost for this query.
	cmr := orca.CallMetricRecorderFromContext(ctx)
	if cmr == nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve call metric recorder (missing ORCA ServerOption?)")
	}
	cmr.SetRequestCost("db_queries", 10)

	return &pb.EchoResponse{Message: in.Message}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Printf("Server listening at %v\n", lis.Addr())

	// Create the gRPC server with the orca.CallMetricsServerOption() option,
	// which will enable per-call metric recording.
	s := grpc.NewServer(orca.CallMetricsServerOption())
	pb.RegisterEchoServer(s, &server{})

	// Register the orca service for out-of-band metric reporting, and set the
	// minimum reporting interval to 3 seconds.  Note that, by default, the
	// minimum interval must be at least 30 seconds, but 3 seconds is set via
	// an internal-only option for illustration purposes only.
	opts := orca.ServiceOptions{MinReportingInterval: 3 * time.Second}
	internal.ORCAAllowAnyMinReportingInterval.(func(so *orca.ServiceOptions))(&opts)
	orcaSvc, err := orca.Register(s, opts)
	if err != nil {
		log.Fatalf("Failed to register ORCA service: %v", err)
	}

	// Simulate CPU utilization reporting.
	go func() {
		for {
			orcaSvc.SetCPUUtilization(.5)
			time.Sleep(2 * time.Second)
			orcaSvc.SetCPUUtilization(.9)
			time.Sleep(2 * time.Second)
		}
	}()

	s.Serve(lis)
}
