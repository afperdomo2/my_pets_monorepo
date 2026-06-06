variable "aws_region" {
  description = "Región de AWS"
  default     = "us-east-1"
}

variable "org_name" {
  description = "Nombre de la organización"
  default     = "felipecorp"
}

variable "project_name" {
  description = "Nombre del proyecto"
  default     = "my-pets-monorepo"
}

variable "environment" {
  description = "Entorno (dev, staging, prod)"
  default     = "prod"
}

variable "ssh_key_name" {
  description = "Nombre del Key Pair en EC2 para acceso SSH"
  default     = "pruebas-felipe-ssh"
}

variable "db_name" {
  description = "Nombre de la base de datos"
  default     = "my_pets"
}

variable "db_user" {
  description = "Usuario de la base de datos"
  default     = "my_pets_user"
}

variable "db_password" {
  description = "Contraseña de la base de datos"
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "Secreto para firmar JWT"
  type        = string
  sensitive   = true
}
