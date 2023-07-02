package templater

type Meta struct {
	Subject string `json:"subject" yaml:"subject"`
}

type Content struct {
	Body string
	Meta Meta
}
