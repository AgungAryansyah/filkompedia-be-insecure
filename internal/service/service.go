package service

import (
	"github.com/yogarn/filkompedia-be/internal/repository"
	"github.com/yogarn/filkompedia-be/pkg/bcrypt"
	"github.com/yogarn/filkompedia-be/pkg/jwt"
	"github.com/yogarn/filkompedia-be/pkg/midtrans"
	"github.com/yogarn/filkompedia-be/pkg/smtp"
	"github.com/yogarn/filkompedia-be/pkg/supabase"
)

type Service struct {
	UserService     IUserService
	AuthService     IAuthService
	BookService     IBookService
	CartService     ICartService
	CommentService  ICommentService
	CheckoutService ICheckoutService
	PaymentService  IPaymentService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.IBcrypt, jwt jwt.IJwt, smtp *smtp.SMTPClient, midtrans midtrans.IMidtrans, supabase supabase.ISupabase) *Service {
	return &Service{
		UserService:     NewUserService(repository.UserRepository, repository.CartRepository, repository.PaymentRepository, repository.AuthRepository, repository.CheckoutRepository, repository.CommentRepository, supabase),
		AuthService:     NewAuthService(repository.AuthRepository, repository.UserRepository, bcrypt, jwt, smtp),
		BookService:     NewBookService(repository.BookRepository, repository.CartRepository, repository.CommentRepository, supabase),
		CartService:     NewCartService(repository.CartRepository, repository.UserRepository, repository.BookRepository),
		CommentService:  NewCommentService(repository.CommentRepository, repository.UserRepository),
		CheckoutService: NewCheckoutService(repository.CheckoutRepository, repository.CartRepository, repository.BookRepository, repository.UserRepository),
		PaymentService:  NewPaymentService(repository.PaymentRepository, midtrans, repository.UserRepository, repository.BookRepository, repository.CheckoutRepository),
	}
}
