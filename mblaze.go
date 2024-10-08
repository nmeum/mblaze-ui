package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	Flagged MailFlag = iota
	Unflagged
	Seen
	Unseen
	Trashed
	Untrashed
)

const (
	// Output format used by mscan(1) (passed via the -f flag).
	mscanFmt = "%0R	%19D <%0f> %0S"

	// Maximum amount of characters to output for the from header.
	maxFrom = 17
)

var (
	// POSIX extended regular expression for parsing 'mscanFmt'.
	mscanRegex = regexp.MustCompilePOSIX(`^([^\t]+)	([0-9]+-[0-9]+-[0-9]+ [0-9][0-9]:[0-9][0-9]:[0-9][0-9]) <([^>]+)> (.+)$`)
)

type Mail struct {
	Path    string
	Date    time.Time
	From    string
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

	cmd := exec.Command("mshow", m.Path)
	cmd.Env = append(os.Environ(), "MBLAZE_PAGER="+pager)

	// Make sure that we use {stdout,stdin,stderr} of the parent
	// process. Need to this explicitly when using os/exec.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (m Mail) Reply() error {
	cmd := exec.Command("mrep", m.Path)

	// Make sure that we use {stdout,stdin,stderr} of the parent
	// process. Need to this explicitly when using os/exec.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (m Mail) Flag(flag MailFlag) error {
	cmd := exec.Command("mflag", flag.CmdOpt(), m.Path)
	return cmd.Run()
}

func (m Mail) String() string {
	from := m.From[0:min(len(m.From), maxFrom)]
	date := adaptiveTime(m.Date)

	out := fmt.Sprintf("%10s %17s %s", date, from, m.Subject)
	return out
}

type MailFlag int

func (f MailFlag) CmdOpt() string {
	switch f {
	case Unflagged:
		return "-f"
	case Flagged:
		return "-F"
	case Unseen:
		return "-s"
	case Seen:
		return "-S"
	case Untrashed:
		return "t"
	case Trashed:
		return "T"
	}

	panic("unreachable")
}

func mscan() ([]Mail, error) {
	var mails []Mail

	cmd := exec.Command("mscan", "-f", mscanFmt, "1:-1")
	cmd.Env = append(os.Environ(), "MBLAZE_PAGER=")
	// TODO: Somehow configure mblaze to not strip at all.
	cmd.Env = append(cmd.Env, "COLUMNS=99999")

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

		fp := subs[1]
		date, err := time.Parse(time.DateTime, subs[2])
		if err != nil {
			return mails, err
		}
		from := strings.TrimSpace(subs[3])
		subject := subs[4]

		mails = append(mails, Mail{
			Path:    fp,
			Date:    date,
			From:    from,
			Subject: subject,
		})
	}

	err = cmd.Wait()
	if err != nil {
		return mails, err
	}

	if len(mails) == 0 {
		return mails, fmt.Errorf("current sequence is empty")
	}

	return mails, nil
}
