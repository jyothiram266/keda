# Lab Exercise 8.3: Implementing TriggerAuthentication Referencing a Secret

This exercise demonstrates how KEDA can authenticate against event sources (such as RabbitMQ) by referencing a Kubernetes Secret directly within a `TriggerAuthentication` resource.

Unlike Lab 21, where KEDA resolved credentials by reading environment variables from the target workload, referencing a Secret directly via `secretTargetRef` provides a cleaner separation of concerns. In this model:
- The target deployment workloads do **not** need to expose connection secrets as environment variables.
- KEDA retrieves authentication details directly from the Kubernetes Secret using its service account permissions.

---

## 🏗️ Architecture & Secret Resolution Flow

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

    Deployment -->|Mounts| ConfigMap
    ScaledObject -->|scaleTargetRef| Deployment
    ScaledObject -->|authenticationRef| TriggerAuth
    TriggerAuth -.->|References secret| Secret
    KEDA -->|Reads spec| ScaledObject
    KEDA -->|Fetches credentials via| TriggerAuth
    KEDA -->|Reads secret values directly| Secret
    KEDA -->|Connects & Polls| RMQ
```

---

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Running RabbitMQ Cluster (deployed under the `rabbitmq` namespace as per previous labs).
3. Completion of Lab Exercise 8.2.

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
Deploys a consumer script ConfigMap and a single-replica Deployment that mounts the script. Notice that the deployment workload **no longer** requires the `RABBITMQ_URL` secret environment variable for scaling verification.
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

### 3. Secret-based TriggerAuthentication & ScaledObject (`scaled-object-trigger-auth-secret.yaml`)
Configures the `TriggerAuthentication` resource using `secretTargetRef` to directly point to the Kubernetes secret and key, and links it to the `ScaledObject` trigger configuration.
```yaml
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: rabbitmq-trigger-auth
spec:
  secretTargetRef:
  - parameter: host
    name: keda-rabbitmq-secret
    key: host
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
   kubectl apply -f scaled-object-trigger-auth-secret.yaml
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
   > Notice the `AUTHENTICATIONS` column points to `rabbitmq-trigger-auth`. Here, KEDA references the secret directly via `secretTargetRef` without relying on container environment variables.

   ![KEDA Secret TriggerAuthentication Verification](1.png)

---

## 🧹 Clean Up

To clean up all resources created in this exercise:
```bash
kubectl delete -f scaled-object-trigger-auth-secret.yaml
kubectl delete -f consumer.yaml
kubectl delete -f secret.yaml
```