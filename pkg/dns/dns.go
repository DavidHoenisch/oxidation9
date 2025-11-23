package dns

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Dns struct {
	domain  string
	rType   string
	dnsType uint16
	server  string
	tcp     bool
}

func (d *Dns) Handler(args []string) error {
	err := d.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	return d.Exec()
}

func (d *Dns) Parse(args []string) error {

	fs := flag.NewFlagSet("Dns", flag.ContinueOnError)

	fs.StringVar(&d.domain, "domain", "", "target domain")
	fs.StringVar(&d.rType, "rtype", "A", "record type to look up")
	fs.StringVar(&d.server, "server", "1.1.1.1:53", "target dns server")
	fs.BoolVar(&d.tcp, "tpc", false, "use tcp instead of udp")

	fs.Usage = d.Doc

	err := fs.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			return err
		}
		return fmt.Errorf("error parsing args %v", err)
	}

	if d.domain == "" {
		return fmt.Errorf("missing -domain parameter")
	}

	return nil
}

func (d *Dns) Exec() error {

	d.domain = dns.Fqdn(d.domain)
	recordType := strings.ToUpper(d.rType)
	d.server = d.server

	dnsType, ok := dns.StringToType[recordType]
	if !ok {
		return fmt.Errorf("could not parse dns type")
	}

	d.dnsType = dnsType

	if !strings.Contains(d.server, ":") {
		d.server = net.JoinHostPort(d.server, "53")
	}

	switch d.rType {
	case "A":
		return a(*d)
	case "AAAA":
		return aaaa(*d)
	case "TXT":
		return txt(*d)
	case "MX":
		return mx(*d)
	case "CNAME":
		return cname(*d)
	case "NS":
		return ns(*d)
	case "SOA":
		return soa(*d)
	}

	return nil
}

func (d *Dns) Doc() {
	fmt.Fprintf(os.Stderr, "Usage of dns:\n")
	fmt.Fprintf(os.Stderr, "  This tool is used for getting dns info.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	// Note: We can't print defaults easily here because 'fs' is local to Parse.
	// Usually, you define flags at the struct level or pass fs here if you want automatic defaults printing.
	fmt.Fprintf(os.Stderr, "  -domain    Target domain (required)\n")
	fmt.Fprintf(os.Stderr, "  -rtype     Dns record to retrieve (default A)\n")
	fmt.Fprintf(os.Stderr, "  -server    Upstream dns server to use (default 1.1.1.1)\n")
	fmt.Fprintf(os.Stderr, "  -tcp       Whether to send request over TCP instead of UDP (default false)\n")

}

func a(d Dns) error {

	msg := new(dns.Msg)
	msg.SetQuestion(d.domain, d.dnsType)
	msg.RecursionDesired = true

	c := new(dns.Client)
	c.Timeout = 5 * time.Second
	if d.tcp {
		c.Net = "tcp"
	}

	r, rtt, err := c.Exchange(msg, d.server)
	if err != nil {
		return fmt.Errorf("Dns exchange failed")
	}

	if r.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("Dns query failed with Rcode: %s", dns.RcodeToString[r.Rcode])
	}

	fmt.Printf("latency %v\n", rtt)
	fmt.Println("-------------------------------------------------")

	if len(r.Answer) == 0 {
		fmt.Println("No answer records found")
	}

	for _, ans := range r.Answer {
		fmt.Println(strings.ReplaceAll(ans.String(), "\t", "   "))
	}

	if len(r.Ns) > 0 {
		fmt.Println("\n -- Authority Section --")
		for _, ns := range r.Ns {
			fmt.Println(strings.ReplaceAll(ns.String(), "\t", "    "))
		}
	}

	if len(r.Extra) > 0 {
		fmt.Println("\n -- Extra Section --")
		for _, extra := range r.Extra {
			fmt.Println(strings.ReplaceAll(extra.String(), "\n", "   "))
		}
	}

	return nil
}

