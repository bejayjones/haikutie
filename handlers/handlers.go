package handlers

import (
	"fmt"
	"haikutie/data"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	db        *data.Helper
	templates *template.Template
}

func New(database *data.Helper) *Handler {
	templates := template.Must(template.ParseGlob("templates/*.html"))
	return &Handler{
		db:        database,
		templates: templates,
	}
}

func (h *Handler) getCurrentUserID(r *http.Request) (int64, bool) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return 0, false
	}
	userID, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		return 0, false
	}
	return userID, true
}

func (h *Handler) getCurrentUser(r *http.Request) (*data.User, bool) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		return nil, false
	}

	user, err := h.db.GetUserByIDHelper(userID)
	if err != nil {
		return nil, false
	}
	return &user, true
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	_, ok := h.getCurrentUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/library", http.StatusSeeOther)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.templates.ExecuteTemplate(w, "login.html", nil)
		return
	}

	// POST
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.db.GetUserHelper(username)
	if err != nil || user.Password != password {
		h.templates.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"Error": "Invalid credentials",
		})
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.FormatInt(user.ID, 10),
		Path:  "/",
	})

	http.Redirect(w, r, "/library", http.StatusSeeOther)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.templates.ExecuteTemplate(w, "register.html", nil)
		return
	}

	// POST
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := h.db.CreateUserHelper(username, password)
	if err != nil {
		h.templates.ExecuteTemplate(w, "register.html", map[string]interface{}{
			"Error": "Username already exists",
		})
		return
	}

	// Auto-login after registration
	user, _ := h.db.GetUserHelper(username)
	http.SetCookie(w, &http.Cookie{
		Name:  "user_id",
		Value: strconv.FormatInt(user.ID, 10),
		Path:  "/",
	})

	http.Redirect(w, r, "/library", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) Library(w http.ResponseWriter, r *http.Request) {
	user, ok := h.getCurrentUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	haikus, err := h.db.GetReceivedHaikusHelper(user.ID)
	if err != nil {
		log.Println("Error fetching haikus:", err)
		haikus = []data.GetReceivedHaikusRow{}
	}

	h.templates.ExecuteTemplate(w, "library.html", map[string]interface{}{
		"User":   user,
		"Haikus": haikus,
	})
}

func (h *Handler) Compose(w http.ResponseWriter, r *http.Request) {
	user, ok := h.getCurrentUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	users, err := h.db.GetAllUsersHelper()
	if err != nil {
		log.Println("Error fetching users:", err)
		users = []data.GetAllUsersRow{}
	}

	// Filter out current user
	var otherUsers []data.GetAllUsersRow
	for _, u := range users {
		if u.ID != user.ID {
			otherUsers = append(otherUsers, u)
		}
	}

	h.templates.ExecuteTemplate(w, "compose.html", map[string]interface{}{
		"User":  user,
		"Users": otherUsers,
	})
}

func (h *Handler) SendHaiku(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/compose", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form:", err)
	}
	fmt.Println(r.FormValue("line2"))
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	toUserID, _ := strconv.ParseInt(r.FormValue("to_user_id"), 10, 64)
	line1 := r.FormValue("line1")
	line2 := r.FormValue("line2")
	line3 := r.FormValue("line3")

	err = h.db.CreateHaikuHelper(userID, toUserID, line1, line2, line3)
	if err != nil {
		log.Println("Error creating haiku:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/library", http.StatusSeeOther)
}

func (h *Handler) ReceivedHaikus(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getCurrentUserID(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	haikus, err := h.db.GetReceivedHaikusHelper(userID)
	if err != nil {
		log.Println("Error fetching haikus:", err)
		haikus = []data.GetReceivedHaikusRow{}
	}

	h.templates.ExecuteTemplate(w, "haikus-list.html", map[string]interface{}{
		"Haikus": haikus,
	})
}
