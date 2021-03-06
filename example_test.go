package dns

import (
	"fmt"
)

// Retrieve the MX records for miek.nl.
func ExampleRR_MX() {
	config, _ := ClientConfigFromFile("/etc/resolv.conf")
	c := new(Client)
	m := new(Msg)
	m.SetQuestion("miek.nl.", TypeMX)
	m.RecursionDesired = true
	r, err := c.Exchange(m, config.Servers[0]+":"+config.Port)
	if err != nil {
		return
	}
	if r.Rcode != RcodeSuccess {
		return
	}
	for _, a := range r.Answer {
		if mx, ok := a.(*RR_MX); ok {
			fmt.Printf("%s\n", mx.String())
		}
	}
}

// Retrieve the DNSKEY records of a zone and convert them
// to DS records for SHA1, SHA256 and SHA384.
func ExampleToDs(zone string) {
	config, _ := ClientConfigFromFile("/etc/resolv.conf")
	c := new(Client)
	m := new(Msg)
	if zone == "" {
		zone = "miek.nl"
	}
	m.SetQuestion(Fqdn(zone), TypeDNSKEY)
	m.SetEdns0(4096, true)
	r, err := c.Exchange(m, config.Servers[0]+":"+config.Port)
	if err != nil {
		return
	}
	if r.Rcode != RcodeSuccess {
		return
	}
	for _, k := range r.Answer {
		if key, ok := k.(*RR_DNSKEY); ok {
			for _, alg := range []int{SHA1, SHA256, SHA384} {
				fmt.Printf("%s; %d\n", key.ToDS(alg).String(), key.Flags)
			}
		}
	}
}

// Show how to setup the authors for 'authors.bind. CH TXT' or 'authors.server. CH  TXT'
// queries.
func ExampleAuthors() {
	// ... server setup is out of scope ...

	// Register the handle funcs.
	HandleFunc("authors.bind.", HandleAuthors)
	HandleFunc("authors.server.", HandleAuthors)

	// To extend the authors list, just append to dns.Authors (a []string)
	Authors = append(Authors, "G. I. Joe")

	// Or ...
	Authors = []string{"Just Me"}
}