func aaaa(d Dns) error {

	msg := new(dns.Msg)
	msg.SetQuestion(d.domain, d.dnsType)
	msg.RecursionDesired = true

	c := new(dns.Client)
	c.Timeout = 5 * time.Second
	if d.tcp {
		c.Net = "tcp"
	}

	r, rtt, err := c.Exchange(msg, d.server)
	if err != nil {
		return fmt.Errorf("Dns exchange failed")
	}

	if r.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("Dns query failed with Rcode: %s", dns.RcodeToString[r.Rcode])
	}

	fmt.Printf("latency %v\n", rtt)
	fmt.Println("-------------------------------------------------")

	if len(r.Answer) == 0 {
		fmt.Println("No answer records found")
	}

	for _, ans := range r.Answer {
		fmt.Println(strings.ReplaceAll(ans.String(), "\t", "   "))
	}

	if len(r.Ns) > 0 {
		fmt.Println("\n -- Authority Section --")
		for _, ns := range r.Ns {
			fmt.Println(strings.ReplaceAll(ns.String(), "\t", "    "))
		}
	}

	if len(r.Extra) > 0 {
		fmt.Println("\n -- Extra Section --")
		for _, extra := range r.Extra {
			fmt.Println(strings.ReplaceAll(extra.String(), "\n", "   "))
		}
	}
	return nil
}

func txt(d Dns) error {

	msg := new(dns.Msg)
	msg.SetQuestion(d.domain, d.dnsType)
	msg.RecursionDesired = true

	c := new(dns.Client)
	c.Timeout = 5 * time.Second
	if d.tcp {
		c.Net = "tcp"
	}

	r, rtt, err := c.Exchange(msg, d.server)
	if err != nil {
		return fmt.Errorf("Dns exchange failed")
	}

	if r.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("Dns query failed with Rcode: %s", dns.RcodeToString[r.Rcode])
	}

	fmt.Printf("latency %v\n", rtt)
	fmt.Println("-------------------------------------------------")

	if len(r.Answer) == 0 {
		fmt.Println("No answer records found")
	}

	for _, ans := range r.Answer {
		fmt.Println(strings.ReplaceAll(ans.String(), "\t", "   "))
	}

	if len(r.Ns) > 0 {
		fmt.Println("\n -- Authority Section --")
		for _, ns := range r.Ns {
			fmt.Println(strings.ReplaceAll(ns.String(), "\t", "    "))
		}
	}

	if len(r.Extra) > 0 {
		fmt.Println("\n -- Extra Section --")
		for _, extra := range r.Extra {
			fmt.Println(strings.ReplaceAll(extra.String(), "\n", "   "))
		}
	}
	return nil
}

func mx(d Dns) error {

	msg := new(dns.Msg)
	msg.SetQuestion(d.domain, d.dnsType)
	msg.RecursionDesired = true

	c := new(dns.Client)
	c.Timeout = 5 * time.Second
	if d.tcp {
		c.Net = "tcp"
	}

	r, rtt, err := c.Exchange(msg, d.server)
	if err != nil {
		return fmt.Errorf("Dns exchange failed")
	}

	if r.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("Dns query failed with Rcode: %s", dns.RcodeToString[r.Rcode])
	}

	fmt.Printf("latency %v\n", rtt)
	fmt.Println("-------------------------------------------------")

	if len(r.Answer) == 0 {
		fmt.Println("No answer records found")
	}

	for _, ans := range r.Answer {
		fmt.Println(strings.ReplaceAll(ans.String(), "\t", "   "))
	}

	if len(r.Ns) > 0 {
		fmt.Println("\n -- Authority Section --")
		for _, ns := range r.Ns {
			fmt.Println(strings.ReplaceAll(ns.String(), "\t", "    "))
		}
	}

	if len(r.Extra) > 0 {
		fmt.Println("\n -- Extra Section --")
		for _, extra := range r.Extra {
			fmt.Println(strings.ReplaceAll(extra.String(), "\n", "   "))
		}
	}
	return nil
}

