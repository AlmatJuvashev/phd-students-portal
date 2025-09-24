# Deployment Manifests

## Kubernetes (Helm)
```
cd deploy/helm/phd-portal
helm upgrade --install phd-portal .   --set image.backend=YOUR_REGISTRY/phd-portal-backend:TAG   --set image.frontend=YOUR_REGISTRY/phd-portal-frontend:TAG
```
Set secrets (e.g. `DATABASE_URL`) via a Secret named `phd-portal-secrets` before install.

## AWS ECS (Fargate)
- Build & push images to ECR
- Customize `deploy/ecs/task-*.json` with your ARNs, VPC subnets, SGs, URLs
- Create services from these task defs (via console or CLI)
