# Lab Exercise 10.3: Querying Metrics from KEDA Metrics Server

In this exercise we will learn how to query external metrics from KEDA Metrics Server. It exposes various metrics that are essential for understanding the behavior of your autoscaling configurations and identifying potential issues.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.
3. Completion of Lab Exercises 10.1 and 10.2.

## Lab Exercise

### 1. Recovering metric names from ScaledObject

KEDA updates each ScaledObject with essential information, including the metric names generated from its triggers. To retrieve these metric names, use:

```bash
kubectl get scaledobject keda-rabbitmq -n default -o jsonpath='{.status.externalMetricNames}'
```

*Expected Output:*
```json
["s0-rabbitmq-testqueue","s1-rabbitmq-testqueue"]
```

### 2. Querying the metric

In the command below, we are querying the external metrics of ScaledObject created in our previous lab (Lab Exercise 10.2).

```bash
kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1/namespaces/default/composite-metric?labelSelector=scaledobject.keda.sh%2Fname%3Dkeda-rabbitmq" | jq
```

*Expected Output:*
```json
{
  "kind": "ExternalMetricValueList",
  "apiVersion": "external.metrics.k8s.io/v1beta1",
  "metadata": {},
  "items": [
    {
      "metricName": "composite-metric",
      "metricLabels": null,
      "timestamp": "2024-02-07T02:56:07Z",
      "value": "15"
    }
  ]
}
```

*(Note: The `value` depends on the actual metric being reported at the time of query.)*

## Summary

In this exercise we learned how to query external metrics from the KEDA Metrics Server, which is integral for monitoring auto scaling behaviors. The exercise covers commands for querying specific metrics associated with a ScaledObject. Additionally, it demonstrates how to extract metric names directly from a ScaledObject's status, aiding in autoscaling diagnostics and configuration verification.

## Clean Up

Run the following commands from the `29_Scaling_Modifiers` directory to clean up the environment:

```bash
kubectl delete -f producer-2.yaml --ignore-not-found
kubectl delete -f secret-2.yaml --ignore-not-found
kubectl delete -f consumer.yaml --ignore-not-found
kubectl delete job -l app=producer-program --ignore-not-found
kubectl delete -f scaled-object-scaling-modifiers.yaml --ignore-not-found
kubectl delete rabbitmqcluster rabbitmq-cluster -n rabbitmq --ignore-not-found
kubectl delete rabbitmqcluster rabbitmq-cluster -n rabbitmq-2 --ignore-not-found
```