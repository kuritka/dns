package main

/*
package tests main DNS functionality
Before you run, find your ethernet interface by call

`ifconfig -a`

Normally it is eth0 but sometimes it could be something else i.e. ens4u1

Than listen on port :53 by calling

1sudo tcpdump -i ens4u1 -n udp port 53`

goes from my ip to top google dns and
16:07:05.464049 IP 192.168.0.123.50179 > 8.8.8.8.domain: 29435+ A? ulozto.cz. (27)
goes back with A record of ulozto.cz
16:07:05.482307 IP 8.8.8.8.domain > 192.168.0.123.50179: 29435 1/0/0 A 77.48.29.200 (43)
*/
import (
	"flag"
	"fmt"
	"github.com/kuritka/dns/1_fqdn/internal/lookup"
	"github.com/kuritka/gext/guard"
	"github.com/kuritka/gext/log"
)


var logger = log.Log

var (
	fDomain   = flag.String("domain", "", "The domain to perform guessing against")
	fWordList = flag.String("wordlist", "", "The word list to use for guessing")
	fWorkers  = flag.Int("c", 100, "The number of workers to use.")
	fDns      = flag.String("dns", "8.8.8.8:53", "The DNS server to use")
)



func main() {
	flag.Parse()
	if *fDomain == "" || *fWordList=="" {
		logger.Fatal().Msgf("-domain and -wordlist required")
	}

	x, err := lookup.Lookup(*fDomain,"8.8.8.8:53")
	guard.FailOnError(err,"domain: %s",*fDomain)
	for _,a := range x {
		fmt.Println(a.ToString())
	}
}
