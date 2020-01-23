package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/hpcloud/tail"
	"gopkg.in/gomail.v2"
)

func main() {
	from := flag.String("from", "logwatch@localhost", "Email FROM")
	to := flag.String("to", "", "Email TO")
	subject := flag.String("subject", "New log entry", "Email subject")
	file := flag.String("file", "", "Log file to watch")
	regex := flag.String("regex", "", "RegEx pattern for matching lines")

	flag.Parse()

	if *file == "" {
		fmt.Fprintln(os.Stderr, "Please specify a file to watch")
		os.Exit(1)
	}

	if *to == "" {
		fmt.Fprintln(os.Stderr, "Please specify a email recipient")
		os.Exit(1)
	}

	pattern, err := regexp.Compile(*regex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not compile your regex pattern: \n\t%v\n", err)
		os.Exit(1)
	}

	t, err := tail.TailFile(*file, tail.Config{
		Follow: true,
		ReOpen: true,
		Location: &tail.SeekInfo{
			Whence: os.SEEK_END,
		},
		Poll:   true,
		Logger: nil,
	})
	if err != nil {
		panic(err)
	}

	for line := range t.Lines {
		if !pattern.MatchString(line.Text) {
			continue
		}

		log.Printf("Sending email: %s", line.Text)
		m := gomail.NewMessage()
		m.SetHeader("From", *from)
		m.SetHeader("To", *to)
		m.SetHeader("Subject", *subject)
		m.SetBody("text/plain", line.Text)
		submitMail(m)
	}

}

func submitMail(m *gomail.Message) (err error) {
	cmd := exec.Command("/usr/sbin/sendmail", "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pw, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	var errs [3]error
	_, errs[0] = m.WriteTo(pw)
	errs[1] = pw.Close()
	errs[2] = cmd.Wait()
	for _, err = range errs {
		if err != nil {
			return
		}
	}
	return
}
