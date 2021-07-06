package progress

import (
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	. "github.com/christiangelone/bang/lib/sugar"
	"github.com/christiangelone/bang/lib/ux/print"
	"github.com/schollz/progressbar/v3"
)

func NewProgress() *Progress {
	delay := 80 * time.Millisecond
	s := spinner.New(
		[]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay,
		spinner.WithHiddenCursor(true),
	)
	s.Prefix = print.Padding
	s.Color("blue", "bold")

	return &Progress{
		spinner: s,
		delay:   delay,
		started: false,
	}
}

type Progress struct {
	bar     *progressbar.ProgressBar
	spinner *spinner.Spinner
	delay   time.Duration
	text    string
	started bool
}

func (s *Progress) Start(text string) {
	s.started = true
	s.Text(text)
	s.startBar()
	s.spinner.Start()
}

func (s *Progress) startBar() {
	Go(func() {
		for s.started {
			if s.bar != nil {
				s.progressText(s.bar.String())
			}
		}
	})
}

func (s *Progress) GetBar(cLength int64) io.Writer {
	s.bar = progressbar.NewOptions64(
		cLength,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWriter(ioutil.Discard),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionThrottle(s.delay),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]=[reset]",
			SaucerHead:    "[blue]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	s.bar.RenderBlank()
	return s.bar
}

func (s *Progress) progressText(progressStr string) {
	s.spinner.Suffix = " " + s.text + strings.Replace(progressStr, "\r", "", -1)
}

func (s *Progress) Text(text string) {
	s.text = text
	s.spinner.Suffix = " " + s.text
}

func (s *Progress) stop() {
	s.started = false
	s.spinner.Stop()
}

func (s *Progress) printSuccess(text, emoji string) {
	success := print.Sprint(print.FgNoColor, text, print.FgGreen, emoji)
	print.Bullet(success)
}

func (s *Progress) printFail(text, emoji string) {
	fail := print.Sprint(print.FgNoColor, text, print.FgRed, emoji)
	print.Bullet(fail)
}

func (s *Progress) StopSuccessWith(text string) {
	s.stop()
	if text == "" {
		s.printSuccess(s.text, "✓")
	} else {
		s.printSuccess(text, "✓")
	}
}

func (s *Progress) StopSuccess() {
	s.stop()
	s.printSuccess(s.text, "✓")
}

func (s *Progress) StopFailWith(text string) {
	s.stop()
	if text == "" {
		s.printFail(s.text, "✗")
	} else {
		s.printFail(text, "✗")
	}
}
