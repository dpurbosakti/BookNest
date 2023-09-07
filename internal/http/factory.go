package http

import (
	eh "book-nest/utils/emailhelper"

	"gorm.io/gorm"

	// auth
	authHdl "book-nest/internal/handlers/auth"
	authRepo "book-nest/internal/repositories/auth"
	authSrv "book-nest/internal/services/auth"

	// users
	userHdl "book-nest/internal/handlers/user"
	userRepo "book-nest/internal/repositories/user"
	userSrv "book-nest/internal/services/user"

	// books
	bookHdl "book-nest/internal/handlers/book"
	bookRepo "book-nest/internal/repositories/book"
	bookSrv "book-nest/internal/services/book"
)

type Presenter struct {
	Auth *authHdl.AuthHandler
	User *userHdl.UserHandler
	Book *bookHdl.BookHandler
}

func NewPresenter(db *gorm.DB) *Presenter {
	// auth
	ar := authRepo.NewAuthRepository()
	as := authSrv.NewAuthService(ar, db)
	ah := authHdl.NewAuthHandler(as)

	// emailhelper
	email := eh.NewEmailHelper()

	// users
	ur := userRepo.NewUserRepository()
	us := userSrv.NewUserService(ur, db, email)
	uh := userHdl.NewUserHandler(us)

	// books
	br := bookRepo.NewBookRepository()
	bs := bookSrv.NewBookService(br, db)
	bh := bookHdl.NewBookHandler(bs)

	return &Presenter{
		Auth: ah.(*authHdl.AuthHandler),
		User: uh.(*userHdl.UserHandler),
		Book: bh.(*bookHdl.BookHandler),
	}
}
