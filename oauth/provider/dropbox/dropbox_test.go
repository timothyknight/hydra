package dropbox

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mock = &dropbox{
	id: "123",
	conf: &oauth2.Config{
		ClientID:     "client",
		ClientSecret: "secret",
		RedirectURL:  "/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "/oauth2/authorize",
			TokenURL: "/oauth2/token",
		},
	},
}

func TestNew(t *testing.T) {
	m := New("321", "client", "secret", "/callback")
	assert.Equal(t, "321", m.id)
	assert.Equal(t, "client", m.conf.ClientID)
	assert.Equal(t, "secret", m.conf.ClientSecret)
	assert.Equal(t, "/callback", m.conf.RedirectURL)
}
func TestGetID(t *testing.T) {
	assert.Equal(t, "123", mock.GetID())
}

func TestGetAuthCodeURL(t *testing.T) {
	require.NotEmpty(t, mock.GetAuthenticationURL("state"))
}

func TestExchangeCode(t *testing.T) {

	router := mux.NewRouter()
	router.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"access_token": "ABCDEFG", "token_type": "bearer", "uid": "12345"}`)
	})
	router.HandleFunc("/users/get_current_account", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"account_id": "dbid:2qrw3etsdtr","name": {"given_name": "Peter","surname": "Peter","familiar_name": "Peter","display_name": "Peter"},"email": "peter@gmail.com","country": "DE","locale": "de","referral_link": "https://db.tt/w34setrdgxf","is_paired": false,"account_type": {".tag": "pro"}}`)
	})
	ts := httptest.NewServer(router)

	mock.api = ts.URL
	mock.conf.Endpoint.TokenURL = ts.URL + mock.conf.Endpoint.TokenURL

	t.Logf("Token URL: %s", mock.conf.Endpoint.TokenURL)
	t.Logf("API URL: %s", mock.api)
	code := "testcode"
	ses, err := mock.Exchange(code)
	require.Nil(t, err)
	assert.Equal(t, "dbid:2qrw3etsdtr", ses.GetRemoteSubject())
}
