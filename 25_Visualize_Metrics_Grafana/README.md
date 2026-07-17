# Lab Exercise 9.1: Visualize KEDA AutoScaling Metrics in Grafana

This lab covers configuring KEDA to export autoscaling metrics in Prometheus format, scraping them using PodMonitors via the Prometheus Operator, and visualizing them in a Grafana dashboard.

---

## 🏗️ Architecture & Metrics Pipeline Flow

```mermaid
graph TD
    subgraph default Namespace
        Consumer[Deployment: consumer-program]
        Secret[Secret: keda-rabbitmq-secret]
        ScaledObject[ScaledObject: keda-rabbitmq]
    end

    subgraph keda Namespace
        KEDA_Operator[KEDA Operator]
        KEDA_Metrics[KEDA Metrics Server]
    end

    subgraph monitoring Namespace
        Prometheus[Prometheus Operator / Server]
        Grafana[Grafana Dashboard]
        PodMonitor_Op[PodMonitor: keda-operator]
        PodMonitor_Webhook[PodMonitor: keda-admission-webhooks]
        PodMonitor_API[PodMonitor: keda-operator-metrics-apiserver]
    end

    KEDA_Operator -->|Polls Queue Length| RMQ[(RabbitMQ Cluster)]
    KEDA_Operator -->|Scales Replicas| Consumer

    %% Scraping Paths
    Prometheus -->|Scrapes via PodMonitor_Op| KEDA_Operator
    Prometheus -->|Scrapes via PodMonitor_API| KEDA_Metrics
    Grafana -->|Queries Metrics| Prometheus
    Grafana -->|Visualizes Metrics| User((User))
```

---

## Prerequisites

1. Kubernetes cluster with KEDA installed.
2. Prometheus Operator (kube-prometheus-stack) installed in the `monitoring` namespace.
3. RabbitMQ Cluster running in the `rabbitmq` namespace.

---

## 📂 Manifests & Configuration Files

All required files are already prepared in this directory:
1. `secret.yaml`: Configures the base64 connection string for RabbitMQ.
2. `consumer.yaml`: Deploys the consumer deployment and script configmap.
3. `producer.yaml`: Job used to publish test messages.
4. `pod-monitors.yaml`: PodMonitor resources targeting KEDA components.
5. `scaledobject.yaml`: ScaledObject referencing RabbitMQ scaler.
6. `keda-dashboard.json`: Pre-built KEDA Grafana Dashboard JSON ready to be imported.

---

## 🛠️ Step-by-Step Lab Walkthrough

### 1. Enable KEDA Metrics Exporting
To enable KEDA to export metrics in Prometheus format, update the KEDA Helm release:
```bash
helm upgrade -i keda kedacore/keda \
  --namespace keda \
  --set prometheus.metricServer.enabled=true \
  --set prometheus.operator.enabled=true \
  --set prometheus.webhooks.enabled=true
```

---

### 2. Configure Prometheus PodMonitors
Apply the `pod-monitors.yaml` manifest so Prometheus begins scraping metrics from the KEDA pods:
```bash
kubectl apply -f pod-monitors.yaml
```

Verify that KEDA targets are discovered by Prometheus:
1. Port-forward the Prometheus dashboard:
   ```bash
   kubectl port-forward service/prometheus-stack-kube-prom-prometheus -n monitoring 9090:9090
   ```
2. Navigate to [http://localhost:9090/targets?search=keda](http://localhost:9090/targets?search=keda) to verify that `keda-operator`, `keda-admission-webhooks`, and `keda-operator-metrics-apiserver` show up.

---

### 3. Setup Grafana Dashboard
1. Port-forward the Grafana service:
   ```bash
   kubectl port-forward service/prometheus-stack-grafana -n monitoring 3000:80
   ```
2. Access Grafana at [http://localhost:3000](http://localhost:3000) using the credentials:
   * **Username:** `admin`
   * **Password:** `prom-operator`
3. Navigate to **Dashboards** -> **New** -> **Import**.
4. Upload or copy-paste the contents of the `keda-dashboard.json` file in this directory, select `Prometheus` as the datasource, and click **Import**.

---

### 4. Deploy Workload & ScaledObject
1. Deploy the secret and consumer app:
   ```bash
   kubectl apply -f secret.yaml
   kubectl apply -f consumer.yaml
   ```

2. Deploy the KEDA ScaledObject:
   ```bash
   kubectl apply -f scaledobject.yaml
   ```

3. Confirm that the ScaledObject is ready:
   ```bash
   kubectl get scaledobjects.keda.sh keda-rabbitmq
   ```
   *Expected Output:*
   ```text
   NAME            SCALETARGETKIND      SCALETARGETNAME    MIN   MAX   READY   ACTIVE    FALLBACK   PAUSED   TRIGGERS   AUTHENTICATIONS                  AGE
   keda-rabbitmq   apps/v1.Deployment   consumer-program   1     10    True    False     False      Unknown  rabbitmq   keda-trigger-auth-rabbitmq-conn  12s
   ```

---

### 5. Generate Messages & Observe AutoScaling
1. Publish 20 messages to the queue:
   ```bash
   sed 's/value: "1"/value: "20"/' producer.yaml | kubectl apply -f -
   ```

2. Monitor the autoscaling behavior inside the Grafana KEDA Dashboard.
   * **Namespace:** `default`
   * **ScaledObject:** `keda-rabbitmq`
   * **Scaler:** `rabbitMQScaler`
   * **Metric:** `s0-rabbitmq-testqueue`
   * Set time range to **Last 5 minutes**.

You should see the metrics value scale up rapidly, and KEDA expanding the replica count accordingly!

---

## 🧹 Clean Up

To clean up all resources created in this exercise:
```bash
kubectl delete -f scaledobject.yaml
kubectl delete -f consumer.yaml
kubectl delete -f secret.yaml
kubectl delete -f pod-monitors.yaml
kubectl delete job rabbitmq-producer || true
```