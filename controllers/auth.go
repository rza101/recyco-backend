package controllers


import{
	"recyco/model"
	"recyco/middleware"
	"recyco/utils"

	"strings"
}


type userInput struct{
	Phone_number string `json:"phone_number" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role string `json:"role" binding:"required"`
}

type userUpdate struct{
	Phone_number string `json:"phone_number" binding:"required"`
	password string `json:"password" binding:"required"`
}


func CreatedUser (u *user) b {
	
}
