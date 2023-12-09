# Chat

# Prerequisites
Before you begin, make sure you have an internet connection and administrative privileges on your machine. **Please run entire step on this document before you go to workshop**.
\
&nbsp;

# Installing Go (Golang)
## macOS
- Download the latest macOS package from the official Go website: https://golang.org/dl/
- Open the downloaded package and follow the installation instructions.
- Verify your Go installation by opening a terminal and running:
```
go version
```

## Windows
- Download the latest Windows installer from the official Go website: https://golang.org/dl/
- Run the installer and follow the on-screen instructions.
- Open the Windows Command Prompt or PowerShell and verify your Go installation by running:
```
go version
```

## Linux
- Open a terminal and use the package manager to install Go. For example, on Ubuntu, you can run:
```
sudo apt-get update
sudo apt-get install golang
```
- Verify your Go installation by running:
```
go version
```
&nbsp;

# Installing Docker
# macOS
- Download the Docker Desktop for Mac from the official Docker website: https://www.docker.com/products/docker-desktop
- Install Docker Desktop by following the on-screen instructions.
- After installation, start Docker Desktop.

# Windows
- Download Docker Desktop for Windows from the official Docker website: https://www.docker.com/products/docker-desktop
- Install Docker Desktop by following the on-screen instructions.
- After installation, start Docker Desktop.

# Linux
- Install Docker on Linux using the official Docker documentation for your specific distribution. For example, on Ubuntu, you can follow these steps: https://docs.docker.com/engine/install/ubuntu/
- After installation, start the Docker service:
```
sudo systemctl start docker
```
- Add your user to the docker group to run Docker commands without sudo:
```
sudo usermod -aG docker $USER
```
- Log out and log back in or reboot your system to apply the group changes.
- Verify your Docker installation by running:
```
docker --version
```
&nbsp;

# Run project
Run this command:
```
docker-compose up
```
&nbsp;

# Healthcheck
- [nsq-admin](http://localhost:3004/)
- [service](http://localhost:8080/check)