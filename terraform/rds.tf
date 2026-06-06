# Grupo de subredes para el RDS (usa la subred pública)
resource "aws_db_subnet_group" "main" {
  name       = "${local.name_prefix}-rds-subnet"
  subnet_ids = [aws_subnet.public.id]

  tags = merge({ Name = "${local.name_prefix}-rds-subnet" }, local.common_tags)
}

# Instancia PostgreSQL 16 administrada
resource "aws_db_instance" "main" {
  identifier        = "${local.name_prefix}-pg"
  engine            = "postgres"
  engine_version    = "16"
  instance_class    = "db.t3.micro"
  allocated_storage = 20
  storage_type      = "gp2"

  db_name  = var.db_name
  username = var.db_user
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  backup_retention_period = 0
  skip_final_snapshot     = true
  publicly_accessible     = false # No expone el RDS a Internet

  tags = merge({ Name = "${local.name_prefix}-pg" }, local.common_tags)
}
