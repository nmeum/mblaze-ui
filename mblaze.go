package main

import (
	"bufio"
	"os/exec"
	"regexp"
	"strconv"
)

type Mail struct {
	ID      uint
	Subject string
}

const (
	// Output format used by mscan(1) (passed via the -f flag).
	mscanFmt = "%n %S"
)

var (
	// POSIX extended regular expression for parsing 'mscanFmt'.
	mscanRegex = regexp.MustCompilePOSIX("^([0-9]+) (.+)$")
)

func mblaze_mscan() ([]Mail, error) {
	var mails []Mail

	cmd := exec.Command("mscan", "-f", mscanFmt, "1:-1")
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return mails, err
	}
	defer reader.Close()

	err = cmd.Start()
	if err != nil {
		return mails, err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		subs := mscanRegex.FindStringSubmatch(scanner.Text())
		if subs == nil {
			// Message might have been moved since last mseq(1).
			// For example, if it was marked as seen via mflag(1).
			continue
		}

		var id uint64
		id, err = strconv.ParseUint(subs[1], 10, 32)
		if err != nil {
			return mails, err
		}
		subject := subs[2]

		mails = append(mails, Mail{uint(id), subject})
	}

	err = cmd.Wait()
	if err != nil {
		return mails, err
	}

	return mails, nil
}
