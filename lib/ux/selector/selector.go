package selector

import (
	"os"

	"github.com/christiangelone/bang/lib/ux/print"
	"github.com/christiangelone/bang/source"
	"github.com/manifoldco/promptui"
)

type Selector struct{}

func New() *Selector {
	return &Selector{}
}

type bellSkipper struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal
// bell character.
func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriterCloser over os.Stderr.
func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}

var stringTemplates = &promptui.SelectTemplates{
	Label:    "    {{ . }}",
	Active:   `   {{ "⮀" }} {{ "•" | blue }} {{ .| cyan }}`,
	Inactive: `     {{ "•" | blue }} {{ . }}`,
	Selected: `    {{ "•" | blue }} Selection: {{ . | yellow }} {{ "✓" | green }}`,
	Help:     `    {{ "Use the arrow keys to navigate:" | faint }} {{ .NextKey | faint }} {{ .PrevKey | faint }}`,
}

func (s *Selector) Select(label string, options []string) (string, error) {
	prompt := promptui.Select{
		Label:     label,
		Items:     options,
		Templates: stringTemplates,
		Stdout:    &bellSkipper{},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *Selector) SelectSourceIndex(label string, options []source.Download) (int, error) {
	var templates = &promptui.SelectTemplates{
		Label:    "    {{ . }}",
		Active:   `   {{ "⮀" }} {{ "•" | blue }} {{ .RepoUrl | cyan }}`,
		Inactive: `     {{ "•" | blue }} {{ .RepoUrl }}`,
		Selected: `    {{ "•" | blue }} Selection: {{ .RepoUrl | yellow }} {{ "✓" | green }}`,
		Help:     `    {{ "Use the arrow keys to navigate:" | faint }} {{ .NextKey | faint }} {{ .PrevKey | faint }}`,
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     options,
		Templates: templates,
		Stdout:    &bellSkipper{},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return i, err
	}

	return i, nil
}

func (s *Selector) SelectIndex(label string, options []string) (int, error) {
	prompt := promptui.Select{
		Label:     label,
		Items:     options,
		Templates: stringTemplates,
		Stdout:    &bellSkipper{},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (s *Selector) SelectFailWith(text string) {
	print.Println(print.FgBlue, print.Padding+"•", print.FgNoColor, text, print.FgRed, "✗")
}
