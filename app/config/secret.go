package config

// CnfAdmin is configuration about initial admin user authorization
type CnfAdmin struct {
	UserName string
	Password string
}

// CnfDatabase is configuration about connectiing user data stored database
type CnfDatabase struct {
	Driver string `validation:"oneof=mysql"`
	Spec   string
}

// CnfTwitter is configuration about Twitter OAuth
type CnfTwitter struct {
	Token        string
	Secret       string
	Request      string `validate:"url"`
	Authenticate string `validate:"url"`
	AccessToken  string `validate:"url"`
}

// CnfOAuth is general configuration about OAuth
type CnfOAuth struct {
	Client string
	Secret string
}

// CnfAuth is auth section of secret.conf
type CnfAuth struct {
	BaseURL string `validate:"url"`
	Salt    string
	Key     string `validate:"len=16"`
	State   string

	Twitter CnfTwitter
	Google  CnfOAuth
	GitHub  CnfOAuth
}

// CnfSecret is root section of secret.conf
type CnfSecret struct {
	Admin CnfAdmin
	Auth  CnfAuth
	DB    CnfDatabase
}
