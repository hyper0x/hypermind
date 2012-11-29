package utils

import (
	"encoding/json"
)

func UnmarshalUser(literals string) (User, error) {
	var user User
	err := json.Unmarshal([]byte(literals), &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func MarshalUser(user User) (string, error) {
	b, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	return string(b), nil
}



