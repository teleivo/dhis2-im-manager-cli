package instance

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Manager struct {
	instanceURL string
	user        string
	pw          string
	client      *http.Client
	token       string
}

func NewManager(host, user, pw string, client *http.Client) *Manager {
	return &Manager{
		instanceURL: host,
		user:        user,
		pw:          pw,
		client:      client,
	}
}

type tokenBody struct {
	Token string `json:"access_token"`
}

func (m *Manager) Login() error {
	req, err := http.NewRequest(http.MethodPost, m.instanceURL+"/tokens", nil)
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

	// TODO remove
	fmt.Println(tb.Token)

	return nil
}
