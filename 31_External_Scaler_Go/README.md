# Lab Exercise 11.1: Developing the External Scaler in Go

## Use-Case: AutoScaling Based on Custom Queue Length

Imagine you have developed a custom queuing system (like RabbitMQ, Kafka) tailored to your organization's specific workflow and requirements. This queue system is used to process various tasks, but it doesn't integrate with any existing scalers that KEDA supports out of the box.

You need to autoscale your Kubernetes deployments based on the length of your custom queue — scaling up when the queue gets too long (indicating a backlog of tasks) and scaling down when the queue length decreases (indicating less work and thus, less resource requirement). We can solve this by developing an external scaler for KEDA that communicates with your custom queue system.

This external scaler would enable you to leverage KEDA's autoscaling capabilities with your custom queue system, ensuring that your processing workloads are always scaled appropriately based on demand.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.

## Lab Environment Setup

### 1. Install Golang
*(Skip the Golang installation steps if you have completed Lab 1 or if Go is already installed on your system).*

**For Linux:**
```bash
curl -sSL -O https://go.dev/dl/go1.21.2.linux-amd64.tar.gz
sudo tar -zxf go1.21.2.linux-amd64.tar.gz -C /usr/local
export PATH=$PATH:/usr/local/go/bin
go version
```

**For Mac:**
```bash
brew install golang
```

### 2. Download Proto CLI

**For Linux:**
```bash
PROTOC_ZIP=protoc-3.14.0-linux-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP
sudo chmod +x /usr/local/bin/protoc
protoc --version
```

**For Mac:**
```bash
brew install protobuf
protoc --version
```

Congratulations, you have successfully set up the required lab environment.

## Lab Exercise

### 1. Create a Golang project

Execute the following commands to create a Golang project for custom scaler.

```bash
mkdir customscaler
cd customscaler
go mod init customscaler
```

### 2. Generate a proto stub

As discussed in the chapter, KEDA talks to scalers via GRPC protocol. We need to implement a GRPC server definition that satisfies the below methods:

```protobuf
service ExternalScaler {
  rpc IsActive(ScaledObjectRef) returns (IsActiveResponse) {}
  rpc StreamIsActive(ScaledObjectRef) returns (stream IsActiveResponse) {}
  rpc GetMetricSpec(ScaledObjectRef) returns (GetMetricSpecResponse) {}
  rpc GetMetrics(GetMetricsRequest) returns (GetMetricsResponse) {}
}
```

The commands below use a pre-defined proto file to generate grpc proto stubs for Golang.

```bash
mkdir externalscaler
cd externalscaler
wget https://raw.githubusercontent.com/kedacore/keda/main/pkg/scalers/externalscaler/externalscaler.proto
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH=$PATH:$(go env GOPATH)/bin
protoc --go_out=. --go-grpc_out=. externalscaler.proto
cd ..
```

### 3. Verify generated proto files

After executing the previous step, ensure that you have these three files in the `externalscaler` directory.

```bash
ls externalscaler
```
*Expected Output:*
```text
externalscaler.pb.go
externalscaler.proto
externalscaler_grpc.pb.go
```

### 4. Implement GRPC server stub

The code below implements the above mentioned ExternalScaler methods. In the `customscaler` directory create a file named `grpc_server.go` with the following contents:

