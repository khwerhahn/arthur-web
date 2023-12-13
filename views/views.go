package views

type SessionDataView struct {
	IsAuthenticated bool
	UserID          uint
	FirstName       string
	LastName        string
	Username        string
	IsAdmin         bool
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
			IsAuthenticated: false,
			UserID:          0,
			FirstName:       "",
			LastName:        "",
			Username:        "",
			IsAdmin:         false,
		},
	}
}

func (v *ViewObj) UpdateSession(session *SessionDataView) {
	v.Session = session
}

func (v *ViewObj) AddError(key, value string) {
	v.Errors[key] = value
}
