package instance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Manager struct {
	url    string
	user   string
	pw     string
	client *http.Client
	token  string
}

func NewManager(URL, user, pw string, client *http.Client) *Manager {
	return &Manager{
		url:    URL,
		user:   user,
		pw:     pw,
		client: client,
	}
}

type tokenBody struct {
	Token string `json:"access_token"`
}

func (m *Manager) Login() error {
	req, err := http.NewRequest(http.MethodPost, m.url+"/tokens", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(m.user, m.pw)
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("login failed: expected HTTP status 201, got %s", resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	tb := &tokenBody{}
	if err := d.Decode(tb); err != nil {
		return err
	}

	if tb.Token == "" {
		return errors.New("login failed: token is empty")
	}
	m.token = tb.Token

	return nil
}

type createBody struct {
	Name    string `json:"name"`
	GroupID int    `json:"groupId"`
	StackID int    `json:"stackID"`
}

// TODO the parameters are probably interesting/worth getting
type createRespBody struct {
	ID      int    `json:"ID"`
	Name    string `json:"name"`
	GroupID int    `json:"groupId"`
	StackID int    `json:"stackID"`
}

func (m *Manager) Create(name string, group, stack int) error {
	c := &createBody{
		Name:    name,
		GroupID: group,
		StackID: stack,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, m.url+"/instances", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	// TODO assumes Login() was called/and token is not expired
	req.Header.Add("Authorization", "Bearer "+m.token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create failed: expected HTTP status 201, got %s", resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	cb := &createRespBody{}
	if err := d.Decode(cb); err != nil {
		return err
	}

	// TODO remove
	fmt.Printf("%+v\n", cb)

	return nil
}

type optionalParam struct {
	ID           int    `json:"ID"`
	Name         string `json:"Name"`
	DefaultValue string `json:"DefaultValue"`
}

type requiredParam struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

// TODO Instances? and parameters
type stackBody struct {
	ID             int             `json:"ID"`
	Name           string          `json:"name"`
	OptionalParams []optionalParam `json:"optionalParameters"`
	RequiredParams []requiredParam `json:"requiredParameters"`
}

func (m *Manager) Stack(id int) error {
	req, err := http.NewRequest(http.MethodGet, m.url+"/stacks/"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}
	// TODO assumes Login() was called/and token is not expired
	req.Header.Add("Authorization", "Bearer "+m.token)
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching stack failed: expected HTTP status 200, got %s", resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	sb := &stackBody{}
	if err := d.Decode(sb); err != nil {
		return err
	}

	// TODO remove
	fmt.Printf("%+v\n", sb)

	return nil
}

type stacksBody struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
}

func (m *Manager) Stacks() error {
	req, err := http.NewRequest(http.MethodGet, m.url+"/stacks/", nil)
	if err != nil {
		return err
	}
	// TODO assumes Login() was called/and token is not expired
	req.Header.Add("Authorization", "Bearer "+m.token)
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching stacks failed: expected HTTP status 200, got %s", resp.Status)
	}

	d := json.NewDecoder(resp.Body)
	// TODO
	var sb []stacksBody
	if err := d.Decode(&sb); err != nil {
		return err
	}

	// TODO remove
	fmt.Printf("%+v\n", sb)

	return nil
}
