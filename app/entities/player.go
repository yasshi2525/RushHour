package entities

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/yasshi2525/RushHour/app/services/auth"
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
)

// AuthList is list of all AuthType
var AuthList []AuthType

// InitAuthList instanciate AuthList
func InitAuthList() {
	AuthList = []AuthType{
		Local,
		Twitter,
		Google,
	}
}

// Player represents user information
type Player struct {
	Base
	Persistence

	Level PlayerType `gorm:"not null" json:"lv"`
	// DisplayName is public attribute and shown to everyone. Owner can change it.
	DisplayName string `gorm:"not null" sql:"type:text" json:"-"`
	// Image is public attribute and shown to everyone. Owner can change it.
	Image string `gorm:"not null" sql:"type:text" json:"-"`
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
	// Token is hidden attribute and used for authentication in RushHour
	Token string `gorm:"not null;index" json:"-"`

	ReRouting bool `gorm:"-"              json:"-"`
	Color     int  `gorm:"not null"       json:"color"`

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
	h := float64((int(o.ID)%6)*60 + (int(o.ID)/6)*15)
	var r, g, b int
	h0 := int(math.Floor(h / 60))

	switch h0 {
	case 0:
		r, g, b = 0xFF, int(h/60*0xFF), 0x00
	case 1:
		r, g, b = int((h-120)/60*0xFF), 0xFF, 0x00
	case 2:
		r, g, b = 0x00, 0xFF, int((h-120)/60*0xFF)
	case 3:
		r, g, b = 0x00, int((240-h)/60*0xFF), 0xFF
	case 4:
		r, g, b = int((240-h)/60*0xFF), 0x00, 0xFF
	case 5:
		r, g, b = 0xFF, 0x00, int((360-h)/60*0xFF)
	}

	o.Color = r<<16 + g<<8 + b

	o.O = o
	o.OwnerID = o.ID
	o.Init(m)
	o.Marshal()
	m.Add(o)
	return o
}

// OAuthSignIn finds or create Player by auth and loginid, then refresh token value.
func (m *Model) OAuthSignIn(authType AuthType, info *auth.UserInfo) (*Player, error) {
	if !info.IsValid() {
		return nil, fmt.Errorf("token is empty")
	}
	loginhash := auth.Digest(info.LoginID)
	if o, found := m.Logins[authType][loginhash]; found {
		enc := info.Enc()
		o.OAuthToken = enc.OAuthToken
		o.OAuthSecret = enc.OAuthSecret
		o.Token = info.Token()
		m.Tokens[o.Token] = o
		return o, nil
	}
	o := m.NewPlayer()
	o.Level = Normal
	o.ImportInfo(authType, info)
	return o, nil
}

// SignOut deletes token value.
func (o *Player) SignOut() {
	delete(o.M.Tokens, o.Token)
	o.OAuthToken = ""
	o.OAuthSecret = ""
	o.Token = ""
}

// PasswordSignIn finds Player by loginid and password, then refresh token value.
// arg must be plain text
func (m *Model) PasswordSignIn(loginid string, password string) (*Player, error) {
	if o, found := m.Logins[Local][auth.Digest(loginid)]; found {
		if encPassword := auth.Digest(password); strings.Compare(o.Password, encPassword) == 0 {
			delete(m.Tokens, o.Token)
			o.Token = auth.Digest(fmt.Sprintf("%s%s", time.Now(), encPassword))
			m.Tokens[o.Token] = o
			return o, nil
		}
	}
	return nil, fmt.Errorf("invalid user name or password")
}

// PasswordSignUp creates Player with loginid and password, then register token value.
// arg must be plain text
func (m *Model) PasswordSignUp(loginid string, password string) (*Player, error) {
	loginhash := auth.Digest(loginid)
	if _, found := m.Logins[Local][loginhash]; found {
		return nil, fmt.Errorf("id is already exists")
	}
	o := m.NewPlayer()
	o.Level = Normal
	o.LoginID = auth.Encrypt(loginid)
	o.Password = auth.Digest(password)
	o.DisplayName = loginid
	o.Auth = Local
	o.Token = auth.Digest(fmt.Sprintf("%s%s", time.Now(), password))
	o.M.Logins[Local][loginhash] = o
	o.M.Tokens[o.Token] = o
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
func (o *Player) ImportInfo(authType AuthType, info *auth.UserInfo) {
	enc := info.Enc()
	delete(o.M.Tokens, o.Token)

	o.DisplayName = enc.DisplayName
	o.Image = enc.Image
	o.LoginID = enc.LoginID
	o.Auth = authType
	o.OAuthToken = enc.OAuthToken
	o.OAuthSecret = enc.OAuthSecret
	o.Token = info.Token()

	o.M.Logins[authType][auth.Digest(info.LoginID)] = o
	o.M.Tokens[o.Token] = o
}

// ExportInfo decrypts user information
func (o *Player) ExportInfo() *auth.UserInfo {
	return (&auth.UserInfo{
		DisplayName: o.DisplayName,
		Image:       o.Image,
		LoginID:     o.LoginID,
		OAuthToken:  o.OAuthToken,
		OAuthSecret: o.OAuthSecret,
		IsEnc:       true,
	}).Dec()
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
		o.ID, o.LoginID, o.Level, o.DisplayName)
}

// Short returns short description
func (o *Player) Short() string {
	return fmt.Sprintf("%s(%d)", o.LoginID, o.ID)
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
