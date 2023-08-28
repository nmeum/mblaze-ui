package main

import (
	"bufio"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	Seen MailFlag = iota
	Flagged
)

const (
	// Output format used by mscan(1) (passed via the -f flag).
	mscanFmt = "%n %S"
)

var (
	// POSIX extended regular expression for parsing 'mscanFmt'.
	mscanRegex = regexp.MustCompilePOSIX("^([0-9]+) (.+)$")
)

type Mail struct {
	ID      uint
	Subject string
}

func (m Mail) Show() error {
	// Use custom command-line options for less to ensure
	// the pager doesn't exit if the output fits on the screen.
	//
	// See also: https://github.com/leahneukirchen/mblaze/blob/v1.2/mshow.c#L818-L822
	pager := os.Getenv("PAGER")
	if pager == "" || strings.HasPrefix(pager, "less") {
		pager = "less --RAW-CONTROL-CHARS"
	}

	cmd := exec.Command("mshow", m.CmdArg())
	cmd.Env = append(os.Environ(), "MBLAZE_PAGER="+pager)

	// Make sure that we use {stdout,stdin,stderr} of the parent
	// process. Need to this explicitly when using os/exec.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (m Mail) Flag(flag MailFlag) error {
	cmd := exec.Command("mflag", flag.CmdOpt(), m.CmdArg())
	return cmd.Run()
}

func (m Mail) CmdArg() string {
	return strconv.FormatUint(uint64(m.ID), 10)
}

type MailFlag int

func (f MailFlag) CmdOpt() string {
	switch f {
	case Seen:
		return "-S"
	case Flagged:
		return "-F"
	}

	panic("unreachable")
}

func mblaze_mscan() ([]Mail, error) {
	var mails []Mail

	cmd := exec.Command("mscan", "-f", mscanFmt, "1:-1")
	cmd.Env = append(os.Environ(), "MBLAZE_PAGER=")

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
