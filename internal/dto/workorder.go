package dto

type AssignMechanic struct {
	MechanicID string `json:"mechanic_id" binding:"required,uuid"`
}
