package vehicle

import (
	"strings"
	"time"
	"workshop-management/internal/domain/vehicle"
	"workshop-management/internal/dto"
	"workshop-management/utils"
)

type ServiceVehicle struct {
	VehicleRepo vehicle.Repository
}

func NewVehicleService(vehicleRepo vehicle.Repository) *ServiceVehicle {
	return &ServiceVehicle{
		VehicleRepo: vehicleRepo,
	}
}

func (s *ServiceVehicle) CreateVehicle(userId string, req dto.AddVehicle) (vehicle.Vehicle, error) {
	data := vehicle.Vehicle{
		Id:           utils.CreateUUID(),
		UserId:       userId,
		Brand:        utils.TitleCase(req.Brand),
		Model:        utils.TitleCase(req.Model),
		Year:         req.Year,
		Color:        utils.TitleCase(req.Color),
		LicensePlate: strings.ToUpper(req.LicensePlate),
		CreatedAt:    time.Now(),
	}

	if err := s.VehicleRepo.Store(data); err != nil {
		return vehicle.Vehicle{}, err
	}

	return data, nil
}

func (s *ServiceVehicle) GetVehicle(id string) (vehicle.Vehicle, error) {
	return s.VehicleRepo.GetVehicle(id)
}

func (s *ServiceVehicle) FetchVehicles(page, limit int, orderBy, orderDir, search, userId string) ([]vehicle.Vehicle, int64, error) {
	return s.VehicleRepo.FetchVehicles(page, limit, orderBy, orderDir, search, userId)
}
