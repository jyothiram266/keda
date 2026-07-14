# Lab Exercise 8.4 Integrating External Secret Management Solutions


# Lab Exercise 8.4: Integrating External

Secret Management Solutions
In this exercise, we will explore how to use the external secret manager store with TriggerAuthentication CRD
for secure authentication.

## Prerequisites

1. Basic understanding of Kubernetes and KEDA.
2. Access to a Kubernetes environment with KEDA and Metric Server installed as per Lab 5.
3. Completion of Lab Exercises 8.1, 8.2, and 8.3.

## Lab Exercise

1. Create secret in Vault:
Execute the command below to create the required secrets that will be referred to in the TriggerAuthentication
CRD.
```bash
kubectl exec -ti -n vault vault-0 -- vault kv put secret/rabbitmq
```
host=amqp://default_user_hmGZFhdewq65P4dIdx7:qc98n4iGD7MYXMBVFcIO2mtB5voDuV_n@rab
bitmq-cluster.rabbitmq.svc.cluster.local:5672
2. Create ScaledObject:
As discussed in the chapter, this KEDA TriggerAuthentication Custom Resource Definition (CRD) is configured
to securely authenticate with HashiCorp Vault for accessing secrets needed by KEDA scalers. It defines a
vault-trigger-auth resource within the default namespace that:
- Specifies the Vault server's address (http://vault.vault.svc.cluster.local:8200).
- Uses token-based authentication with the root token.
- Retrieves the host parameter from the rabbitmq secret located at secret/data/rabbitmq in Vault,
mapping it to the host key in the KEDA scaler's configuration.
Create a file name scaled-object-hashi-secret.yaml with the following contents and apply it using the
command below.
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
```bash
kubectl apply -f scaled-object-hashi-secret.yaml
```
3. Verify ScaledObject:
```bash
kubectl get scaledobjects.keda.sh keda-rabbitmq
```
NAME SCALETARGETKIND SCALETARGETNAME MIN MAX TRIGGERS
AUTHENTICATION READY ACTIVE FALLBACK PAUSED AGE
keda-rabbitmq apps/v1.Deployment consumer-program rabbitmq
True False Unknown Unknown 7s
4. Clean up:
```bash
kubectl delete -f scaled-object-hashi-secret.yaml
```

## Summary

In this exercise, we learned how to use an external secret manager (HashiCorp Vault) with the
TriggerAuthentication CRD. The exercise involved creating a RabbitMQ secret in Vault, configuring
TriggerAuthentication to authenticate via Vault, and linking it to a ScaledObject.