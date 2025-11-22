package spam

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

type Spam struct {
	count int
	url   string
}

func (s *Spam) Handler(args []string) error {
	err := s.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	return s.Exec()
}

func (s *Spam) Parse(args []string) error {

	fs := flag.NewFlagSet("spam", flag.ContinueOnError)

	fs.IntVar(&s.count, "count", 10, "number of requests to make")
	fs.StringVar(&s.url, "url", "", "url to make requests to")

	fs.Usage = s.Doc

	err := fs.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			return err
		}
		return fmt.Errorf("error parsing args %v", err)
	}

	if s.url == "" {
		return fmt.Errorf("missing -url parameter")
	}

	if s.count < 1 {
		return fmt.Errorf("missing -count parameter")
	}

	return nil

}

func (s *Spam) Exec() error {

	for x := range s.count {
		fmt.Printf("================= REQUEST [%d] ================= \n", x)
		data, err := http.Get(s.url)
		if err != nil {
			return err
		}

		for name, values := range data.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)

			}
		}
	}

	return nil
}

func (s *Spam) Doc() {
	fmt.Fprintf(os.Stderr, "Usage of spam:\n")
	fmt.Fprintf(os.Stderr, "  This tool spams a URL to test headers.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	// Note: We can't print defaults easily here because 'fs' is local to Parse.
	// Usually, you define flags at the struct level or pass fs here if you want automatic defaults printing.
	fmt.Fprintf(os.Stderr, "  -url    Target URL (required)\n")
	fmt.Fprintf(os.Stderr, "  -count  Number of requests (default 10)\n")
}
