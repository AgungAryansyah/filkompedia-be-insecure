package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/internal/repository"
	"github.com/yogarn/filkompedia-be/model"
)

type ICheckoutService interface {
	GetUserCheckouts(userId uuid.UUID) (*[]entity.Checkout, error)
	GetCheckoutCarts(checkoutId uuid.UUID) (*[]entity.Cart, error)
	Checkout(checkoutReq model.CheckoutRequest, userId uuid.UUID, checkoutId uuid.UUID) (float64, error)
}

type CheckoutService struct {
	checkoutRepo repository.ICheckoutRepository
	cartRepo     repository.ICartRepository
	bookRepo     repository.IBookRepository
	userRepo     repository.IUserRepository
}

func NewCheckoutService(checkoutRepo repository.ICheckoutRepository, cartRepo repository.ICartRepository, bookRepo repository.IBookRepository, userRepo repository.IUserRepository) ICheckoutService {
	return &CheckoutService{
		checkoutRepo: checkoutRepo,
		cartRepo:     cartRepo,
		bookRepo:     bookRepo,
		userRepo:     userRepo,
	}
}

func (s *CheckoutService) GetUserCheckouts(userId uuid.UUID) (*[]entity.Checkout, error) {
	var user entity.User
	if err := s.userRepo.GetUser(&user, userId); err != nil {
		return nil, err
	}

	return s.checkoutRepo.GetUserCheckouts(userId)
}

func (s *CheckoutService) GetCheckoutCarts(checkoutId uuid.UUID) (*[]entity.Cart, error) {
	_, err := s.checkoutRepo.GetCheckout(checkoutId)
	if err != nil {
		return nil, err
	}

	return s.checkoutRepo.GetCheckoutCarts(checkoutId)
}

func (s *CheckoutService) Checkout(checkoutReq model.CheckoutRequest, userId uuid.UUID, checkoutId uuid.UUID) (totalPrice float64, err error) {
	var user entity.User
	if err := s.userRepo.GetUser(&user, userId); err != nil {
		return 0, err
	}

	for _, cart_id := range checkoutReq.CartsId {
		var cart entity.Cart
		if err := s.cartRepo.GetCart(&cart, cart_id); err != nil {
			return 0, err
		}

		if cart.CheckoutId != uuid.Nil {
			return 0, errors.New("invalid input at " + cart_id.String() + " where it's already being checked out")
		}

		var book entity.Book
		if err := s.bookRepo.GetBook(&book, cart.BookId); err != nil {
			return 0, err
		}

		if cart.UserId != userId {
			return 0, errors.New("invalid input at " + cart_id.String() + " where it's not belong to the user")
		}

		totalPrice += (float64(cart.Amount) * book.Price)
	}

	if err := s.checkoutRepo.NewCheckout(checkoutId, userId); err != nil {
		return 0, err
	}

	for _, cart_id := range checkoutReq.CartsId {
		if err := s.checkoutRepo.AddCheckoutId(cart_id, checkoutId); err != nil {
			return 0, err
		}
	}

	return totalPrice, nil
}
