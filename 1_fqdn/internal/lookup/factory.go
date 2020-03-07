package lookup

import (
	"errors"

	"github.com/kuritka/gext/guard"

	"github.com/miekg/dns"
)


//https://www.sohamkamani.com/blog/golang/2018-06-20-golang-factory-patterns/
//function factory
func Factory(dnsType uint16) func(host, dnsAddr string) ([]string, error) {
	switch dnsType {
	case dns.TypeA:
		return lookupA
	case dns.TypeCNAME:
		return lookupCName
	default:
		guard.FailOnError(errors.New(""),"not implemented factory for such dns type")
	}
	return nil
}

//ulozto.cz,8.8.8.8:53
func lookupA(host, dnsAddr string) ([]string, error){
	var msg dns.Msg
	var ips []string
	fqdn := dns.Fqdn(host)
	msg.SetQuestion(fqdn, dns.TypeA)
	//8.8.8.8 is the primary DNS server for Google
	in, err := dns.Exchange(&msg, dnsAddr)
	if err != nil {
		return ips, err
	}
	if len(in.Answer) < 1 {
		return ips,  errors.New("no answer")
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips,a.A.String())
		}
	}
	return ips,nil
}

func lookupCName(host, dnsAddr string) ([]string, error){
	var msg dns.Msg
	var fqdns []string
	fqdn := dns.Fqdn(host)
	msg.SetQuestion(fqdn, dns.TypeCNAME)
	//8.8.8.8 is the primary DNS server for Google
	in, err := dns.Exchange(&msg, dnsAddr)
	if err != nil {
		return fqdns, err
	}
	if len(in.Answer) < 1 {
		return fqdns,  errors.New("no answer")
	}
	for _, answer := range in.Answer {
		if c, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns,c.Target)
		}
	}
	return fqdns,nil
}