# docker-machine-dns

DNS for Docker machines, allows to access them with the following domain format : `{machine}.docker`.

## How it works

It spins up a simplistic DNS server, only listening for questions about A records.

Behind the scene it will run `docker-machine ip {machine}` in order to resolve the IP address of a given machine.

## Installation

    go get github.com/bamarni/docker-machine-dns

## Usage

    docker-machine-dns [options...]

    Options:
      -port   Port to listen on (defaults to "10053").

## Usage example : Mac OSX

As Docker only runs on Linux, Mac users need a local VM, using VirtualBox for example.

The thing is that when creating this machine, docker-machine will assign to it a dynamic IP address.
It'd be more convenient to access it through a domain name instead. And here comes `docker-machine-dns`!

### Create a DNS resolver configuration

The first thing to do is to create a DNS resolver for the `.docker` TLD.

This is done by creating the following configuration file : `sudo nano /etc/resolver/docker`

With the following content :

    nameserver 127.0.0.1
    port 10053

### Run the program

The program can now be run : `docker-machine-dns -port 10053`

To make sure it works properly (let's say for a machine called `dev`) :

    > dig -p 10053 @localhost dev.docker +short
    192.168.99.100

Enjoy!
