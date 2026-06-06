# 🏗️ Terraform — Deploy en AWS (free tier)

> Guía completa para desplegar este proyecto en AWS usando Terraform,
> aprovechando la capa gratuita.

---

## 📋 Tabla de contenido

- [Arquitectura](#arquitectura)
- [Recursos AWS](#recursos-aws)
- [Prerrequisitos](#prerrequisitos)
- [Configuración](#configuración)
- [Deploy](#deploy)
- [Outputs](#outputs)
- [Probar la aplicación](#probar-la-aplicación)
- [Solución de problemas](#solución-de-problemas)
- [Destruir todo](#destruir-todo)
- [Estructura de Terraform](#estructura-de-terraform)
- [Costo estimado](#costo-estimado)

---

## 🏛️ Arquitectura

```
                    Internet
                        │
                        ▼
              ┌─────────────────┐
              │   EC2 t3.micro  │
              │  (Amazon Linux  │
              │     2023)       │
              │                 │
              │  ┌───────────┐  │
              │  │  nginx    │  │  Puerto 80
              │  │           │  │
              │  │  /api/*   ───┼──→ Go API (:8080)
              │  │  /*       │  │  (Vue SPA)
              │  └───────────┘  │
              │                 │
              │  ┌───────────┐  │
              │  │  Go API   │  │
              │  │  (Gin)    │  │
              │  └─────┬─────┘  │
              └────────┼────────┘
                       │
                       ▼
              ┌─────────────────┐
              │  RDS PostgreSQL │
              │   db.t3.micro   │
              │     20GB        │
              └─────────────────┘
```

### Flujo de una petición

```
1. Usuario abre http://54.xx.xx.xx en el navegador
                        │
2. nginx recibe la petición en puerto 80
                        │
           ┌────────────┴────────────┐
           ▼                         ▼
    ¿Ruta /api/*?              Ruta /* (SPA)
           │                         │
           ▼                         ▼
    nginx redirige a          nginx sirve archivos
    Go API en :8080           estáticos (Vue)
           │
           ▼
    Go consulta RDS
    PostgreSQL
```

---

## 📦 Recursos AWS

| Recurso | Tipo | Descripción |
|---|---|---|
| **VPC** | `10.0.0.0/16` | Red virtual aislada |
| **Subnet** | `10.0.1.0/24` | Subred pública en `us-east-1a` |
| **Internet Gateway** | — | Permite tráfico saliente/entrante a Internet |
| **Route Table** | — | Enruta `0.0.0.0/0` → Internet Gateway |
| **Security Group EC2** | — | Permite HTTP (80) y SSH (22) desde cualquier IP |
| **Security Group RDS** | — | Permite PostgreSQL (5432) solo desde el EC2 |
| **EC2** | `t3.micro` | Servidor principal: nginx + Go API + Docker |
| **Elastic IP** | — | IP pública fija asociada al EC2 |
| **RDS** | `db.t3.micro` | PostgreSQL 16 administrado, 20GB gp2 |

### ¿Por qué estos recursos?

| Decisión | Motivo |
|---|---|
| **EC2 + Docker Compose** | Más simple que ECS/Fargate para un proyecto de prueba |
| **RDS separado** | No queremos PostgreSQL dentro del EC2; lo administra AWS |
| **Elastic IP** | IP fija para no tener que cambiar la URL cada vez que el EC2 reinicia |
| **Un solo EC2** | nginx + Go en el mismo servidor, sin balanceador |

---

## 📋 Prerrequisitos

### 1. Cuenta de AWS

Si no tienes, crea una en [aws.amazon.com](https://aws.amazon.com/free/).

### 2. Terraform instalado

```bash
# Linux (Ubuntu/Debian)
wget -O - https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform

# macOS
brew tap hashicorp/tap
brew install hashicorp/tap/terraform

# Windows (con chocolatey)
choco install terraform
```

Verificar: `terraform --version`

### 3. AWS Access Keys

Crea un par de claves en [IAM → Users → Security credentials](https://console.aws.amazon.com/iam/):

```bash
export AWS_ACCESS_KEY_ID="AKIAXXXXXXXXXXX"
export AWS_SECRET_ACCESS_KEY="xxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### 4. Key Pair EC2

Creado desde la consola AWS → EC2 → Key Pairs → Create key pair.

- **Nombre:** `pruebas-felipe-ssh`
- **Tipo:** RSA
- **Formato:** `.pem` (Linux/macOS) o `.ppk` (Windows)

Guardar el archivo en `~/.ssh/pruebas-felipe-ssh.pem` y darle permisos:

```bash
chmod 400 ~/.ssh/pruebas-felipe-ssh.pem
```

---

## ⚙️ Configuración

### Variables de Terraform

Hay 6 variables. Dos de ellas son **sensibles** y no tienen valor por defecto (Terraform las pedirá al ejecutar):

| Variable | Default | Descripción |
|---|---|---|
| `aws_region` | `us-east-1` | Región de AWS |
| `ssh_key_name` | `pruebas-felipe-ssh` | Nombre del Key Pair EC2 |
| `db_name` | `my_pets` | Nombre de la base de datos |
| `db_user` | `my_pets_user` | Usuario de la base de datos |
| `db_password` | *(obligatorio)* | Contraseña de la base de datos |
| `jwt_secret` | *(obligatorio)* | Secreto para firmar tokens JWT |

### Exportar variables sensibles

Para que Terraform no las pregunte en cada `apply`:

```bash
export TF_VAR_db_password="MiClaveSuperSegura2024"
export TF_VAR_jwt_secret="MiJwtSecretoNadieLoSabe"
```

> **Alternativa (persistente):** crear un archivo `terraform.tfvars` a partir de [`terraform.tfvars.example`](terraform/terraform.tfvars.example):
>
> ```bash
> cp terraform.tfvars.example terraform.tfvars
> # Editar terraform.tfvars con los valores reales
> ```
>
> Terraform lo carga automáticamente. Así evitas tener que exportar las variables cada vez que abres una terminal nueva.

---

## 🚀 Deploy

Paso a paso:

```bash
# 1. Ir al directorio de Terraform
cd terraform

# 2. Exportar variables sensibles
export TF_VAR_db_password="MiClaveSuperSegura2024"
export TF_VAR_jwt_secret="MiJwtSecretoNadieLoSabe"

# 3. Inicializar Terraform (descarga plugins)
terraform init

# 4. Ver el plan (opcional, muestra qué va a crear)
terraform plan

# 5. Aplicar (crear todo en AWS)
terraform apply -auto-approve
```

### ¿Qué pasa durante el apply?

```
1. Crea la VPC y subred          (~10s)
2. Crea Internet Gateway         (~5s)
3. Carga Security Groups         (~5s)
4. Crea RDS PostgreSQL           (~3 min) ← lo más lento
5. Crea Elastic IP               (~5s)
6. Crea EC2 t3.micro             (~30s)
7. Asocia Elastic IP al EC2      (~5s)
8. El EC2 arranca y ejecuta:
   a. Instala Docker Compose
   b. Clona el repositorio
   c. Crea archivo .env
   d. docker compose up -d       (~2 min)

Tiempo total: ~5-7 minutos
```

---

## 📤 Outputs

Después de `terraform apply`, verás:

```
ec2_public_ip = "54.198.xxx.xxx"
rds_endpoint  = "my-pets-db.c9xxxxx.us-east-1.rds.amazonaws.com:5432"
url           = "http://54.198.xxx.xxx"
ssh_command   = "ssh -i pruebas-felipe-ssh.pem ec2-user@54.198.xxx.xxx"
```

Para ver los outputs en cualquier momento:

```bash
cd terraform
terraform output
```

O un output específico:

```bash
terraform output -raw ec2_public_ip
terraform output -raw url
```

---

## ✅ Probar la aplicación

### Desde terminal

```bash
# Obtener la URL
URL=$(terraform output -raw url)

# Health check de la API
curl $URL/api/v1/health

# Respuesta esperada:
# {"service":"my-pets-api","status":"ok"}
```

### Desde el navegador

Abrir `http://$(terraform output -raw url)`

Deberías ver la aplicación Vue cargando (puede tardar unos segundos si el frontend aún está compilándose en el EC2).

### Ver logs del servidor

```bash
ssh -i ~/.ssh/pruebas-felipe-ssh.pem ec2-user@$(terraform output -raw ec2_public_ip)

# Ver logs de Docker
docker compose -f docker-compose.cloud.yml logs

# Ver logs de la API
docker compose -f docker-compose.cloud.yml logs api

# Ver logs de nginx
docker compose -f docker-compose.cloud.yml logs web
```

---

## 🔍 Solución de problemas

### La API responde 502 o 503

```bash
# El Go API puede estar arrancando todavía
ssh -i ~/.ssh/pruebas-felipe-ssh.pem ec2-user@$(terraform output -raw ec2_public_ip)

# Revisar si los contenedores están corriendo
docker compose -f docker-compose.cloud.yml ps

# Ver logs
docker compose -f docker-compose.cloud.yml logs api
```

### La página web carga en blanco

```bash
# Revisar nginx
docker compose -f docker-compose.cloud.yml logs web

# Revisar que la build del frontend esté completa
docker compose -f docker-compose.cloud.yml exec web ls /usr/share/nginx/html/
```

### Error de conexión a la base de datos

```bash
# Verificar que el RDS esté accesible desde el EC2
ssh -i ~/.ssh/pruebas-felipe-ssh.pem ec2-user@$(terraform output -raw ec2_public_ip)

# Probar conexión a PostgreSQL
docker run --rm postgres:16-alpine psql -h $(terraform output -raw rds_endpoint | cut -d: -f1) -U my_pets_user -d my_pets
```

### Quiero re-deployar desde cero

```bash
# Destruir todo
terraform destroy -auto-approve

# Esperar a que termine (~3 min)

# Volver a crear
terraform apply -auto-approve
```

---

## 🗑️ Destruir todo

Para eliminar todos los recursos y no generar más costos:

```bash
cd terraform
terraform destroy -auto-approve
```

### ¿Qué se elimina?

| Recurso | ¿Se elimina? |
|---|---|
| EC2 | ✅ Sí |
| RDS + volumen 20GB | ✅ Sí |
| Elastic IP | ✅ Sí |
| VPC + subred | ✅ Sí |
| Security Groups | ✅ Sí |
| Key Pair | ❌ **No** (es independiente, creado aparte) |

### Tiempo estimado de destrucción

```
RDS deletion:       ~2 min  (sin snapshot final)
EC2 termination:    ~30s
VPC teardown:       ~10s
Total:              ~3 min
```

---

## 📁 Estructura de Terraform

```
terraform/
├── main.tf              # Provider AWS
├── variables.tf         # 6 variables (2 sensibles)
├── outputs.tf           # IP, RDS endpoint, URL, SSH
├── network.tf           # VPC, subnet, IGW, route table
├── security.tf          # Security groups EC2 + RDS
├── rds.tf               # RDS PostgreSQL 16
├── ec2.tf               # EC2 + Elastic IP + user_data
└── templates/
    └── deploy.sh.tpl    # Script bootstrap del EC2
```

### Explicación de cada archivo

#### `main.tf`

Configura el provider de AWS y la región.

```hcl
provider "aws" {
  region = var.aws_region
}
```

#### `variables.tf`

Define las 6 variables del proyecto. Dos son `sensitive = true` para que nunca se muestren en pantalla.

| Variable | Dónde se usa |
|---|---|
| `aws_region` | `main.tf` |
| `ssh_key_name` | `ec2.tf` (para SSH) |
| `db_name`, `db_user`, `db_password` | `rds.tf` + `deploy.sh.tpl` |
| `jwt_secret` | `deploy.sh.tpl` (se inyecta en el `.env`) |

#### `outputs.tf`

Muestra información útil al terminar el `apply`.

#### `network.tf`

- **VPC**: red 10.0.0.0/16 con DNS habilitado
- **Subnet**: subred pública 10.0.1.0/24
- **Internet Gateway**: conecta la VPC a Internet
- **Route Table**: envía todo el tráfico 0.0.0.0/0 al IGW

#### `security.tf`

- **EC2 SG**: permite HTTP (80) y SSH (22) desde cualquier IP
- **RDS SG**: permite PostgreSQL (5432) **solo desde el EC2**

El RDS no es accesible desde Internet, solo desde el servidor.

#### `rds.tf`

- PostgreSQL 16 en `db.t3.micro` (free tier elegible)
- 20GB de almacenamiento gp2
- Sin backups automáticos (`retention_period = 0`)
- Sin snapshot final al destruir (`skip_final_snapshot = true`)
- No es accesible públicamente

#### `ec2.tf`

- Amazon Linux 2023 (viene con Docker instalado)
- `t3.micro` (free tier)
- Asocia Elastic IP fija
- Ejecuta `deploy.sh.tpl` al arrancar (bootstrap)

#### `templates/deploy.sh.tpl`

Script que se ejecuta en el EC2 cuando arranca por primera vez:

```bash
#!/bin/bash
set -e

# 1. Instalar Docker Compose plugin
dnf install -y docker-compose-plugin

# 2. Iniciar Docker
systemctl enable docker
systemctl start docker

# 3. Clonar el repositorio
cd /home/ec2-user
git clone https://github.com/afperdomo2/my_pets_monorepo.git
cd my_pets_monorepo

# 4. Crear .env con los datos de RDS + JWT
cat > .env <<EOF
DATABASE_URL=host=${rds_endpoint} user=${db_user} ...
JWT_SECRET=${jwt_secret}
...
EOF

# 5. Desplegar con Docker Compose
docker compose -f docker-compose.cloud.yml up -d --build
```

Las variables `${rds_endpoint}`, `${db_user}`, etc. son reemplazadas por Terraform antes de subir el script al EC2.

---

## 💰 Costo estimado

### Con free tier activo (cuenta < 12 meses)

| Recurso | Costo |
|---|---|
| EC2 t3.micro (750h/mes) | **$0.00** |
| RDS db.t3.micro (750h/mes, 20GB) | **$0.00** |
| Elastic IP (asociada) | **$0.00** |
| Transferencia de datos | ~$0.00 (uso mínimo) |
| **1 día de prueba** | **~$0.00** |

### Sin free tier (cuenta > 12 meses)

| Recurso | Costo por 24h |
|---|---|
| EC2 t3.micro | ~$0.17 |
| RDS db.t3.micro + 20GB | ~$0.50 |
| Elastic IP | $0.00 |
| **Total por 24h** | **~$0.67** |

### Recomendaciones para no generar costos

1. Usar `terraform destroy` inmediatamente después de probar
2. No dejar el RDS corriendo más de 24h
3. El EC2 solo cuesta mientras está encendido (stop no incurre en cargo de cómputo)
4. El almacenamiento RDS (20GB) cuesta aunque el servidor esté detenido

### Script para calcular costo en tiempo real

```bash
# Costo estimado por hora del RDS
echo "RDS: ~$0.021/hora"
echo "EC2: ~$0.007/hora"
echo "Total: ~$0.028/hora x 24h = ~$0.67/día"
```

---

## 🆚 Comparativa con otras opciones

| Opción | Costo/mes | Esfuerzo | Ideal para |
|---|---|---|---|
| **Terraform + EC2 + RDS** (esta guía) | ~$0 (free tier) | Medio | Prácticas, proyectos personales |
| Docker Compose local | $0 | Bajo | Desarrollo |
| ECS Fargate | ~$25+ | Alto | Producción real |
| Kubernetes (EKS) | ~$75+ | Muy alto | Escalado masivo |

---

## 📚 Referencias

| Recurso | Enlace |
|---|---|
| Documentación de Terraform AWS | https://registry.terraform.io/providers/hashicorp/aws/latest |
| Free Tier AWS | https://aws.amazon.com/free/ |
| Docker Compose | https://docs.docker.com/compose/ |
| Código del proyecto | [`docker-compose.cloud.yml`](../docker-compose.cloud.yml) |
