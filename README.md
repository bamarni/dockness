# docker-machine-dns

DNS for Docker machines, allows to access them with the following domain format : `{machine}.docker`.

## How it works

It spins up a simplistic DNS server, only listening for questions about A records.

Behind the scene it will run `docker-machine ip {machine}` in order to resolve the IP address of a given machine.

## Installation

### Prebuilt binaries

Prebuilt binaries are available in the [releases](https://github.com/bamarni/docker-machine-dns/releases).

### From source (requires Go)

    go get github.com/bamarni/docker-machine-dns

## Usage

### Mac OSX

To develop on Mac you probably have a local VM, using VirtualBox for example.
However this machine gets assigned a dynamic IP address.

The program can be up and running in one command :

    > sudo docker-machine-dns
    2016/02/18 10:39:52 Creating configuration file at /etc/resolver/docker...
    2016/02/18 10:39:52 Listening on :53...

In another terminal, to make sure it works properly (let's say for a machine called `dev`) :

    > dig @localhost dev.docker +short
    192.168.99.100

## Linux

Even though Linux users might not need a development VM, it can be useful for remote machines
(eg. `ssh staging2.docker`, ...).

Here and unlike Mac OSX, there is no quick trick to make it work out of the box. Linux distributions come
with more advanced DNS management where manually tweaking `resolv.conf` is usually not an option.

What should work in most cases is to use [Dnsmasq](http://www.thekelleys.org.uk/dnsmasq/doc.html),
which provides a lightweight DNS server.

If you have it installed, you can run :

    > sudo docker-machine-dns -port 10053
    2016/02/18 10:40:43 Listening on :10053...

Then let Dnsmasq now about the resolver by running those commands :

    > echo "server=/docker/127.0.0.1#10053" | sudo tee -a /etc/dnsmasq.conf
    > sudo /etc/init.d/dnsmasq restart

### Usage reference

    docker-machine-dns [options...]

    Options:
      -port         Port to listen on (defaults to "53")
      -server-only  Server only, doesn't try to create a resolver configuration file
      -user         Execute the "docker-machine ip" command as a different user (defaults to "SUDO_USER")
