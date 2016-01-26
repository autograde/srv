package web

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/autograde/srv/config"
	"github.com/google/go-github/github"
	// "golang.org/x/oauth2/github"
)

// Constants and methods for using GitHub as OAuth provider for Autograder.

const (
	githubAppURL          = "https://github.com/settings/applications/new"
	githubTokenURL        = "https://github.com/login/oauth/access_token"
	githubAuthURL         = "https://github.com/login/oauth/authorize"
	githubAdminScope      = "admin:org,repo,admin:repo_hook"
	githubSessionName     = "auth_credentials"
	githubTokenSessionKey = "access_token"
)

var (
	// ErrNoAccessToken indicates that a empty access token was provided
	ErrNoAccessToken = errors.New("non-empty OAuth access token required")
)

// ghConns map from access token to the corresponding GitHub connection object.
var ghConns = make(map[string]*GitHubConn)

// GitHubConn contains connection details for GitHub.
type GitHubConn struct {
	// private fields will not be stored in the database
	accessToken string
	scope       string
	client      *github.Client
	user        *github.User
}

// NewGitHubConn creates a new GitHubConn and connects to GitHub and gets the
// GitHub user associated with the given access token.
func NewGitHubConn(token, scope string) (g *GitHubConn, err error) {
	if token == "" {
		return nil, ErrNoAccessToken
	}
	g = &GitHubConn{accessToken: token, scope: scope}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	g.client = github.NewClient(tc)
	g.user, _, err = g.client.Users.Get("")
	if err != nil {
		return nil, err
	}
	return g, nil
}

// oauthScopeURL returns a URL for redirecting to obtain a new scope.
// This is used to authenticate as teacher.
func oauthScopeURL() string {
	u, err := url.Parse(githubAuthURL)
	if err != nil {
		// the redirection URL must be a valid URL
		panic(err)
	}
	values := u.Query()
	values.Set("client_id", config.Get().OAuthClientID)
	values.Set("scope", githubAdminScope)
	u.RawQuery = values.Encode()
	return u.String()
}

// oauthURL returns a URL for redirecting to the GitHub login page.
func oauthURL() string {
	u, err := url.Parse(githubAuthURL)
	if err != nil {
		// the redirection URL must be a valid URL
		panic(err)
	}
	values := u.Query()
	values.Set("client_id", config.Get().OAuthClientID)
	u.RawQuery = values.Encode()
	return u.String()
}

// loginHandler redirects to GitHub's OAuth login page.
func loginHandler() http.Handler {
	return http.RedirectHandler(oauthURL(), http.StatusTemporaryRedirect)
}

// appHandler redirects to GitHub's application registration page.
func appHandler() http.Handler {
	return http.RedirectHandler(githubAppURL, http.StatusTemporaryRedirect)
}

// oauthHandler is the OAuth handler for the GitHub response to a login request.
func oauthHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	code := v.Get("code")
	errstr := v.Get("error")
	if len(errstr) > 0 {
		logAndRedirect(w, r, front, "Failed to obtain temporary OAuth code: "+errstr)
		return
	}

	resp, err := http.PostForm(githubTokenURL, url.Values{
		"client_id":     {config.Get().OAuthClientID},
		"client_secret": {config.Get().OAuthClientSecret},
		"code":          {code},
	})
	if err != nil {
		// failed to issue POST request
		logErrorAndRedirect(w, r, front, err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// failed to read response body
		logErrorAndRedirect(w, r, front, err)
		return
	}

	q, err := url.ParseQuery(string(data))
	if err != nil {
		// failed to parse query data
		logErrorAndRedirect(w, r, front, err)
		return
	}

	errstr = q.Get("error")
	if len(errstr) > 0 {
		logAndRedirect(w, r, front, "Failed to obtain access token: "+errstr)
		return
	}

	accessToken := q.Get(githubTokenSessionKey)
	// save access token for this session
	session, _ := cookieStore.Get(r, githubSessionName)
	session.Values[githubTokenSessionKey] = accessToken
	err = session.Save(r, w)
	if err != nil {
		logErrorAndRedirect(w, r, front, err)
		return
	}

	scope := q.Get("scope")
	if scope != "" {
		gh, found := ghConns[accessToken]
		if !found {
			gh, err = NewGitHubConn(accessToken, scope)
			if err != nil {
				// failed to connect to GitHub.
				logErrorAndRedirect(w, r, front, err)
				return
			}
			ghConns[accessToken] = gh
		}
	}
	user := ghConns[accessToken].user.Login
	logAndRedirect(w, r, home, fmt.Sprintf("Yeay! User log in: %v", *user))
}

// logoutHandler revokes the current login session and
// redirects the user to front page.
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// delete access token for this session
	session, _ := cookieStore.Get(r, githubSessionName)
	accessToken := session.Values[githubTokenSessionKey].(string)
	session.Values[githubTokenSessionKey] = ""
	err := session.Save(r, w)
	if err != nil {
		logErrorAndRedirect(w, r, front, err)
		return
	}
	user := ghConns[accessToken].user.Login
	logAndRedirect(w, r, front, fmt.Sprintf("Logout: %v", *user))
}
