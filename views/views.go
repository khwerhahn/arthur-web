package views

import (
	"arthur-web/database"
	"arthur-web/globals"
	"arthur-web/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionDataView struct {
	IsAuthenticated   bool
	UserID            uint
	FirstName         string
	LastName          string
	Username          string
	IsAdmin           bool
	ProfilePictureUrl string
}

type ViewObj struct {
	Title   string
	Link    string
	Errors  map[string]string
	Session *SessionDataView
}

func NewViewObj(title string, link string) *ViewObj {
	return &ViewObj{
		Title:  title,
		Link:   link,
		Errors: make(map[string]string),
		Session: &SessionDataView{
			IsAuthenticated:   false,
			UserID:            0,
			FirstName:         "",
			LastName:          "",
			Username:          "",
			IsAdmin:           false,
			ProfilePictureUrl: "",
		},
	}
}

func (v *ViewObj) UpdateViewObjSession(c *gin.Context) (*ViewObj, error) {
	DB := database.DB
	session := sessions.Default(c)
	isAuthenticated := session.Get(globals.IsAuthenticated)
	isAdmin := session.Get(globals.IsAdmin)
	userID := session.Get(globals.UserID)
	if isAuthenticated == nil {
		isAuthenticated = false
	}
	if isAdmin == nil {
		isAdmin = false
	}
	if userID == nil {
		userID = 0
	}
	v.Session.IsAuthenticated = isAuthenticated.(bool)
	v.Session.IsAdmin = isAdmin.(bool)

	// get user data
	if userID != 0 {
		var user models.User
		DB.Where("id = ?", userID).First(&user)
		if user.Username == "" {
			return v, nil
		}
		v.Session.UserID = user.ID
		v.Session.Username = user.Username
		v.Session.FirstName = user.FirstName
		v.Session.LastName = user.LastName
		v.Session.ProfilePictureUrl = user.ProfileImageUrl
	} else {
		v.Session.UserID = 0
		v.Session.Username = ""
		v.Session.FirstName = ""
		v.Session.LastName = ""
		v.Session.ProfilePictureUrl = ""
	}
	return v, nil
}

func (v *ViewObj) AddError(key, value string) {
	v.Errors[key] = value
}
