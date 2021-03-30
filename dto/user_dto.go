package dto

import "com.jumbo/ginessential/model"

type UserDto struct {
	UserId uint `json:"id"`
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.User) UserDto  {
	return UserDto{
		UserId: user.ID,
		Name:      user.Name,
		Telephone: user.Telephone,
	}
}