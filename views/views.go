package views

type ViewObj struct {
	Title  string
	Errors map[string]string
}

func NewViewObj(title string) *ViewObj {
	return &ViewObj{
		Title:  title,
		Errors: make(map[string]string),
	}
}
