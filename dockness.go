package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/miekg/dns"
	//"github.com/pkg/profile"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

var api *libmachine.Client
var ttl uint

func findIp(machineName string) (string, error) {
	machine, err := api.Filestore.Load(machineName)
	if err != nil {
		return "", fmt.Errorf("couldn't load machine : %s", err)
	}

	var baseDriver drivers.BaseDriver
	err = json.Unmarshal(machine.RawDriver, &baseDriver)
	if err != nil {
		return "", fmt.Errorf("couldn't load driver %s", machine.DriverName)
	}

	return baseDriver.IPAddress, nil
}

func lookup(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	if m.Question[0].Qtype != dns.TypeA {
		return
	}

	domLevels := strings.Split(m.Question[0].Name, ".")
	domLevelsLen := len(domLevels)
	if domLevelsLen < 3 {
		log.Printf("Couldn't parse the DNS question '%s'", m.Question[0].Name)
		return
	}

	machineName := domLevels[len(domLevels)-3]
	ip, err := findIp(machineName)
	if err != nil {
		log.Printf("Couldn't find IP for machine '%s' : %s", machineName, err)
	} // else {
	//log.Printf("Found IP %s for machine '%s'", ip, machineName)
	//}

	rr := &dns.A{
		Hdr: dns.RR_Header{
			Name:   m.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    uint32(ttl),
		},
		A: net.ParseIP(ip).To4(),
	}

	m.Answer = append(m.Answer, rr)

	w.WriteMsg(m)
}

func main() {
	tld := flag.String("tld", "docker", "Top-level domain to use")
	flag.UintVar(&ttl, "ttl", 0, "Time to Live for DNS records")
	port := flag.String("port", "53", "Port to listen on")
	_ = flag.Bool("debug", false, "Enable debugging")
	flag.Parse()

	//if *debug {
	//	p := profile.Start(profile.MemProfile, profile.ProfilePath("."))
	//	defer p.Stop()
	//}

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
