# docker-machine-dns

DNS for Docker machines, allows to access them with the following domain format : `{machine}.docker`.

## Installation

    go get github.com/bamarni/docker-machine-dns

## Usage

    docker-machine-dns [options...]

    Options:
      -port   Port to listen on (defaults to "53").
      -user   Execute the "docker-machine ip" command as this user (defaults to "SUDO_USER")

## Usage example : Mac OSX

As Docker only runs on Linux, Mac users need a local VM, using VirtualBox for example.

The thing is that when creating this machine, docker-machine will assign to it a dynamic IP address.
It'd be more convenient to access it through a domain name instead. Here comes `docker-machine-dns`!

The program can be run like this : `sudo docker-machine-dns`

*Root privileges are required because of the default low port and the fact that a DNS resolver configuration file
has to be created at `/etc/resolver/docker`.*

To make sure it works properly (let's say for a machine called `dev`) :

    > dig @localhost dev.docker +short
    192.168.99.100
    > ping dev.docker
    192.168.99.100

## How it works

It spins up a simplistic DNS server, only listening for questions about A records.

Behind the scene it will run `docker-machine ip {machine}` in order to resolve the IP address of a given machine.
