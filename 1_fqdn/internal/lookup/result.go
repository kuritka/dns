package lookup

import "fmt"

type Result struct {
	IPAddr string
	Host   string
}

func (r *Result) ToString() string {
	return fmt.Sprintf("%s %s", r.IPAddr, r.Host)
}
