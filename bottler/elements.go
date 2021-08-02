package bottler

// bottler uses a collection data structures to
// programatically generate a Python3 "Bottle" application.
// These enscribed in JSON or XML. They are basted on
// HTML form element plus custom "Web Components" explcitily
// supported by bottler.

// ELement
type Element interface {
	ToHTML() string
	ToJSON() string
}

// Form presents an HTML form
type Form struct {
	Class   string    `json:"class,omitempty" xml:"class,omitempty"`
	Id      string    `json:"id,omitempty" xml:"id,omitempty"`
	Action  string    `json:"action,omitempty" xml:"action,omitempty"`
	Method  string    `json:"method,omitempty" xml:"method,omitempty"`
	Content []Element `json:"content,omitempty", xml:"content,omitempty"`
}

// Div groups elements together
type Div struct {
	Class   string    `json:"class,omitemtpy" xml:"class,omitempty"`
	Id      string    `json:"id,omitempty" xml:"id,omitempty"`
	Contain []Element `json:"content,omitempty", xml:"content,omitempty"`
}

// Label element
type Label struct {
	Class   string    `json:"class,omitemtpy" xml:"class,omitempty"`
	Id      string    `json:"id,omitempty" xml:"id,omitempty"`
	For     string    `json:"for,omitempty" xml:"for,omitempty"`
	Contain []Element `json:"content,omitempty", xml:"content,omitempty"`
}

// Input element
type Input struct {
	Class       string    `json:"class,omitemtpy" xml:"class,omitempty"`
	Id          string    `json:"id,omitempty" xml:"id,omitempty"`
	Type        string    `json:"type,omitempty" xml:"type,omitempty"`
	Name        string    `json:"name,omitempty" xml:"name,omitempty"`
	Value       string    `json:"value,omitempty" xml:"value,omitempty"`
	Title       string    `json:"title,omitempty" xml:"title,omitempty"`
	Placeholder string    `json:"placeholder,omitempty" xml:"placeholder,omitempty"`
	Content     []Element `json:"content,omitempty" xml:"content,omitempty"`
}
