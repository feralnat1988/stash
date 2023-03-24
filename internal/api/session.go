package api

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/session"
)

const returnURLParam = "returnURL"

func getLoginPage(loginUIBox fs.FS) []byte {
	data, err := fs.ReadFile(loginUIBox, "login.html")
	if err != nil {
		panic(err)
	}
	return data
}

type loginTemplateData struct {
	URL   string
	Error string
}

func serveLoginPage(loginUIBox fs.FS, w http.ResponseWriter, r *http.Request, returnURL string, loginError string) {
	loginPage := string(getLoginPage(loginUIBox))
	prefix := getProxyPrefix(r)
	loginPage = strings.ReplaceAll(loginPage, "/%BASE_URL%", prefix)

	templ, err := template.New("Login").Parse(loginPage)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = templ.Execute(w, loginTemplateData{URL: returnURL, Error: loginError})
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), http.StatusInternalServerError)
	}
}

func handleLogin(loginUIBox fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		returnURL := r.URL.Query().Get(returnURLParam)

		if !config.GetInstance().HasCredentials() {
			if returnURL != "" {
				http.Redirect(w, r, returnURL, http.StatusFound)
			} else {
				prefix := getProxyPrefix(r)
				http.Redirect(w, r, prefix+"/", http.StatusFound)
			}
			return
		}

		serveLoginPage(loginUIBox, w, r, returnURL, "")
	}
}

func handleLoginPost(loginUIBox fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue(returnURLParam)
		if url == "" {
			url = getProxyPrefix(r) + "/"
		}

		err := manager.GetInstance().SessionStore.Login(w, r)
		if errors.Is(err, session.ErrInvalidCredentials) {
			// serve login page with an error
			serveLoginPage(loginUIBox, w, r, url, "Username or password is invalid")
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}

func handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := manager.GetInstance().SessionStore.Logout(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirect to the login page if credentials are required
		prefix := getProxyPrefix(r)
		if config.GetInstance().HasCredentials() {
			http.Redirect(w, r, prefix+"/login", http.StatusFound)
		} else {
			http.Redirect(w, r, prefix+"/", http.StatusFound)
		}
	}
}
