#!/bin/bash

echo "================================"
echo "  Cloud Run Proxy — Deploy Tool"
echo "================================"
echo ""

# 1. Project ID
read -p "Введи Project ID нового аккаунта: " PROJECT_ID
gcloud config set project "$PROJECT_ID"
echo ""

# 2. Включить необходимые API
echo "Включаю необходимые API..."
gcloud services enable run.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com --project="$PROJECT_ID"
echo ""

# 3. Имя сервиса
read -p "Введи имя сервиса (например ru, fi, ca): " SERVICE_NAME
echo ""

# 4. Регион
echo "Выбери регион:"
echo "  1) europe-west1     (Belgium)   — для ВПС в СПб/Европе"
echo "  2) europe-north1    (Finland)   — для ВПС в Финляндии/Швеции"
echo "  3) northamerica-northeast1 (Montreal) — для ВПС в Канаде"
echo "  4) Ввести вручную"
read -p "Твой выбор (1-4): " REGION_CHOICE

case $REGION_CHOICE in
  1) REGION="europe-west1" ;;
  2) REGION="europe-north1" ;;
  3) REGION="northamerica-northeast1" ;;
  4) read -p "Введи регион вручную: " REGION ;;
  *) REGION="europe-west1" ;;
esac
echo ""

# 5. Docker образ
read -p "Введи Docker образ (например docker.io/riestru/google-cloud-run-proxy): " IMAGE
echo ""

# 6. IP сервера V2Ray
read -p "Введи IP твоего V2Ray сервера: " V2RAY_IP
echo ""

# 7. Деплой
echo "Разворачиваю сервис..."
echo ""

gcloud run deploy "$SERVICE_NAME" \
  --image="$IMAGE" \
  --platform=managed \
  --region="$REGION" \
  --allow-unauthenticated \
  --set-env-vars="V2RAY_SERVER_IP=$V2RAY_IP" \
  --memory=128Mi \
  --cpu=1 \
  --min-instances=0 \
  --max-instances=1 \
  --concurrency=1000 \
  --timeout=3600 \
  --project="$PROJECT_ID"

echo ""
echo "================================"
echo "  Готово!"
echo "================================"
