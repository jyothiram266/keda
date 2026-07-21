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
