# Lab Exercise 8.4: Integrating External Secret Management Solutions

This exercise explores how KEDA can securely fetch connection secrets directly from an external secret manager, specifically **HashiCorp Vault**, utilizing KEDA's native Vault integration within a `TriggerAuthentication` resource.

This approach provides a robust security architecture by removing the need to manage Kubernetes Secret resources for autoscaling. KEDA queries the external Vault API directly at runtime to retrieve sensitive parameters (e.g., connection strings).

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
        TriggerAuth[TriggerAuthentication: vault-trigger-auth]
    end

    subgraph vault Namespace
        Vault[HashiCorp Vault: vault-0]
    end

    subgraph RabbitMQ Namespace
        RMQ[RabbitMQ Cluster]
    end

    Deployment -->|References| Secret
    Deployment -->|Mounts| ConfigMap
    ScaledObject -->|scaleTargetRef| Deployment
    ScaledObject -->|authenticationRef| TriggerAuth
    TriggerAuth -.->|Fetches via token| Vault
    KEDA -->|Reads spec| ScaledObject
    KEDA -->|Queries secrets via| TriggerAuth
    KEDA -->|Reads 'host' key from secret/rabbitmq| Vault
    KEDA -->|Connects & Polls| RMQ
```

---

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Running RabbitMQ Cluster (deployed under the `rabbitmq` namespace as per previous labs).
3. HashiCorp Vault installed in the `vault` namespace. If it's not running, install it using:
   ```bash
   helm upgrade -i vault hashicorp/vault --set "server.dev.enabled=true" -n vault --create-namespace
   ```
4. Completion of Lab Exercise 8.3.

---

## 📂 Manifests

### 1. RabbitMQ Credentials Secret (`secret.yaml`)
Stores the base64-encoded AMQP connection string used by the consumer pod itself.
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
Deploys the consumer ConfigMap and Deployment.
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

### 3. Vault-based TriggerAuthentication & ScaledObject (`scaled-object-hashi-secret.yaml`)
Configures `TriggerAuthentication` to fetch credentials from Vault using the token method, and associates it with the `ScaledObject` trigger.
```yaml
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: vault-trigger-auth
  namespace: default
spec:
  hashiCorpVault:
    address: "http://vault.vault.svc.cluster.local:8200"
    authentication: token
    credential:
      token: "root"
    secrets:
    - parameter: "host"
      key: "host"
      path: "secret/data/rabbitmq"
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
      name: vault-trigger-auth
```

---

## 🛠️ Step-by-Step Lab Walkthrough

### 1. Write the RabbitMQ connection string to HashiCorp Vault
1. Write the secret key-value to the dev Vault KV engine:
   ```bash
   kubectl exec -ti -n vault vault-0 -- vault kv put secret/rabbitmq host=amqp://default_user_hmGZFhdewq65P4dIdx7:qc98n4iGD7MYXMBVFcIO2mtB5voDuV_n@rabbitmq-cluster.rabbitmq.svc.cluster.local:5672
   ```

2. Verify that the secret was successfully stored in Vault:
   ```bash
   kubectl exec -ti -n vault vault-0 -- vault kv get secret/rabbitmq
   ```

### 2. Deploy the Workload
1. Deploy the Secret, ConfigMap, and the Deployment:
   ```bash
   kubectl apply -f secret.yaml
   kubectl apply -f consumer.yaml
   ```

2. Confirm the consumer pod is running:
   ```bash
   kubectl get pods
   ```

### 3. Deploy the Triggered Scaler
1. Apply the combined `TriggerAuthentication` and `ScaledObject` configuration:
   ```bash
   kubectl apply -f scaled-object-hashi-secret.yaml
   ```

2. Verify that the ScaledObject successfully authenticated and is in `READY: True` state:
   ```bash
   kubectl get scaledobjects.keda.sh keda-rabbitmq
   ```
   *Expected Output:*
   ```text
   NAME            SCALETARGETKIND      SCALETARGETNAME    MIN   MAX   READY   ACTIVE    FALLBACK   PAUSED   TRIGGERS   AUTHENTICATIONS      AGE
   keda-rabbitmq   apps/v1.Deployment   consumer-program               True    Unknown   False      False    rabbitmq   vault-trigger-auth   8s
   ```
   > [!NOTE]
   > KEDA connects to Vault, uses the root token to authenticate, retrieves the secret from `secret/data/rabbitmq`, maps the `host` key to the RabbitMQ scaler parameter, and verifies the readiness of the ScaledObject.

---

## 🧹 Clean Up

To clean up all resources created in this exercise:
```bash
kubectl delete -f scaled-object-hashi-secret.yaml
kubectl delete -f consumer.yaml
kubectl delete -f secret.yaml
kubectl exec -ti -n vault vault-0 -- vault kv delete secret/rabbitmq
```