package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

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
	name := fs.String("name", "", "Name of the instance to be created/deployed")
	// TODO allow passing in the name of the group
	group := fs.Int("group", 2, "Group to perform actions on the instance manager")
	// TODO allow passing in the name of the stack
	stack := fs.Int("stack", 1, "Group to perform actions on the instance manager")
	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}
	if *url == "" || *user == "" || *pw == "" || *name == "" {
		return errors.New("url, user, pw and name are required")
	}

	// TODO set some timeouts
	client := &http.Client{}
	im := instance.NewManager(*url, *user, *pw, client)
	err = im.Login()
	if err != nil {
		return err
	}
	_ = group
	// err = im.Create(*name, *group, *stack)
	// if err != nil {
	// 	return err
	// }
	fmt.Println("stacks")
	err = im.Stacks()
	if err != nil {
		return err
	}
	fmt.Println("stack", *stack)
	err = im.Stack(*stack)
	if err != nil {
		return err
	}

	_ = out
	return nil
}
