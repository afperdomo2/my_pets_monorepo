#!/bin/sh
set -e

# ── Requisitos ─────────────────────────────────────────
: "${AWS_ACCOUNT_ID:?Definir AWS_ACCOUNT_ID}"
: "${AWS_REGION:=us-east-1}"
: "${IMAGE_TAG:=latest}"

ECR="$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"

echo "=== Login ECR ==="
aws ecr get-login-password --region "$AWS_REGION" \
  | docker login --username AWS --password-stdin "$ECR"

echo "=== Construyendo api ==="
docker build \
  -f apps/api/Dockerfile.prod \
  -t "$ECR/my-pets-api:$IMAGE_TAG" \
  -t "$ECR/my-pets-api:latest" \
  apps/api

echo "=== Pusheando api ==="
docker push "$ECR/my-pets-api:$IMAGE_TAG"
docker push "$ECR/my-pets-api:latest"

echo "=== Construyendo web ==="
docker build \
  -f apps/web/Dockerfile.prod \
  -t "$ECR/my-pets-web:$IMAGE_TAG" \
  -t "$ECR/my-pets-web:latest" \
  .

echo "=== Pusheando web ==="
docker push "$ECR/my-pets-web:$IMAGE_TAG"
docker push "$ECR/my-pets-web:latest"

echo "✅ Imágenes publicadas:"
echo "  $ECR/my-pets-api:$IMAGE_TAG"
echo "  $ECR/my-pets-web:$IMAGE_TAG"
