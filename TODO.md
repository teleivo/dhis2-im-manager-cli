# TODO

* create a frame with a header (tabs), and footer (auth info and help)
* create help component like https://raw.githubusercontent.com/charmbracelet/bubbletea/master/examples/help/main.go
* create auth component showing user (group) and host

## Help component

* it should have a set of global key maps like ? for help and q for quitting
  or should these be defined in the frame and passed by it
  other keymaps should be passed by whatever is rendered in the frame

## Auth component

* should it be responsible for keeping me signed in? refreshing the token
* it could show a message when login/refresh is in progress

