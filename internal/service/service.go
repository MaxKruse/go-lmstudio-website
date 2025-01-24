package service

type Service[T any] interface {
	GetAll() ([]T, error)
	GetById(id int32) (T, error)
	Create(entity T) (T, error)
	Update(entity T) error
	Delete(id int32) error
}
