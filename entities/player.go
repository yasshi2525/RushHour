package entities

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/yasshi2525/RushHour/auth"
)

// PlayerType represents authenticate level
type PlayerType uint

// PlayerType represents authenticate level
const (
	Admin PlayerType = iota + 1
	Normal
	Guest
)

// AuthType represents which SNS account player sigin in
type AuthType uint

// AuthType represents which SNS account player sigin in
const (
	Local AuthType = iota + 1
	Twitter
	Google
	GitHub
)

// MarshalJSON converts AuthType to string
func (a AuthType) MarshalJSON() ([]byte, error) {
	switch a {
	case Local:
		return json.Marshal("RushHour")
	case Twitter:
		return json.Marshal("Twitter")
	case Google:
		return json.Marshal("Google")
	case GitHub:
		return json.Marshal("GitHub")
	}
	return json.Marshal("Unknown Service")
}

// AuthList is list of all AuthType
var AuthList []AuthType

// InitAuthList instanciate AuthList
func InitAuthList() {
	AuthList = []AuthType{
		Local,
		Twitter,
		Google,
		GitHub,
	}
}

// Player represents user information
type Player struct {
	Base
	Persistence

	Level PlayerType `gorm:"not null" json:"-"`
	// OAuthDisplayName is public attribute and shown to everyone. Owner can change it.
	OAuthDisplayName     string `gorm:"not null" sql:"type:text" json:"-"`
	CustomDisplayName    string `gorm:"not null" sql:"type:text" json:"-"`
	UseCustomDisplayName bool   `gorm:"not null" json:"-"`

	// OAuthImage is public attribute and shown to everyone. Owner can change it.
	OAuthImage     string `gorm:"not null" sql:"type:text" json:"-"`
	CustomImage    string `gorm:"not null" sql:"type:text" json:"-"`
	UseCustomImage bool   `gorm:"not null" json:"-"`

	// LoginID is hidden attribute and used for identification in OAuth App.
	LoginID string `gorm:"not null" sql:"type:text" json:"-"`
	// Password is hidden attribute, but owner can change it.
	// Password is empty in OAuth authentication.
	Password string   `gorm:"not null" sql:"type:text" json:"-"`
	Auth     AuthType `gorm:"not null;index" json:"-"`
	// OAuthToken is hidden attribute and used for OAuth authentication (access token).
	OAuthToken string `gorm:"not null" sql:"type:text" json:"-"`
	// OAuthSecret is hidden attribute and used for OAuth authentication (access token secret).
	OAuthSecret string `gorm:"not null" sql:"type:text" json:"-"`

	ReRouting bool `gorm:"-" json:"-"`

	// Hue is hue attribute on HSV model.
	Hue int `gorm:"not null" json:"hue"`

	RailNodes map[uint]*RailNode `gorm:"-" json:"-"`
	RailEdges map[uint]*RailEdge `gorm:"-" json:"-"`
	Stations  map[uint]*Station  `gorm:"-" json:"-"`
	Gates     map[uint]*Gate     `gorm:"-" json:"-"`
	Platforms map[uint]*Platform `gorm:"-" json:"-"`
	RailLines map[uint]*RailLine `gorm:"-" json:"-"`
	LineTasks map[uint]*LineTask `gorm:"-" json:"-"`
	Trains    map[uint]*Train    `gorm:"-" json:"-"`
}

// NewPlayer creates instance.
func (m *Model) NewPlayer() *Player {
	o := &Player{
		Base:        m.NewBase(PLAYER),
		Persistence: NewPersistence(),
	}

	o.O = o
	o.OwnerID = o.ID
	o.Init(m)
	o.Marshal()
	m.Add(o)
	return o
}

// OAuthSignIn finds or create Player by auth and loginid, then refresh token value.
func (m *Model) OAuthSignIn(authType AuthType, info *auth.OAuthInfo) (*Player, error) {
	if !info.IsValid() {
		return nil, fmt.Errorf("token is empty")
	}
	loginhash := m.auther.Digest(info.LoginID)
	if o, found := m.Logins[authType][loginhash]; found {
		enc, err := info.Enc()
		if err != nil {
			return nil, err
		}
		o.OAuthToken = enc.OAuthToken
		o.OAuthSecret = enc.OAuthSecret
		return o, nil
	}
	o := m.NewPlayer()
	o.Level = Normal
	o.Hue = rand.Intn(360)
	o.ImportInfo(authType, info)
	return o, nil
}

// SignOut deletes token value.
func (o *Player) SignOut() {
	o.OAuthToken = ""
	o.OAuthSecret = ""
}

// PasswordSignIn finds Player by loginid and password, then refresh token value.
// arg must be plain text
func (m *Model) PasswordSignIn(loginid string, password string) (*Player, error) {
	if o, found := m.Logins[Local][m.auther.Digest(loginid)]; found {
		if encPassword := m.auther.Digest(password); o.Password == encPassword {
			return o, nil
		}
	}
	return nil, fmt.Errorf("invalid user name or password")
}

// PasswordSignUp creates Player with loginid and password, then register token value.
// arg must be plain text
func (m *Model) PasswordSignUp(loginid string, password string, lv PlayerType) (*Player, error) {
	loginhash := m.auther.Digest(loginid)
	if _, found := m.Logins[Local][loginhash]; found {
		return nil, fmt.Errorf("id is already exists")
	}
	o := m.NewPlayer()
	o.Level = lv
	o.LoginID = m.auther.Encrypt(loginid)
	o.Password = m.auther.Digest(password)
	o.Auth = Local
	o.M.Logins[Local][loginhash] = o
	return o, nil
}

// B returns base information of this elements.
func (o *Player) B() *Base {
	return &o.Base
}

