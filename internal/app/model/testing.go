package model

func TestUser() *User {
	return &User{}
}

func TestUsers() []*User {
	var users []*User = make([]*User, 10)

	u := &User{}

	u.DBID = 1
	u.UserId = 12345
	u.FirstName = "Username_1"
	u.Registered = true

	users = append(users, u)

	return users
}

func TestQuestion() *Question {
	return &Question{
		Question: "Как зовут мою любимку?",
		Answer:   "Алёнушка",
	}
}
