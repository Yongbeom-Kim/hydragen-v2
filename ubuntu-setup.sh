#!/bin/bash
# Copy paste this into the local server (ubuntu).

##### Firewall #####
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

##### User #####
sudo useradd -m -s /bin/bash app
sudo passwd app
sudo usermod -aG sudo app
su app
cd /home/app

# Copy authorized_keys
sudo mkdir -p /home/app/.ssh
sudo cp ~/.ssh/authorized_keys /home/app/.ssh/authorized_keys
sudo chown -R app:app /home/app/.ssh
sudo chmod 700 /home/app/.ssh
sudo chmod 600 /home/app/.ssh/authorized_keys

##### Docker #####
sudo apt remove $(dpkg --get-selections docker.io docker-compose docker-compose-v2 docker-doc podman-docker containerd runc | cut -f1)

# Add Docker's official GPG key:
sudo apt update
sudo apt install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Signed-By: /etc/apt/keyrings/docker.asc
EOF

sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo systemctl start docker
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker

sudo systemctl enable docker.service
sudo systemctl enable containerd.service

# make

sudo apt install make

## Finally, add to GitHub Actions (offline-notion label)