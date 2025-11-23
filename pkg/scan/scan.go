package scan

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Scan struct {
	domain string
}

func (s *Scan) Handler(args []string) error {

	err := s.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	return s.Exec()
}

func (s *Scan) Parse(args []string) error {

	fs := flag.NewFlagSet("Dns", flag.ContinueOnError)

	fs.StringVar(&s.domain, "domain", "", "target domain")

	fs.Usage = s.Doc

	err := fs.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			return err
		}
		return fmt.Errorf("error parsing args %v", err)
	}

	if s.domain == "" {
		return fmt.Errorf("missing -domain parameter")
	}

	return nil
}

func (s *Scan) Exec() error {
	targetPorts := 1024
	timeout := 500 * time.Millisecond
	concurrencyLimit := 1000

	// Channels
	ports := make(chan int, concurrencyLimit)
	results := make(chan int)
	var openports []int

	var wg sync.WaitGroup

	for i := 0; i < concurrencyLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(s.domain, ports, results, timeout)
		}()
	}

	go func() {
		for i := 1; i <= targetPorts; i++ {
			ports <- i
		}
		close(ports)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for port := range results {
		if port != 0 {
			openports = append(openports, port)
		}
	}

	sort.Ints(openports)

	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}

	return nil
}

func (s *Scan) Doc() {}

func worker(addr string, ports <-chan int, results chan<- int, timeout time.Duration) {
	for p := range ports {
		address := net.JoinHostPort(addr, strconv.Itoa(p))

		// CRITICAL CHANGE: DialTimeout
		// If it takes longer than 500ms, we kill it.
		conn, err := net.DialTimeout("tcp", address, timeout)

		if err != nil {
			// Port is closed or filtered
			results <- 0
			continue
		}

		conn.Close()
		results <- p
	}
}
