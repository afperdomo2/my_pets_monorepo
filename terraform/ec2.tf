data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

resource "aws_instance" "main" {
  ami                    = data.aws_ami.amazon_linux_2023.id
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

  tags = { Name = "my-pets-ec2" }
}

resource "aws_eip" "main" {
  domain = "vpc"

  tags = { Name = "my-pets-eip" }
}

resource "aws_eip_association" "main" {
  allocation_id = aws_eip.main.id
  instance_id   = aws_instance.main.id
}
