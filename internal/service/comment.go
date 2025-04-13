package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/internal/repository"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

type ICommentService interface {
	GetComment(id uuid.UUID) (*model.CommentRes, error)
	GetCommentByBook(bookId uuid.UUID) (*[]model.CommentRes, error)
	CreateComment(commentReq *model.CreateComment, userId uuid.UUID) error
	UpdateComment(commentReq *model.UpdateComment, userId uuid.UUID, bookId uuid.UUID, commentId uuid.UUID) error
	DeleteComment(id uuid.UUID, userId uuid.UUID) error
}

type CommentService struct {
	commentRepository repository.ICommentRepository
	userRepository    repository.IUserRepository
}

func NewCommentService(commentRepository repository.ICommentRepository, userRepository repository.IUserRepository) ICommentService {
	return &CommentService{
		commentRepository: commentRepository,
		userRepository:    userRepository,
	}
}

func (s *CommentService) GetComment(id uuid.UUID) (*model.CommentRes, error) {
	comment, err := s.commentRepository.GetComment(id)
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = s.userRepository.GetUser(&user, comment.UserId)
	if err != nil {
		return nil, err
	}

	return &model.CommentRes{
		Id:             comment.Id,
		UserId:         comment.UserId,
		Username:       user.Username,
		ProfilePicture: user.ProfilePicture,
		BookId:         comment.BookId,
		Comment:        comment.Comment,
		Rating:         comment.Rating,
		CreatedAt:      comment.CreatedAt,
	}, nil
}

func (s *CommentService) GetCommentByBook(bookId uuid.UUID) (*[]model.CommentRes, error) {
	comments, err := s.commentRepository.GetCommentByBook(bookId)
	if err != nil {
		return nil, err
	}

	var commentResponses []model.CommentRes
	for _, comment := range *comments {
		var user entity.User
		err := s.userRepository.GetUser(&user, comment.UserId)
		if err != nil {
			return nil, err
		}

		commentResponses = append(commentResponses, model.CommentRes{
			Id:             comment.Id,
			UserId:         comment.UserId,
			Username:       user.Username,
			ProfilePicture: user.ProfilePicture,
			BookId:         comment.BookId,
			Comment:        comment.Comment,
			Rating:         comment.Rating,
			CreatedAt:      comment.CreatedAt,
		})
	}

	return &commentResponses, nil
}

func (s *CommentService) CreateComment(commentReq *model.CreateComment, userId uuid.UUID) error {
	comment := &entity.Comment{
		Id:        uuid.New(),
		UserId:    userId,
		BookId:    commentReq.BookId,
		Comment:   commentReq.Comment,
		Rating:    commentReq.Rating,
		CreatedAt: time.Now(),
	}

	return s.commentRepository.CreateComment(comment)
}

func (s *CommentService) UpdateComment(commentReq *model.UpdateComment, userId uuid.UUID, bookId uuid.UUID, commentId uuid.UUID) error {
	comment := &entity.Comment{
		Id:        commentId,
		UserId:    userId,
		BookId:    bookId,
		Comment:   commentReq.Comment,
		Rating:    commentReq.Rating,
		CreatedAt: time.Now(),
	}

	return s.commentRepository.UpdateComment(comment)
}

func (s *CommentService) DeleteComment(id uuid.UUID, userId uuid.UUID) error {
	var user entity.User

	err := s.userRepository.GetUser(&user, userId)
	if err != nil {
		return err
	}

	comment, err := s.GetComment(id)
	if err != nil {
		return err
	}

	if comment.UserId != userId && user.RoleId != 1 {
		return &response.RoleUnauthorized
	}

	return s.commentRepository.DeleteComment(id)
}
