//I realized one think here: you must block main thread when till  promise ends . Which brings another complexity
//You open at  one thread per promise, one thread per record  and promise often timeouts in this case ! (15s )was not enough in some runs
//execution time 17.387858418s
//*/
//
//package main
//
//import (
//	"bufio"
//	"flag"
//	"fmt"
//	"github.com/kuritka/gext/concurency"
//	"os"
//	"runtime"
//	"sync"
//	"time"
//
//	"github.com/kuritka/gext/guard"
//	"github.com/kuritka/gext/log"
//
//	"github.com/kuritka/dns/1_fqdn/internal/lookup"
//)
//
//var logger = log.Log
//
//var (
//	fDomain       = flag.String("domain", "", "The domain to perform guessing against")
//	fWordListPath = flag.String("wordlist", "", "The word list to use for guessing")
//	fDns          = flag.String("dns", "8.8.8.8:53", "The DNS server to use")
//)
//
//
//func main() {
//	start := time.Now()
//	flag.Parse()
//	if *fDomain == "" || *fWordListPath == "" {
//		logger.Fatal().Msgf("-domain and -wordlist required")
//	}
//	fmt.Println("checking dns")
//	runtime.GOMAXPROCS(10)
//
//	f, err := os.Open(*fWordListPath)
//	guard.FailOnError(err, "can't open file")
//	defer f.Close()
//
//	wg := new(sync.WaitGroup)
//	scanner := bufio.NewScanner(f)
//	for scanner.Scan() {
//		test :=  fmt.Sprintf("%s.%s", scanner.Text(), *fDomain)
//
//		call(test, *fDns).Then(
//			func(i interface{}) error {
//				wg.Add(1)
//				results := i.([]lookup.Result)
//				for _, r := range results {
//					fmt.Println(r.ToString())
//				}
//				wg.Done()
//				return nil
//			}, func(err error) {
//				logger.Err(err)
//		})
//	}
//	wg.Wait()
//
//	fmt.Printf("\nexecution time %s\n", time.Since(start))
//	_,_ = fmt.Scanln()
//}
//
//
//
//func call(hostname string, dns string) *concurency.Promise {
//	result := new(concurency.Promise)
//	result.SuccessChannel = make(chan interface{},1)
//	result.ErrorChannel = make(chan error,1)
//
//	go func() {
//		results, err := lookup.Get(hostname, dns)
//		if err != nil {
//			result.ErrorChannel <- err
//			return
//		}
//		result.SuccessChannel <- results
//	}()
//	return result
//}
//
