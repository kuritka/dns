///*
//package tests main DNS functionality
//Before you run, find your ethernet interface by call
//
//`ifconfig -a`
//
//Normally it is eth0 but sometimes it could be something else i.e. ens4u1
//
//Than listen on port :53 by calling
//
//1sudo tcpdump -i ens4u1 -n udp port 53`
//
//goes from my ip to top google dns and
//16:07:05.464049 IP 192.168.0.123.50179 > 8.8.8.8.domain: 29435+ A? ulozto.cz. (27)
//goes back with A record of ulozto.cz
//16:07:05.482307 IP 8.8.8.8.domain > 192.168.0.123.50179: 29435 1/0/0 A 77.48.29.200 (43)
//
//
//https://github.com/kuritka/threading/blob/master/c_promises/main.go
//
//output:
//Logger configured
//checking dns
//52.157.177.204 asterix.onho.cz
//137.117.240.153 dev.onho.cz
//51.145.247.10 hello.onho.cz
//52.157.177.204 ne.onho.cz
//137.117.240.153 dev.onho.cz
//52.157.177.204 ne.onho.cz
//52.157.177.204 httpbin.onho.cz
//
//execution time 17.387858418s
//*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/kuritka/gext/guard"
	"github.com/kuritka/gext/log"

	"github.com/kuritka/dns/1_fqdn/internal/lookup"
)

const workerCount = 100

var logger = log.Log

var (
	fDomain       = flag.String("domain", "", "The domain to perform guessing against")
	fWordListPath = flag.String("wordlist", "", "The word list to use for guessing")
	fDns          = flag.String("dns", "8.8.8.8:53", "The DNS server to use")
)

var (
	hostCh = make(chan string)
)

func producer(wordListPath string, host string, hostCh chan<- string) {
	f, err := os.Open(wordListPath)
	guard.FailOnError(err, "can't open file")
	defer f.Close()
	defer close(hostCh)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hostCh <- fmt.Sprintf("%s.%s", scanner.Text(), host)
	}
}

func worker(hosts <-chan string, dnsServer string, wg *sync.WaitGroup) {
	wg.Add(1)
	for fqdn := range hosts {

		results, _ := lookup.Get(fqdn, dnsServer)
		for _, r := range results {
			fmt.Println(r.ToString())
		}
	}
	wg.Done()
}

func main() {
	start := time.Now()
	flag.Parse()
	if *fDomain == "" || *fWordListPath == "" {
		logger.Fatal().Msgf("-domain and -wordlist required")
	}
	fmt.Println("checking dns")
	runtime.GOMAXPROCS(10)
	wg := new(sync.WaitGroup)

	//specifying worker pool first
	for i := 0; i < workerCount; i++ {
		go worker(hostCh, *fDns, wg)
	}

	go producer(*fWordListPath, *fDomain, hostCh)

	wg.Wait()

	fmt.Printf("\nexecution time %s\n", time.Since(start))
	//x, err := lookup.Get(*fDomain,"8.8.8.8:53")
	//guard.FailOnError(err,"domain: %s",*fDomain)
	//for _,a := range x {
	//	fmt.Println(a.ToString())
	//}
}
