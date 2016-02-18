# docker-machine-dns

DNS for Docker machines, allows to access them with the following domain format : `{machine}.docker`.

## How it works

It spins up a simplistic DNS server, only listening for questions about A records.

Behind the scene it will run `docker-machine ip {machine}` in order to resolve the IP address of a given machine.

## Installation

### Prebuilt binary (Mac OSX only)

For Mac OSX, a prebuilt binary is available in the [releases](https://github.com/bamarni/docker-machine-dns/releases).

### From source (requires Go)

    go get github.com/bamarni/docker-machine-dns

## Usage

    sudo docker-machine-dns [options...]

    Options:
      -port         Port to listen on (defaults to "53")
      -server-only  Server only, doesn't try to create a resolver configuration file
      -user         Execute the "docker-machine ip" command as this user (defaults to "SUDO_USER")

If you don't feel like running the program with `sudo`, see the [-server-only flag](#server-only)

## Usage example : Mac OSX

As Docker only runs on Linux, Mac users need a local VM, using VirtualBox for example.

The thing is that when creating this machine, docker-machine will assign to it a dynamic IP address.
It'd be more convenient to access it through a domain name instead. Here comes `docker-machine-dns`!

Run the program :

    > sudo docker-machine-dns
    2016/02/18 10:39:52 Creating configuration file at /etc/resolver/docker...
    2016/02/18 10:39:52 Listening on :53...

In another terminal, to make sure it works properly (let's say for a machine called `dev`) :

    > dig @localhost dev.docker +short
    192.168.99.100
    > ping dev.docker
    192.168.99.100

## Server only

By default, root privileges are required because a DNS resolver configuration file has to be created at `/etc/resolver/docker`.
This is a working setup for Mac OSX, other OS are not yet supported out of the box.

If you're not on Mac OSX or don't want to run the program as root, you can pass the `-server-only` flag :

    > docker-machine-dns -server-only -port 10053
    2016/02/18 10:40:43 Listening on :10053...

*In that case you should also use a high port, the default port (53) would most likely require root privileges.*

To make sure the server runs correctly :

    > dig @localhost -p 10053 dev.docker +short
    192.168.99.100
