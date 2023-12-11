package views

type SessionDataView struct {
	Authenticated bool
	UserID        uint
	FirstName     string
	LastName      string
	IsAdmin       bool
}

type ViewObj struct {
	Title   string
	Errors  map[string]string
	Session *SessionDataView
}

func NewViewObj(title string) *ViewObj {
	return &ViewObj{
		Title:  title,
		Errors: make(map[string]string),
		Session: &SessionDataView{
			Authenticated: false,
			UserID:        0,
			FirstName:     "",
			LastName:      "",
			IsAdmin:       false,
		},
	}
}

func (v *ViewObj) UpdateSession(session *SessionDataView) {
	v.Session = session
}
