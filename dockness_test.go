package main

import (
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func CreateDockness(t *testing.T) (*Dockness, string) {
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Errorf("cannot listen for UDP packets")
	}

	dockness := &Dockness{
		//Debug:  true,
		Tld:    "docker",
		Client: libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir()),
		Server: &dns.Server{
			PacketConn: pc,
		},
	}

	go func() {
		dockness.Listen()
		pc.Close()
	}()

	return dockness, pc.LocalAddr().String()
}

func TestExistingMachine(t *testing.T) {
	dockness, addr := CreateDockness(t)
	defer dockness.Shutdown()

	req := new(dns.Msg)
	req.SetQuestion("test."+dockness.Tld+".", dns.TypeA)

	client := new(dns.Client)
	resp, _, err := client.Exchange(req, addr)
	assert.NoError(t, err)

	respA, ok := resp.Answer[0].(*dns.A)
	if !ok {
		t.Errorf("expected an A record")
	}
	assert.Equal(t, respA.A, net.IPv4(1, 2, 3, 4).To4())
	err = dockness.Shutdown()
	assert.NoError(t, err)
}

func TestUnexistingMachine(t *testing.T) {
	dockness, addr := CreateDockness(t)
	defer dockness.Shutdown()

	req := new(dns.Msg)
	req.SetQuestion("lochness."+dockness.Tld+".", dns.TypeA)

	client := new(dns.Client)
	resp, _, err := client.Exchange(req, addr)
	assert.NoError(t, err)

	assert.Equal(t, resp.Rcode, dns.RcodeNameError)
}

func TestInvalidQuestionType(t *testing.T) {
	dockness, addr := CreateDockness(t)
	defer dockness.Shutdown()

	req := new(dns.Msg)
	req.SetQuestion("invalid."+dockness.Tld+".", dns.TypeTXT)

	client := new(dns.Client)
	resp, _, err := client.Exchange(req, addr)
	assert.NoError(t, err)

	assert.Equal(t, resp.Rcode, dns.RcodeNotImplemented)
}

func TestInvalidQuestionZone(t *testing.T) {
	dockness, addr := CreateDockness(t)
	defer dockness.Shutdown()

	req := new(dns.Msg)
	req.SetQuestion(dockness.Tld+".", dns.TypeA)

	client := new(dns.Client)
	resp, _, err := client.Exchange(req, addr)
	assert.NoError(t, err)

	assert.Equal(t, resp.Rcode, dns.RcodeFormatError)
}
