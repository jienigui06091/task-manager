package model

type TaskPage struct {
	List []Task `json:"list"`
	PageDTO
}
