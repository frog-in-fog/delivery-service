package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/frog-in-fog/delivery-system/gateway-service/models/dto"
)

var loginTemplate = template.Must(template.New("login").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Login Page</title>
</head>
<body>
	<h2>Login</h2>
	<form action="/login" method="post">
		<label for="username">Username:</label>
		<input type="text" id="username" name="username" required><br>
		<label for="password">Password:</label>
		<input type="password" id="password" name="password" required><br>
		<input type="submit" value="Login">
	</form>
</body>
</html>
`))

const (
	username = "demo"
	password = "password"
)

// func LoginPage(w http.ResponseWriter, r *http.Request) {
// 	loginTemplate.Execute(w, nil)
// }

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var result dto.OneLineResp
	req, err := http.NewRequest(http.MethodGet, "http://host.docker.internal:4000/api/v0/tokens", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Write(body)
	// if user == username && pass == password {
	// 	http.SetCookie(w, &http.Cookie{
	// 		Name:  "auth",
	// 		Value: "true",
	// 	})
	// 	http.Redirect(w, r, "/logger", http.StatusSeeOther)
	// } else {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprintln(w, "Invalid credentials")
	// }
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// cookie, err := r.Cookie("auth")
		// if err != nil || cookie.Value != "true" {
		// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
		// 	return
		// }

		// next(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func Proxy(path string, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL := target + r.URL.Path
		req, err := http.NewRequest(r.Method, targetURL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		var response dto.OneLineResp

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		if err := json.Unmarshal(body, &response); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		w.Write([]byte(response.Data))
	}
}