// P returns time information for database.
func (o *Player) P() *Persistence {
	return &o.Persistence
}

// ClearTracks eraces track infomation.
func (o *Player) ClearTracks() {
	for _, rn := range o.RailNodes {
		rn.Tracks = make(map[uint]map[uint]bool)
	}
}

// Init do nothing
func (o *Player) Init(m *Model) {
	o.Base.Init(PLAYER, m)
	o.RailNodes = make(map[uint]*RailNode)
	o.RailEdges = make(map[uint]*RailEdge)
	o.Stations = make(map[uint]*Station)
	o.Gates = make(map[uint]*Gate)
	o.Platforms = make(map[uint]*Platform)
	o.RailLines = make(map[uint]*RailLine)
	o.LineTasks = make(map[uint]*LineTask)
	o.Trains = make(map[uint]*Train)
}

// ImportInfo encrypts user information
func (o *Player) ImportInfo(authType AuthType, info *auth.OAuthInfo) error {
	enc, err := info.Enc()
	if err != nil {
		return err
	}

	o.OAuthDisplayName = enc.DisplayName
	o.OAuthImage = enc.Image
	o.LoginID = enc.LoginID
	o.Auth = authType
	o.OAuthToken = enc.OAuthToken
	o.OAuthSecret = enc.OAuthSecret

	o.M.Logins[authType][o.M.auther.Digest(info.LoginID)] = o
	return nil
}

// ExportInfo decrypts user information for OAuth
func (o *Player) ExportInfo() (*auth.OAuthInfo, error) {
	return (&auth.OAuthInfo{
		Handler:     o.M.auther,
		DisplayName: o.GetDisplayName(),
		Image:       o.GetImage(),
		LoginID:     o.LoginID,
		OAuthToken:  o.OAuthToken,
		OAuthSecret: o.OAuthSecret,
		IsEnc:       true,
	}).Dec()
}

// ExportJWTInfo decrypts user information for json web token
func (o *Player) ExportJWTInfo() *auth.JWTInfo {
	return &auth.JWTInfo{
		ID:    o.ID,
		Name:  o.M.auther.Decrypt(o.GetDisplayName()),
		Image: o.M.auther.Decrypt(o.GetImage()),
		Admin: o.Level == Admin,
		Hue:   o.Hue,
	}
}

// GetDisplayName returns customized name if player do.
func (o *Player) GetDisplayName() string {
	if o.UseCustomDisplayName {
		return o.CustomDisplayName
	}
	return o.OAuthDisplayName
}

// GetImage returns customized name if player do.
func (o *Player) GetImage() string {
	if o.UseCustomImage {
		return o.CustomImage
	}
	return o.OAuthImage
}

// Resolve set reference.
func (o *Player) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *RailNode:
			o.RailNodes[obj.ID] = obj
		case *RailEdge:
			o.RailEdges[obj.ID] = obj
		case *Station:
			o.Stations[obj.ID] = obj
		case *Gate:
			o.Gates[obj.ID] = obj
		case *Platform:
			o.Platforms[obj.ID] = obj
		case *RailLine:
			o.RailLines[obj.ID] = obj
		case *LineTask:
			o.LineTasks[obj.ID] = obj
		case *Train:
			o.Trains[obj.ID] = obj
		default:
			panic(fmt.Errorf("invalid type %v %+v", obj, obj))
		}
	}
}

// UnResolve unregisters specified refernce.
func (o *Player) UnResolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *RailNode:
			delete(o.RailNodes, obj.ID)
		case *RailEdge:
			delete(o.RailEdges, obj.ID)
		case *Station:
			delete(o.Stations, obj.ID)
		case *Gate:
			delete(o.Gates, obj.ID)
		case *Platform:
			delete(o.Platforms, obj.ID)
		case *RailLine:
			delete(o.RailLines, obj.ID)
		case *LineTask:
			delete(o.LineTasks, obj.ID)
		case *Train:
			delete(o.Trains, obj.ID)
		default:
			panic(fmt.Errorf("invalid type %v %+v", obj, obj))
		}
	}
}

// Marshal do nothing for implementing Resolvable
func (o *Player) Marshal() {
	// do-nothing
}

// UnMarshal set reference from id.
func (o *Player) UnMarshal() {

}

// CheckDelete check remain relation.
func (o *Player) CheckDelete() error {
	return nil
}

// BeforeDelete deletes related reference
func (o *Player) BeforeDelete() {
}

// Delete removes this entity with related ones.
func (o *Player) Delete() {
	o.M.Delete(o)
}

// String represents status
func (o *Player) String() string {
	o.Marshal()
	return fmt.Sprintf("%s(%d):nm=%s,lv=%v:%s", o.Type().Short(),
		o.ID, o.LoginID, o.Level, o.OAuthDisplayName)
}

// Short returns short description
func (o *Player) Short() string {
	return fmt.Sprintf("%s(%d)", o.LoginID, o.ID)
}

// MarshalJSON returns plain text data.
func (o *Player) MarshalJSON() ([]byte, error) {
	type Alias Player
	return json.Marshal(&struct {
		DisplayName string `json:"name"`
		Image       string `json:"image"`
		Admin       bool   `json:"admin,omitempty"`
		*Alias
	}{
		DisplayName: o.M.auther.Decrypt(o.GetDisplayName()),
		Image:       o.M.auther.Decrypt(o.GetImage()),
		Alias:       (*Alias)(o),
		Admin:       o.Level == Admin,
	})
}

func (pt PlayerType) String() string {
	switch pt {
	case Admin:
		return "admin"
	case Normal:
		return "normal"
	case Guest:
		return "guest"
	}
	return "???"
}
