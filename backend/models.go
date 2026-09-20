package main

type Option struct {
	ID    string `json:"id" bson:"id"`
	Text  string `json:"text" bson:"text"`
	Votes int    `json:"votes" bson:"votes"`
}

type Poll struct {
	ID       string   `json:"id" bson:"_id,omitempty"`
	Creator  string   `json:"creator" bson:"creator"`
	Question string   `json:"question" bson:"question"`
	Options  []Option `json:"options" bson:"options"`
	VotedBy  []string `json:"voted_by" bson:"voted_by,omitempty"`
}

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
