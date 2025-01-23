package converters

import (
	"github.com/maxkruse/go-lmstudio-website/internal/models/dtos"
	"github.com/maxkruse/go-lmstudio-website/internal/models/entities"
)

func BookDtoToEntity(dto dtos.Book) entities.Book {
	return entities.Book{
		Id:            dto.Id,
		Title:         dto.Title,
		Author:        dto.Author,
		Description:   dto.Description,
		ImageUrl:      dto.ImageUrl,
		PublishedDate: dto.PublishedDate,
		Isbn:          dto.Isbn,
		Price:         dto.Price,
	}
}

func BookEntityToDto(entity entities.Book) dtos.Book {
	return dtos.Book{
		Id:            entity.Id,
		Title:         entity.Title,
		Author:        entity.Author,
		Description:   entity.Description,
		ImageUrl:      entity.ImageUrl,
		PublishedDate: entity.PublishedDate,
		Isbn:          entity.Isbn,
		Price:         entity.Price,
	}
}

func BookEntityToDtoSlice(entities []entities.Book) []dtos.Book {
	dtos := make([]dtos.Book, len(entities))
	for i, entity := range entities {
		dtos[i] = BookEntityToDto(entity)
	}
	return dtos
}

func BookDtoToEntitySlice(dtos []dtos.Book) []entities.Book {
	entities := make([]entities.Book, len(dtos))
	for i, dtos := range dtos {
		entities[i] = BookDtoToEntity(dtos)
	}
	return entities
}
