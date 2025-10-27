package workorder

type RepoWorkOrder interface {
	Create(workOrder WorkOrder, svcWorkOrders []SvcWorkOrder) error
	GetById(id string) (WorkOrder, error)
	Update(workOrder WorkOrder, data map[string]interface{}) (int64, error)
}
