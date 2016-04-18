package main

import (
	"flag"
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	mcnlog "github.com/docker/machine/libmachine/log"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

var api *libmachine.Client
var ttl uint
var user string

func lookup(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)

	m.SetReply(r)

	var rr dns.RR
	for _, q := range m.Question {
		if q.Qtype != dns.TypeA {
			continue
		}

		domLevels := strings.Split(q.Name, ".")
		domLevelsLen := len(domLevels)
		if domLevelsLen < 3 {
			log.Printf("Couldn't parse the DNS question '%s'", q.Name)
			continue
		}
		machineName := domLevels[len(domLevels)-3]

		machine, err := api.Load(machineName)
		if err != nil {
			log.Printf("Couldn't load machine '%s' : %s", machineName, err)
			continue
		}

		ip, err := machine.Driver.GetIP()
		if err != nil {
			log.Printf("Couldn't find IP for machine '%s' : %s", machineName, err)
			continue
		}

		rr = &dns.A{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(ttl),
			},
			A: net.ParseIP(ip).To4(),
		}

		m.Answer = append(m.Answer, rr)
	}

	w.WriteMsg(m)
}

func main() {
	tld := flag.String("tld", "docker", "Top-level domain to use")
	flag.UintVar(&ttl, "ttl", 0, "Time to Live for DNS records")
	port := flag.String("port", "53", "Port to listen on")
	debug := flag.Bool("debug", false, "Enable debugging")
	flag.Parse()

	mcnlog.SetDebug(*debug)
	api = libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())
	defer api.Close()

	addr := ":" + *port
	server := &dns.Server{
		Addr: addr,
		Net:  "udp",
	}

	dns.HandleFunc(*tld+".", lookup)

	log.Printf("Listening on %s...", addr)
	go log.Fatal(server.ListenAndServe())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
