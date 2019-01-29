provider "aws" {
  region = "us-east-2"
}

resource "aws_security_group" "goad_traffic" {
  name        = "goad_traffic"
  description = "Allow all inbound traffic"

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8081
    to_port     = 8081
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "goad" {
  ami           = "ami-06cb4dd657e406acb"
  instance_type = "t2.nano"
  key_name = "jenkins"
  security_groups = ["goad_traffic", "default"]
}

output "goad_dns" {
  value = ["${aws_instance.goad.public_dns}"]
}