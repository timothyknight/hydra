package provider

type Provider interface {
	GetAuthenticationURL(state string) string
	Exchange(code string) (Session, error)
	GetID() string
}
