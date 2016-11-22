package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/miekg/dns"
)

type Dockness struct {
	Debug  bool
	Tld    string
	Ttl    uint
	Server *dns.Server
	Client *libmachine.Client
}

func (dockness *Dockness) Log(msg string) {
	if dockness.Debug {
		log.Println(msg)
	}
}

func (dockness *Dockness) Listen() error {
	dns.HandleFunc(dockness.Tld+".", dockness.lookup)
	if dockness.Server.PacketConn != nil {
		return dockness.Server.ActivateAndServe()
	}
	return dockness.Server.ListenAndServe()
}

// todo : collect errors and return a multierror
func (dockness *Dockness) Shutdown() error {
	dns.HandleRemove(dockness.Tld + ".")
	dockness.Client.Close()

	return dockness.Server.Shutdown()
}

func (dockness *Dockness) lookup(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	// check if it's a question from type A
	if m.Question[0].Qtype != dns.TypeA {
		dockness.Log("Unsupported question type.")
		m.SetRcode(r, dns.RcodeNotImplemented)
		w.WriteMsg(m)
		return
	}

	// parse the domain name
	domLevels := strings.Split(m.Question[0].Name, ".")
	domLevelsLen := len(domLevels)
	if domLevelsLen < 3 {
		dockness.Log(fmt.Sprintf("Couldn't parse the DNS question '%s'.", m.Question[0].Name))
		m.SetRcode(r, dns.RcodeFormatError)
		w.WriteMsg(m)
		return
	}

	// lookup for machine config
	machineName := domLevels[len(domLevels)-3]
	machine, err := dockness.Client.Filestore.Load(machineName)
	if err != nil {
		dockness.Log(fmt.Sprintf("Couldn't load machine : %s.", err))
		m.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(m)
		return
	}

	// parse machine config
	var baseDriver drivers.BaseDriver
	err = json.Unmarshal(machine.RawDriver, &baseDriver)
	if err != nil {
		dockness.Log(fmt.Sprintf("Couldn't load driver '%s' : %s.", machine.DriverName, err))
		m.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(m)
		return
	}

	dockness.Log(fmt.Sprintf("Found IP %s for machine '%s'", baseDriver.IPAddress, machineName))

	rr := &dns.A{
		Hdr: dns.RR_Header{
			Name:   m.Question[0].Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    uint32(dockness.Ttl),
		},
		A: net.ParseIP(baseDriver.IPAddress).To4(),
	}

	m.Answer = append(m.Answer, rr)

	w.WriteMsg(m)
}

func main() {
	port := flag.String("port", "53", "Port to listen on")
	tld := flag.String("tld", "docker", "Top-level domain to use")
	ttl := flag.Uint("ttl", 0, "Time to Live for DNS records")
	debug := flag.Bool("debug", false, "Enable debugging")
	flag.Parse()

	dockness := &Dockness{
		Ttl:    *ttl,
		Tld:    *tld,
		Debug:  *debug,
		Client: libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir()),
		Server: &dns.Server{
			Addr: ":" + *port,
			Net:  "udp",
		},
	}
	defer dockness.Shutdown()

	go log.Fatal(dockness.Listen())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
