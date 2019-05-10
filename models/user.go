package models

//type User struct {
//	ID           uint   `json:"-"`
//	Login        string `json:"login"`
//	Email        string `json:"email"`
//	Password     string `json:"password,omitempty"`
//	HashPassword string `json:"-"`
//	Score        uint   `json:"score"`
//	Picture      string `json:"avatarUrl"`
//	Games        uint   `json:"count"`
//}

var AuthManager AuthCheckerClient

type Session struct {
	Email     string `json:"email"`
	SessionID string `json:"sessionid"`
}

var Sessions map[string]string

func init() {
	Sessions = make(map[string]string)
}

func ReturnCountOfSessions() int {
	return len(Sessions)
}

type LeadersInfo1 struct {
	Count       uint `json:"count"`
	UsersOnPage uint `json:"usersOnPage"`
}
