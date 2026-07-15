# Lab Exercise 8.2: Implementing TriggerAuthentication

This exercise demonstrates how to decouple authentication parameters from a `ScaledObject` using KEDA's `TriggerAuthentication` Custom Resource Definition (CRD).

By utilizing a `TriggerAuthentication` resource, we can define how to authenticate against event sources (such as RabbitMQ) and link it to our `ScaledObject`. This separation of concerns allows the `ScaledObject` to define *what* workload to scale and *how* to scale it, while the `TriggerAuthentication` specifies *how* to connect to the trigger securely.

---

## 🏗️ Architecture & Authentication Flow

```mermaid
graph TD
    subgraph Kubernetes Cluster
        KEDA[KEDA Operator]
        Deployment[Deployment: consumer-program]
        Secret[Secret: keda-rabbitmq-secret]
        ConfigMap[ConfigMap: consumer-script-config]
        ScaledObject[ScaledObject: keda-rabbitmq]
        TriggerAuth[TriggerAuthentication: rabbitmq-trigger-auth]
    end

    subgraph RabbitMQ Namespace
        RMQ[RabbitMQ Cluster]
    end

    Deployment -->|References| Secret
    Deployment -->|Mounts| ConfigMap
    ScaledObject -->|scaleTargetRef| Deployment
    ScaledObject -->|authenticationRef| TriggerAuth
    TriggerAuth -.->|Resolves env: host from RABBITMQ_URL| Deployment
    KEDA -->|Reads spec| ScaledObject
    KEDA -->|Fetches Auth parameters via| TriggerAuth
    KEDA -->|Retrieves secret value| Secret
    KEDA -->|Connects & Polls| RMQ
```

---

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Running RabbitMQ Cluster (deployed under the `rabbitmq` namespace as per previous labs).
3. Completion of Lab Exercise 8.1.

---

## 📂 Manifests

### 1. RabbitMQ Credentials Secret (`secret.yaml`)
Stores the base64-encoded AMQP connection string.
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: keda-rabbitmq-secret
type: Opaque
data:
  host: YW1xcDovL2RlZmF1bHRfdXNlcl9obUdaRmhkZXdxNjVQNGRJZHg3OnFjOThuNGlHRDdNWVhNQlZGY0lPMm10QjV2b0R1Vl9uQHJhYmJpdG1xLWNsdXN0ZXIucmFiYml0bXEuc3ZjLmNsdXN0ZXIubG9jYWw6NTY3Mg==
```

### 2. Consumer Workload (`consumer.yaml`)
Deploys a consumer script ConfigMap and a single-replica Deployment that references the credentials secret.
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: consumer-script-config
data:
  consumer-script.sh: |
    #!/bin/bash
    currentMessage=""
    handle_sigterm() {
      if [ -n "$currentMessage" ]; then
        echo "SIGTERM signal received while processing a message."
        curl -X POST http://result-analyzer-service:8080/kill/count -s
        echo "Kill count HTTP request sent."
      else
        echo "SIGTERM signal received, but no message was being processed."
      fi
      exit 0
    }
    trap 'handle_sigterm' SIGTERM
    while true; do
      echo -e "Waiting for message...\n"
      if ! currentMessage=$(amqp-consume --url="$RABBITMQ_URL" -q "testqueue" -c 1 cat); then
        echo -e "Error occurred during message consumption. Exiting...\n"
        continue
      fi
      echo -e "Message received, processing: $currentMessage \n"
      i=1
      while [ $i -le 360 ]; do
        echo "Encoding video $i"
        sleep 1
        i=$((i+1))
      done
      currentMessage=""
      curl -X POST http://result-analyzer-service:8080/create/count -s
      echo -e "Waiting for next message...\n"
    done
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: consumer-program
spec:
  replicas: 1
  selector:
    matchLabels:
      app: consumer-program
  template:
    metadata:
      labels:
        app: consumer-program
    spec:
      containers:
      - name: consumer-program
        image: ghcr.io/kedify/blog05-cli-consumer-program:latest
        command: ["/bin/bash"]
        args: ["/scripts/consumer-script.sh"]
        volumeMounts:
        - name: script-volume
          mountPath: /scripts
        env:
        - name: RABBITMQ_URL
          valueFrom:
            secretKeyRef:
              name: keda-rabbitmq-secret
              key: host
      volumes:
      - name: script-volume
        configMap:
          name: consumer-script-config
```

### 3. TriggerAuthentication & ScaledObject (`scaled-object-trigger-auth-env.yaml`)
Configures the `TriggerAuthentication` resource resolving the RabbitMQ host address from environment variables and links it to the `ScaledObject` trigger configuration.
```yaml
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: rabbitmq-trigger-auth
spec:
  env:
  - parameter: host
    name: RABBITMQ_URL
    containerName: consumer-program
---
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: keda-rabbitmq
spec:
  scaleTargetRef:
    name: consumer-program
  triggers:
  - type: rabbitmq
    metadata:
      protocol: amqp
      queueName: testqueue
      queueLength: "5"
    authenticationRef:
      name: rabbitmq-trigger-auth
```

---

## 🛠️ Step-by-Step Lab Walkthrough

### 1. Deploy the Workload
1. Deploy the Secret, ConfigMap, and the Deployment:
   ```bash
   kubectl apply -f secret.yaml
   kubectl apply -f consumer.yaml
   ```

2. Confirm the consumer pod is running:
   ```bash
   kubectl get pods
   ```

### 2. Deploy the Triggered Scaler
1. Apply the combined `TriggerAuthentication` and `ScaledObject` configuration:
   ```bash
   kubectl apply -f scaled-object-trigger-auth-env.yaml
   ```

2. Verify that the ScaledObject successfully authenticated and is in `READY: True` state:
   ```bash
   kubectl get scaledobjects.keda.sh keda-rabbitmq
   ```
   *Expected Output:*
   ```text
   NAME            SCALETARGETKIND      SCALETARGETNAME    MIN   MAX   READY   ACTIVE    FALLBACK   PAUSED   TRIGGERS   AUTHENTICATIONS       AGE
   keda-rabbitmq   apps/v1.Deployment   consumer-program               True    Unknown   False      False    rabbitmq   rabbitmq-trigger-auth 8s
   ```
   > [!NOTE]
   > Notice the `AUTHENTICATIONS` column points to `rabbitmq-trigger-auth`. KEDA uses this reference to resolve credentials dynamically from the workload container environment.

---

## 🧹 Clean Up

To clean up all resources created in this exercise:
```bash
kubectl delete -f scaled-object-trigger-auth-env.yaml
kubectl delete -f consumer.yaml
kubectl delete -f secret.yaml
```