package blubber

import (
	"fmt"
	"net"
	"bytes"
	"errors"
	"time"
)

type workerJob struct {
	Target string
	Provider Provider
}

type workerResult struct {
	wid        int
	provider   string
	code	   string
	message    string
	time 		time.Duration
}

type ipVersion int
const NoIP ipVersion = 0
const IPv4 ipVersion = 4
const IPv6 ipVersion = 6

func checkIpVersion(ip net.IP) (ipVersion,error) {
	v4InV6Prefix := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}
	
	if ip == nil {
		return NoIP,errors.New("Not valid IP")
	}

	if bytes.HasPrefix(ip, v4InV6Prefix) {
		return IPv4,nil
	} else {
		return IPv6,nil
	}
}

func reverseSlice(s []byte, from int, to int) {
	for i, j := from, to; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}

func makeRequestProvider(p Provider, target string) {
	target = fmt.Sprintf("%s.%s", target, p.Hostname)
	fmt.Println(target)
}

func prepareTarget(target string) string {
	ip := net.ParseIP(target)
	ipVersion, err := checkIpVersion(ip)

	if err == nil {
		switch ipVersion {
			case IPv4:
				reverseSlice(ip, len(ip)-4, len(ip)-1)
			case IPv6:
				reverseSlice(ip, 0, len(ip)-1)
		}
		target = fmt.Sprintf("%s", ip)
	}

	return target
}

func (bl *Blubber)MakeRequest(target string, Nworkers int) {
	jobs := make(chan workerJob)
	results := make(chan workerResult)

	for w := 1; w <= Nworkers; w++ {
		go request_worker(w, jobs, results)
	}

	target = prepareTarget(target)

	go func() {
		for _,p := range bl.ProviderList {
			jobs <- workerJob{
					Target: target,
					Provider: p,
			}
		}
	}()

	for i := 0; i < len(bl.ProviderList); i++ {
		result := <- results
		fmt.Println(result)
	}
}

func request_worker(id int, jobs <-chan workerJob, results chan<- workerResult) {
	lid := 0
	for job := range jobs {
		lid = lid + 1
		result := workerResult{
			wid: id,
			provider: job.Provider.Hostname,
		}
		
		hostname := fmt.Sprintf("%s.%s", job.Target, job.Provider.Hostname)

		timeInit := time.Now()
		res, _ := net.LookupHost(hostname)
		if len(res) > 0 {
			result.code = res[0]
			restxt, _ := net.LookupTXT(hostname)
			if len(restxt) > 0 {
				result.message = restxt[0]
			}
		}
		timeFini := time.Now()

		result.time = time.Duration(timeFini.Nanosecond() - timeInit.Nanosecond())/time.Millisecond

		results <- result
	}
}