func cname(d Dns) error {

	msg := new(dns.Msg)
	msg.SetQuestion(d.domain, d.dnsType)
	msg.RecursionDesired = true

	c := new(dns.Client)
	c.Timeout = 5 * time.Second
	if d.tcp {
		c.Net = "tcp"
	}

	r, rtt, err := c.Exchange(msg, d.server)
	if err != nil {
		return fmt.Errorf("Dns exchange failed")
	}

	if r.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("Dns query failed with Rcode: %s", dns.RcodeToString[r.Rcode])
	}

	fmt.Printf("latency %v\n", rtt)
	fmt.Println("-------------------------------------------------")

	if len(r.Answer) == 0 {
		fmt.Println("No answer records found")
	}

	for _, ans := range r.Answer {
		fmt.Println(strings.ReplaceAll(ans.String(), "\t", "   "))
	}

	if len(r.Ns) > 0 {
		fmt.Println("\n -- Authority Section --")
		for _, ns := range r.Ns {
			fmt.Println(strings.ReplaceAll(ns.String(), "\t", "    "))
		}
	}

	if len(r.Extra) > 0 {
		fmt.Println("\n -- Extra Section --")
		for _, extra := range r.Extra {
			fmt.Println(strings.ReplaceAll(extra.String(), "\n", "   "))
		}
	}
	return nil
}

func ns(d Dns) error {

	msg := new(dns.Msg)
	msg.SetQuestion(d.domain, d.dnsType)
	msg.RecursionDesired = true

	c := new(dns.Client)
	c.Timeout = 5 * time.Second
	if d.tcp {
		c.Net = "tcp"
	}

	r, rtt, err := c.Exchange(msg, d.server)
	if err != nil {
		return fmt.Errorf("Dns exchange failed")
	}

	if r.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("Dns query failed with Rcode: %s", dns.RcodeToString[r.Rcode])
	}

	fmt.Printf("latency %v\n", rtt)
	fmt.Println("-------------------------------------------------")

	if len(r.Answer) == 0 {
		fmt.Println("No answer records found")
	}

	for _, ans := range r.Answer {
		fmt.Println(strings.ReplaceAll(ans.String(), "\t", "   "))
	}

	if len(r.Ns) > 0 {
		fmt.Println("\n -- Authority Section --")
		for _, ns := range r.Ns {
			fmt.Println(strings.ReplaceAll(ns.String(), "\t", "    "))
		}
	}

	if len(r.Extra) > 0 {
		fmt.Println("\n -- Extra Section --")
		for _, extra := range r.Extra {
			fmt.Println(strings.ReplaceAll(extra.String(), "\n", "   "))
		}
	}
	return nil
}

func soa(d Dns) error {

	msg := new(dns.Msg)
	msg.SetQuestion(d.domain, d.dnsType)
	msg.RecursionDesired = true

	c := new(dns.Client)
	c.Timeout = 5 * time.Second
	if d.tcp {
		c.Net = "tcp"
	}

	r, rtt, err := c.Exchange(msg, d.server)
	if err != nil {
		return fmt.Errorf("Dns exchange failed")
	}

	if r.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("Dns query failed with Rcode: %s", dns.RcodeToString[r.Rcode])
	}

	fmt.Printf("latency %v\n", rtt)
	fmt.Println("-------------------------------------------------")

	if len(r.Answer) == 0 {
		fmt.Println("No answer records found")
	}

	for _, ans := range r.Answer {
		fmt.Println(strings.ReplaceAll(ans.String(), "\t", "   "))
	}

	if len(r.Ns) > 0 {
		fmt.Println("\n -- Authority Section --")
		for _, ns := range r.Ns {
			fmt.Println(strings.ReplaceAll(ns.String(), "\t", "    "))
		}
	}

	if len(r.Extra) > 0 {
		fmt.Println("\n -- Extra Section --")
		for _, extra := range r.Extra {
			fmt.Println(strings.ReplaceAll(extra.String(), "\n", "   "))
		}
	}
	return nil
}
