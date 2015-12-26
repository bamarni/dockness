package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"os/exec"
	"strings"
)

func lookup(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)

	m.SetReply(r)

	var rr dns.RR
	for _, q := range m.Question {
		if q.Qtype != dns.TypeA {
			continue
		}

		domLevels := strings.Split(q.Name, ".")
		machine := domLevels[len(domLevels)-3]
		stdoutBytes, err := exec.Command("docker-machine", "ip", machine).Output()
		if err != nil {
			continue
		}
		ip := string(stdoutBytes[:len(stdoutBytes)-1])

		rr = &dns.A{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			A: net.ParseIP(ip).To4(),
		}

		m.Answer = append(m.Answer, rr)
	}

	w.WriteMsg(m)
}

func main() {
	port := flag.String("port", "10053", "Port to listen on")
	flag.Parse()

	addr := ":"
	addr += *port

	server := &dns.Server{
		Addr: addr,
		Net:  "udp",
	}

	dns.HandleFunc(".", lookup)

	fmt.Printf("Listening on %s...", addr)
	log.Fatal(server.ListenAndServe())
}
