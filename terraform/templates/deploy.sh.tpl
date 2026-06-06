#!/bin/bash
set -e

# 1. Instalar Docker Compose plugin si no existe
if ! command -v docker compose &> /dev/null; then
  dnf install -y docker-compose-plugin
fi

systemctl enable docker
systemctl start docker

# 2. Clonar repositorio
cd /home/ec2-user

REPO_DIR="my_pets_monorepo"
if [ ! -d "$REPO_DIR" ]; then
  git clone https://github.com/afperdomo2/my_pets_monorepo.git
fi

cd "$REPO_DIR"
git pull origin main || true

# 3. Obtener IP pública (desde AWS metadata o servicio externo)
EC2_IP=$(curl -s http://checkip.amazonaws.com 2>/dev/null || curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo "localhost")

# 4. Crear archivo .env
cat > .env <<EOF
DATABASE_URL=host=${rds_endpoint} user=${db_user} password=${db_password} dbname=${db_name} port=5432 sslmode=disable
JWT_SECRET=${jwt_secret}
PORT=8080
GIN_MODE=release
APP_URL=http://$${EC2_IP}
FRONTEND_URL=http://$${EC2_IP}
EOF

# 5. Desplegar con Docker Compose
docker compose -f docker-compose.cloud.yml pull 2>/dev/null || true
docker compose -f docker-compose.cloud.yml up -d --build
