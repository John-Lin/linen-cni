# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.hostname = 'dev'

  config.vm.provision "shell", privileged: false, inline: <<-SHELL
    set -e -x -u
    sudo apt-get update
    sudo apt-get install -y vim git build-essential openvswitch-switch bridge-utils

    # Install Golang
    wget --quiet https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz
    sudo tar -zxf go1.9.1.linux-amd64.tar.gz -C /usr/local/

    echo 'export GOROOT=/usr/local/go' >> /home/ubuntu/.bashrc
    echo 'export GOPATH=$HOME/go' >> /home/ubuntu/.bashrc
    echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> /home/ubuntu/.bashrc
    source /home/ubuntu/.bashrc

    mkdir -p /home/ubuntu/go/src

    rm -rf /home/ubuntu/go1.9.1.linux-amd64.tar.gz

    # Download CNI and CNI plugins binaries
    wget --quiet https://github.com/containernetworking/cni/releases/download/v0.6.0/cni-amd64-v0.6.0.tgz
    wget --quiet https://github.com/containernetworking/plugins/releases/download/v0.6.0/cni-plugins-amd64-v0.6.0.tgz
    mkdir cni
    tar -zxf cni-amd64-v0.6.0.tgz -C /home/ubuntu/cni
    tar -zxf cni-plugins-amd64-v0.6.0.tgz -C /home/ubuntu/cni

    rm -rf /home/ubuntu/cni-plugins-amd64-v0.6.0.tgz /home/ubuntu/cni-amd64-v0.6.0.tgz

    # Download linen CNI source
    git clone https://github.com/John-Lin/linen-cni.git

  SHELL

  config.vm.provider :virtualbox do |v|
    v.customize ["modifyvm", :id, "--cpus", 2]
    # enable this when you want to have more memory 
    # v.customize ["modifyvm", :id, "--memory", 4096]
    v.customize ['modifyvm', :id, '--nicpromisc1', 'allow-all']
  end
end
