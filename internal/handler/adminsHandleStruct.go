package handler

type reqWithEmail struct {
	Email string `json:"email" binding:"required"`
}

type requestAdmin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type resetAdminPassword struct {
	Email     string `json:"email" binding:"required"`
	EmailCode string `json:"emailCode" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
