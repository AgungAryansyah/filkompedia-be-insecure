package service

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/internal/repository"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/response"
	"github.com/yogarn/filkompedia-be/pkg/supabase"
)

type IBookService interface {
	GetBook(bookId uuid.UUID) (*model.BookResponse, error)
	SearchBooks(bookSearch model.BookSearch) (*[]model.BookResponse, error)
	CreateBook(create *model.CreateBook) error
	DeleteBook(bookId uuid.UUID) error
	EditBook(edit model.EditBook) error
	UploadBookCover(file *multipart.FileHeader) (string, error)
}

type BookService struct {
	bookRepo    repository.IBookRepository
	cartRepo    repository.ICartRepository
	commentRepo repository.ICommentRepository
	Supabase    supabase.ISupabase
}

func NewBookService(bookRepo repository.IBookRepository, cartRepo repository.ICartRepository, commentRepo repository.ICommentRepository, Supabase supabase.ISupabase) IBookService {
	return &BookService{
		bookRepo:    bookRepo,
		cartRepo:    cartRepo,
		commentRepo: commentRepo,
		Supabase:    Supabase,
	}
}

func (s *BookService) GetBook(bookId uuid.UUID) (*model.BookResponse, error) {
	var book entity.Book
	err := s.bookRepo.GetBook(&book, bookId)
	if err != nil {
		return nil, err
	}

	bookResponse := model.BookToBookResponse(book)

	return &bookResponse, nil
}

func (s *BookService) SearchBooks(bookSearch model.BookSearch) (*[]model.BookResponse, error) {
	var booksEntity []entity.Book
	err := s.bookRepo.SearchBooks(&booksEntity, bookSearch.Page, bookSearch.PageSize, bookSearch.SearchParam)
	if err != nil {
		return nil, err
	}

	booksResponse := make([]model.BookResponse, len(booksEntity))
	for i, book := range booksEntity {
		booksResponse[i] = model.BookToBookResponse(book)
	}

	return &booksResponse, nil
}

func (s *BookService) CreateBook(create *model.CreateBook) error {
	return s.bookRepo.CreateBook(&entity.Book{
		Id:           uuid.New(),
		Title:        create.Title,
		Description:  create.Description,
		Author:       create.Author,
		ReleaseDate:  create.ReleaseDate,
		Price:        create.Price,
		Introduction: create.Introduction,
		Image:        create.Image,
		File:         create.File,
	})
}

func (s *BookService) DeleteBook(bookId uuid.UUID) error {
	var book entity.Book
	if err := s.bookRepo.GetBook(&book, bookId); err != nil {
		return err
	}

	if err := s.commentRepo.DeleteCommentByBook(bookId); err != nil {
		return err
	}

	err := s.cartRepo.DeleteCartByBook(bookId)
	if err != nil && err != &response.CartNotFound {
		return err
	}

	if err := s.bookRepo.DeleteBook(bookId); err != nil {
		return err
	}

	return nil
}

func (s *BookService) EditBook(edit model.EditBook) error {
	var book entity.Book
	if err := s.bookRepo.GetBook(&book, edit.Id); err != nil {
		return err
	}

	//todo improve this
	if edit.Title == "" {
		edit.Title = book.Title
	}
	if edit.Description == "" {
		edit.Description = book.Description
	}
	if edit.Introduction == "" {
		edit.Introduction = book.Introduction
	}
	if edit.Image == "" {
		edit.Image = book.Image
	}
	if edit.Author == "" {
		edit.Author = book.Author
	}
	if edit.ReleaseDate == "" {
		edit.ReleaseDate = book.ReleaseDate
	}
	if edit.Price == 0 {
		edit.Price = book.Price
	}

	if err := s.bookRepo.EditBook(&edit); err != nil {
		return err
	}

	return nil
}

func (s *BookService) UploadBookCover(file *multipart.FileHeader) (string, error) {
	url, err := s.Supabase.UploadFile(file, "cover")
	return url, err
}
