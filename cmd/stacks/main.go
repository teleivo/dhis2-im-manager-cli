package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	instance "github.com/teleivo/dhis2-im-manager-cli"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stdout, "Failed due to: %s\n", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	url := fs.String("url", "", "Instance manager URL")
	user := fs.String("user", "", "User to login and perform actions on the instance manager")
	pw := fs.String("pw", "", "Password of user")
	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}
	if *url == "" || *user == "" || *pw == "" {
		return errors.New("url, user and pw are required")
	}

	// TODO set some timeouts
	client := &http.Client{}
	im := instance.NewManager(*url, *user, *pw, client)
	err = im.Login()
	if err != nil {
		return err
	}

	m := instance.NewStacks(im)
	p := tea.NewProgram(m, tea.WithAltScreen())

	_ = out
	return p.Start()
}