```go
package main

import (
	"context"
	pb "customscaler/externalscaler"
	"fmt"
	"log"
	"strconv"
)

type ExternalScaler struct {
	pb.UnimplementedExternalScalerServer
}

func (es *ExternalScaler) IsActive(ctx context.Context, scaledObjectRef *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	log.Println("Executing method IsActive")
	return &pb.IsActiveResponse{Result: CustomQueueLength > 0}, nil
}

func (es *ExternalScaler) GetMetricSpec(ctx context.Context, scaledObjectRef *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {
	log.Println("Executing method GetMetricSpec")
	metricThreshold, err := strconv.ParseInt(scaledObjectRef.ScalerMetadata["queueLength"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value for metric threshold - %s", err)
	}
	return &pb.GetMetricSpecResponse{
		MetricSpecs: []*pb.MetricSpec{
			{MetricName: "custom-queue", TargetSize: metricThreshold},
		},
	}, nil
}

func (es *ExternalScaler) GetMetrics(ctx context.Context, metricRequest *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	log.Println("Executing method GetMetrics")
	if metricRequest.MetricName != "custom-queue" {
		return nil, fmt.Errorf("invalid metric name - %s", metricRequest.MetricName)
	}
	return &pb.GetMetricsResponse{
		MetricValues: []*pb.MetricValue{
			{MetricName: "custom-queue", MetricValue: CustomQueueLength},
		},
	}, nil
}

func (es *ExternalScaler) StreamIsActive(scaledObjectRef *pb.ScaledObjectRef, epsServer pb.ExternalScaler_StreamIsActiveServer) error {
	log.Println("Executing method StreamIsActive")
	return nil
}
```

### 5. Implement Custom Queue API

We don’t have an actual custom implementation of a queue, so we will be mocking the custom queue behavior via the below REST API code.

Here we have defined a variable called `CustomQueueLength` which stores the current length of the queue. We will use REST APIs to control the queue length. Whenever someone hits this API `/api/queue`, the value provided in request is set to the `CustomQueueLength`.

Also as this is a mock implementation, we have defined a function `reduceCustomQueueLength()` which runs in background and reduces the `CustomQueueLength` by one every minute. This helps us in observing the actual auto scaling behavior whenever someone sets the queue length via API.

Create a file named `api.go` with the following contents:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func setValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	number := vars["number"]
	CustomQueueLength, _ = strconv.ParseInt(number, 10, 64)
	log.Printf("new value: %d\n", CustomQueueLength)
}

func RunManagementApi() {
	r := mux.NewRouter()
	r.HandleFunc("/api/queue/{number:[0-9]+}", setValue).Methods("POST")
	http.Handle("/", r)
	fmt.Printf("Running http management server on port: %d\n", 9090)
	http.ListenAndServe(":9090", nil)
}

var CustomQueueLength int64 = 0

func reduceCustomQueueLength() {
	for {
		if CustomQueueLength > 0 {
			CustomQueueLength--
			log.Printf("Reduced queue length value: %d\n", CustomQueueLength)
			time.Sleep(1 * time.Minute)
		}
	}
}
```

### 6. Putting in all together in main.go

Create a file named `main.go` with the below contents:

```go
package main

import (
	"log"
	"net"
	pb "customscaler/externalscaler"

	"google.golang.org/grpc"
)

func main() {
	go RunManagementApi()
	go reduceCustomQueueLength()

	grpcServer := grpc.NewServer()
	lis, _ := net.Listen("tcp", ":6000")
	pb.RegisterExternalScalerServer(grpcServer, &ExternalScaler{})

	log.Println("Listening external scaler on :6000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
```

### 7. Build and verify

```bash
go mod tidy
go build && ./customscaler
```
*Expected Output:*
```text
Running http management server on port: 9090
2024/02/08 09:01:06 Listening external scaler on :6000
```

Execute the following curl request in another terminal and you will observe the following logs in the custom scaler. This output verifies our mock custom queue application is working as intended.

```bash
curl -X POST localhost:9090/api/queue/3
```
*Expected Custom Scaler Output:*
```text
2024/02/08 09:01:27 new value: 3
2024/02/08 09:08:27 Reduced queue length value: 2
2024/02/08 09:09:27 Reduced queue length value: 1
2024/02/08 09:10:27 Reduced queue length value: 0
```

## Summary

In this exercise we learned how to develop an external scaler in Go to enable autoscaling of Kubernetes deployments based on a custom queue length, not supported by KEDA out of the box. It involves creating a Golang project and setting up a GRPC server that implements the `ExternalScaler` interface, which communicates with a mocked custom queue system to dynamically scale resources according to the queue length. The process includes generating GRPC proto stubs, implementing GRPC server stub methods to interact with the custom queue, and mocking queue behavior via a REST API. The exercise concludes with building and verifying the custom scaler.