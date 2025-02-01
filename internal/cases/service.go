package cases

// Service содержит логику работы сервиса

type Service struct {
	storage Storage //проверить на nil
	client  Client
}

// NewService создает новый сервис
func NewService(storage Storage, client Client) *Service {
	//проверить на nil
	return &Service{storage: storage, client: client}
}
