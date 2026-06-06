# Busca la AMI más reciente de Ubuntu LTS
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-*-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# Servidor EC2 con Docker Compose para la API y frontend
resource "aws_instance" "main" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = "t3.micro"
  subnet_id              = aws_subnet.public.id
  vpc_security_group_ids = [aws_security_group.ec2.id]
  key_name               = var.ssh_key_name

  user_data = templatefile("${path.module}/templates/deploy.sh.tpl", {
    rds_endpoint = aws_db_instance.main.endpoint
    db_user      = var.db_user
    db_password  = var.db_password
    db_name      = var.db_name
    jwt_secret   = var.jwt_secret
  })

  tags = merge({ Name = "${local.name_prefix}-api-01" }, local.common_tags)
}

# IP pública fija para el EC2
resource "aws_eip" "main" {
  domain = "vpc"

  tags = merge({ Name = "${local.name_prefix}-eip" }, local.common_tags)
}

# Asocia la IP pública al EC2
resource "aws_eip_association" "main" {
  allocation_id = aws_eip.main.id
  instance_id   = aws_instance.main.id
}
