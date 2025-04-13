package service

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/internal/repository"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/midtrans"
)

type IPaymentService interface {
	GetPayment(paymentId uuid.UUID) (*entity.Payment, error)
	CreatePayment(userId uuid.UUID, checkoutId uuid.UUID, totalPrice float64) (*snap.Response, error)
	UpdatePaymentStatus(PaymentDetails map[string]any) error
	CheckUserBookPurchase(userId uuid.UUID, bookId uuid.UUID) (*bool, error)
	GetPayments(req model.PaymentReq) ([]entity.Payment, error)
	GetPaymentByCheckout(checkoutId uuid.UUID) (*entity.Payment, error)
	GetPaymentByUser(userId uuid.UUID) (*[]entity.Payment, error)
}

type PaymentService struct {
	paymentRepo repository.IPaymentRepository
	userRepo    repository.IUserRepository
	bookRepo    repository.IBookRepository
	midtrans    midtrans.IMidtrans
	chekoutRepo repository.ICheckoutRepository
}

func NewPaymentService(paymentRepo repository.IPaymentRepository, midtrans midtrans.IMidtrans, userRepo repository.IUserRepository, bookRepo repository.IBookRepository, chekoutRepo repository.ICheckoutRepository) IPaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		midtrans:    midtrans,
		userRepo:    userRepo,
		bookRepo:    bookRepo,
		chekoutRepo: chekoutRepo,
	}
}

func (s *PaymentService) GetPayment(paymentId uuid.UUID) (*entity.Payment, error) {
	return s.paymentRepo.GetPayment(paymentId)
}

func (s *PaymentService) CreatePayment(userId uuid.UUID, checkoutId uuid.UUID, totalPrice float64) (*snap.Response, error) {
	paymentId := uuid.New()

	var user entity.User
	if err := s.userRepo.GetUser(&user, userId); err != nil {
		return nil, err
	}

	var snapRes *snap.Response
	snapRes, err := s.midtrans.NewTransactionToken(paymentId.String(), int64(totalPrice), &user)
	if err != nil {
		return nil, err
	}

	if snapRes == nil {
		return nil, errors.New("nil response")
	}

	token, err := uuid.Parse(snapRes.Token)
	if err != nil {
		return nil, err
	}

	if err := s.paymentRepo.CreatePayment(entity.Payment{
		Id:         paymentId,
		Token:      token,
		UserId:     userId,
		CheckoutId: checkoutId,
		TotalPrice: totalPrice,
		StatusId:   0,
		CreatedAt:  time.Now(),
	}); err != nil {
		return nil, err
	}

	return snapRes, nil
}

func (s *PaymentService) UpdatePaymentStatus(PaymentDetails map[string]any) error {
	paymentIDs, ok := PaymentDetails["order_id"].(string)
	if !ok {
		return errors.New("invalid payment details")
	}

	paymentId, err := uuid.Parse(paymentIDs)
	if err != nil {
		return err
	}

	if payment, err := s.paymentRepo.GetPayment(paymentId); err != nil || payment == nil {
		return err
	}

	statusCode, ok := PaymentDetails["status_code"].(string)
	if !ok {
		return errors.New("invalid payment details")
	}

	grossAmount, ok := PaymentDetails["gross_amount"].(string)
	if !ok {
		return errors.New("invalid payment details")
	}

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	signatureKey, ok := PaymentDetails["signature_key"].(string)
	if !ok {
		return errors.New("invalid payment details")
	}

	hash := sha512.New()
	hash.Write([]byte(paymentIDs + statusCode + grossAmount + serverKey))
	verify := hex.EncodeToString(hash.Sum(nil))

	if signatureKey != verify {
		return errors.New("invalid payment details")
	}

	//todo improve this
	status, ok := PaymentDetails["transaction_status"]
	if !ok {
		return errors.New("invalid payment details")
	}

	fraud, ok := PaymentDetails["fraud_status"]
	if !ok {
		return errors.New("invalid payment details")
	}

	if status == "capture" {
		if fraud == "challenge" {
			if err := s.paymentRepo.UpdatePaymentStatus(4, paymentId); err != nil {
				return err
			}
		} else if fraud == "accept" {
			if err := s.paymentRepo.UpdatePaymentStatus(1, paymentId); err != nil {
				return err
			}
		}
	} else if status == "settlement" {
		if err := s.paymentRepo.UpdatePaymentStatus(5, paymentId); err != nil {
			return err
		}
	} else if status == "deny" {
		if err := s.paymentRepo.UpdatePaymentStatus(2, paymentId); err != nil {
			return err
		}
	} else if status == "cancel" || status == "expire" {
		if err := s.paymentRepo.UpdatePaymentStatus(3, paymentId); err != nil {
			return err
		}
	}

	return nil
}

func (s *PaymentService) CheckUserBookPurchase(userId uuid.UUID, bookId uuid.UUID) (*bool, error) {
	var book entity.Book
	if err := s.bookRepo.GetBook(&book, bookId); err != nil {
		return nil, err
	}

	var user entity.User
	if err := s.userRepo.GetUser(&user, userId); err != nil {
		return nil, err
	}

	return s.paymentRepo.CheckUserBookPurchase(userId, bookId)
}

func (s *PaymentService) GetPayments(req model.PaymentReq) ([]entity.Payment, error) {
	return s.paymentRepo.GetPayments(req.Page, req.PageSize)
}

func (s *PaymentService) GetPaymentByCheckout(checkoutId uuid.UUID) (*entity.Payment, error) {
	_, err := s.chekoutRepo.GetCheckout(checkoutId)
	if err != nil {
		return nil, err
	}

	return s.paymentRepo.GetPaymentByCheckout(checkoutId)
}

func (s *PaymentService) GetPaymentByUser(userId uuid.UUID) (*[]entity.Payment, error) {
	var user entity.User
	if err := s.userRepo.GetUser(&user, userId); err != nil {
		return nil, err
	}

	return s.paymentRepo.GetPaymentByUser(userId)
}
