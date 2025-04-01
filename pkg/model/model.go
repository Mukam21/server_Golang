package model

type Person struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         *int    `json:"age,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
}

type PersonRequest struct {
	Name       string  `json:"name" binding:"required"`
	Surname    string  `json:"surname" binding:"required"`
	Patronymic *string `json:"patronymic,omitempty"`
}
