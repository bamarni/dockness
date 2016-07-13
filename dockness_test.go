package main

import (
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"net"
	"reflect"
	"testing"
)

func CreateDockness() (*Dockness, string) {
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		return nil, ""
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
	dockness, addr := CreateDockness()
	if dockness == nil {
		t.Fatal("couldn't create dockness")
	}
	defer dockness.Shutdown()

	req := new(dns.Msg)
	req.SetQuestion("test."+dockness.Tld+".", dns.TypeA)

	client := new(dns.Client)
	resp, _, err := client.Exchange(req, addr)
	assert.NoError(t, err)

	if 1 != len(resp.Answer) {
		t.Fatal("expected an answer")
	}

	respA, ok := resp.Answer[0].(*dns.A)
	if !ok {
		t.Fatal("expected an A record")
	}
	assert.Equal(t, respA.A, net.IPv4(192, 0, 2, 0).To4())
}

func TestUnexistingMachine(t *testing.T) {
	dockness, addr := CreateDockness()
	if dockness == nil {
		t.Fatal("couldn't create dockness")
	}
	defer dockness.Shutdown()

	req := new(dns.Msg)
	req.SetQuestion("lochness."+dockness.Tld+".", dns.TypeA)

	client := new(dns.Client)
	resp, _, err := client.Exchange(req, addr)
	assert.NoError(t, err)

	assert.Equal(t, resp.Rcode, dns.RcodeNameError)
}

func TestInvalidQuestionType(t *testing.T) {
	dockness, addr := CreateDockness()
	if dockness == nil {
		t.Fatal("couldn't create dockness")
	}
	defer dockness.Shutdown()

	req := new(dns.Msg)
	req.SetQuestion("invalid."+dockness.Tld+".", dns.TypeTXT)

	client := new(dns.Client)
	resp, _, err := client.Exchange(req, addr)
	assert.NoError(t, err)

	assert.Equal(t, resp.Rcode, dns.RcodeNotImplemented)
}

func TestInvalidQuestionZone(t *testing.T) {
	dockness, addr := CreateDockness()
	if dockness == nil {
		t.Fatal("couldn't create dockness")
	}
	defer dockness.Shutdown()

	req := new(dns.Msg)
	req.SetQuestion(dockness.Tld+".", dns.TypeA)

	client := new(dns.Client)
	resp, _, err := client.Exchange(req, addr)
	assert.NoError(t, err)

	assert.Equal(t, resp.Rcode, dns.RcodeFormatError)
}

func BenchmarkServer(b *testing.B) {
	dockness, addr := CreateDockness()
	if dockness == nil {
		b.Fatal("couldn't create dockness")
	}
	defer dockness.Shutdown()

	req := new(dns.Msg)
	req.SetQuestion("test."+dockness.Tld+".", dns.TypeA)

	client := new(dns.Client)
	expectedIp := net.IPv4(192, 0, 2, 0).To4()

	for i := 0; i < b.N; i++ {
		resp, _, _ := client.Exchange(req, addr)
		respA, _ := resp.Answer[0].(*dns.A)
		if !reflect.DeepEqual(expectedIp, respA.A) {
			b.Fail()
		}
	}
}
