# K8s

## Postgres
```
cloudcart namespace
       |
       ├── Deployment: postgres
       |       |
       |       └── Pod: postgres-65c47d6d64-xlcd6
       |
       ├── Service: postgres
       |       |
       |       └── ClusterIP: 10.102.7.111
       |
       └── PersistentVolumeClaim: postgres-pvc
               |
               └── PersistentVolume: 2Gi
```

One important concept:

Your application Pods will connect to PostgreSQL using:
```postgres:5432```
They don't need to know its ClusterIP address. Kubernetes DNS resolves the Service name automatically.

A few important concepts to understand from postgres.yaml:

- PersistentVolumeClaim: requests 2 GB of persistent storage.
- Secret: keeps database credentials outside the Deployment definition.
- Deployment: manages the PostgreSQL Pod.
- ClusterIP: makes PostgreSQL accessible to other Pods inside the cluster.
- readinessProbe: checks whether PostgreSQL is ready to accept connections.
- livenessProbe: checks whether PostgreSQL is alive.

Test PostgreSQL inside Kubernetes:
Once the Pod is running:
```
kubectl exec -it deployment/postgres -n cloudcart -- \
  psql -U postgres -d cloudcart
```
Inside PostgreSQL:
```
SELECT version();
```

Then:
```
\l
```
You should see the cloudcart database.

Exit:
```
\q
```

## 1. ConfigMaps & Secrets
We'll keep environment-specific values separate from application deployment manifests.

Two files created:
1. configmap.yaml – non-sensitive configuration.
2. secrets.yaml – passwords and JWT secret.

## 2. Auth Service, Product Service & Order Service
Create k8s/auth-service.yaml

We'll configure:
- One replica initially.
- Pull local Docker image(change this later).
- ConfigMap and Secret injection.
- ClusterIP Service.
- Readiness and liveness probes.
- CPU and memory requests/limits.

## 3. Notifiction Service
Unlike the other three services, Notification Service is a background Kafka consumer. It doesn't expose an HTTP endpoint, so we don't need a Kubernetes Service, containerPort, or HTTP probes.

Notice:
- Kafka broker and topic come from the common ConfigMap.
- Consumer group ID is specific to Notification Service.
- No JWT Secret is needed because this service currently only consumes Kafka events.
- No HTTP Service is required.

