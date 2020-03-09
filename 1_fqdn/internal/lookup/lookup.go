package lookup

func Get(host, dnsAddr string) ([]Result, error) {
	var results []Result
	var cfqdn = host
	for {
		cnames, err := factory(typeCName)(cfqdn, dnsAddr)
		if err == nil && len(cnames) > 0 {
			cfqdn = cnames[0]
			continue
		}
		ips, err := factory(typeA)(cfqdn, dnsAddr)
		if err != nil {
			//there are no A records for hostname
			return results, err
		}
		for _, ip := range ips {
			results = append(results, Result{ip, host})
		}
		break
	}
	return results, nil
}
