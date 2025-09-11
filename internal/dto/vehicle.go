package dto

type AddVehicle struct {
	Brand        string `json:"brand" binding:"required"`
	Model        string `json:"model" binding:"required"`
	Year         string `json:"year" binding:"required"`
	LicensePlate string `json:"license_plate" binding:"required"`
	Color        string `json:"color" binding:"required"`
}
