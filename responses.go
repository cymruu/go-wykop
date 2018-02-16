package wykop

import (
	"fmt"
	"strings"
	"time"
)

const WYKOP_TRUE_RESPONSE = "[true]"

type WykopTime struct {
	time.Time
}

const (
	errMalformedReponse uint16 = 0xFFFF - iota
)

func (wt *WykopTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	wt.Time, err = time.Parse("2006-01-02 15:04:05", s) //IMO its's really retearded idea to use actual date values as placeholders
	if err != nil {
		wt.Time, _ = time.Parse(time.RFC3339Nano, s)
	}
	return
}

type ErrorResponse struct {
	ErrorObject struct {
		Message string `json:"message"`
		Code    uint16 `json:"code"`
	} `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s [%d]", e.ErrorObject.Message, e.ErrorObject.Code)
}

type userGroup uint16

func (userGroup userGroup) String() string {
	return userGroups[uint16(userGroup)]
}

var userGroups = map[uint16]string{
	0:    "zielony",
	1:    "pomarańczowy",
	2:    "bordowy",
	5:    "administraator",
	1001: "zbanowany",
	1002: "usunięty",
	2001: "klient",
}

type AuthorizationResponse struct {
	Login   string `json:"login"`
	Email   string `json:"email"`
	Userkey string `json:"userkey"`
}
type EntryResponse struct {
	ID           int             `json:"id"`
	Author       string          `json:"author"`
	AuthorAvatar string          `json:"author_avatar"`
	AuthorGroup  userGroup       `json:"author_group"`
	Date         WykopTime       `json:"date"`
	Body         string          `json:"body"`
	URL          string          `json:"url"`
	Comments     []EntryResponse `json:"comments"`
	VoteCount    uint32          `json:"vote_count"`
	Voters       []Voters        `json:"voters"`
	Embed        Embed           `json:"embed"`
}

type Voters struct {
	Author       string    `json:"author"`
	AuthorAvatar string    `json:"author_avatar"`
	AuthorGroup  userGroup `json:"author_group"`
}
type Embed struct {
	Type    string `json:"id"`
	Preview string `json:"preview"`
	URL     string `json:"url"`
	Source  string `json:"source"`
	Plus18  bool   `json:"plus18"`
}
type UserResponse struct {
	Login       string    `json:"login"`
	Email       string    `json:"email"`
	PEmail      string    `json:"public_email"`
	Name        string    `json:"name"`
	WWW         string    `json:"www"`
	About       string    `json:"about"`
	AuthorGroup userGroup `json:"author_group"`
	Avatar      string    `json:"avatar"`
	Sex         string    `json:"sex"`
}
type ConversationListItem struct {
	LastUpdate  WykopTime `json:"last_update"`
	Author      string    `json:"author"`
	Avatar      string    `json:"author_avatar"`
	AuthorGroup userGroup `json:"author_group"`
	Status      string    `json:"status"`
}
type Notification struct {
	Author       string    `json:"author"`
	AuthorAvatar string    `json:"author_avatar"`
	AuthorGroup  userGroup `json:"author_group"`
	Date         WykopTime `json:"date"`
	Body         string    `json:"body"`
	Type         string    `json:"type"`
}
