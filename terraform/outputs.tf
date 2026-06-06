output "ec2_public_ip" {
  description = "IP pública del EC2"
  value       = aws_eip.main.public_ip
}

output "rds_endpoint" {
  description = "Endpoint del RDS (host:port)"
  value       = aws_db_instance.main.endpoint
}

output "url" {
  description = "URL de la aplicación"
  value       = "http://${aws_eip.main.public_ip}"
}

output "ssh_command" {
  description = "Comando SSH para conectarse al EC2"
  value       = "ssh -i ${var.ssh_key_name}.pem ubuntu@${aws_eip.main.public_ip}"
}
