package http

import (
	"book-nest/clients/biteship"
	"book-nest/clients/gomail"
	"book-nest/clients/midtrans"

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

	// orders
	orderHdl "book-nest/internal/handlers/order"
	orderRepo "book-nest/internal/repositories/order"
	orderSrv "book-nest/internal/services/order"

	// address
	addressHdl "book-nest/internal/handlers/address"
	addressRepo "book-nest/internal/repositories/address"
	addressSrv "book-nest/internal/services/address"

	// couriers
	courierHdl "book-nest/internal/handlers/courier"
	courierRepo "book-nest/internal/repositories/courier"
	courierSrv "book-nest/internal/services/courier"
)

type Presenter struct {
	Auth    *authHdl.AuthHandler
	User    *userHdl.UserHandler
	Book    *bookHdl.BookHandler
	Order   *orderHdl.OrderHandler
	Address *addressHdl.AddressHandler
	Courier *courierHdl.CourierHandler
}

func NewPresenter(db *gorm.DB) *Presenter {
	// clients
	gomail := gomail.NewGomailClient()
	midtrans := midtrans.NewMidtransClient()
	biteship := biteship.NewBiteshipClient()

	// users
	ur := userRepo.NewUserRepository()
	us := userSrv.NewUserService(ur, db, gomail)
	uh := userHdl.NewUserHandler(us)

	// auth
	ar := authRepo.NewAuthRepository()
	as := authSrv.NewAuthService(ar, ur, db)
	ah := authHdl.NewAuthHandler(as)

	// books
	br := bookRepo.NewBookRepository()
	bs := bookSrv.NewBookService(br, db)
	bh := bookHdl.NewBookHandler(bs)

	// orders
	or := orderRepo.NewOrderRepository()
	os := orderSrv.NewOrderService(or, br, ur, db, gomail, midtrans)
	oh := orderHdl.NewOrderHandler(os)

	// address
	adr := addressRepo.NewAddressRepository()
	ads := addressSrv.NewAddressService(adr, db)
	adh := addressHdl.NewAddressHandler(ads)

	// couriers
	cr := courierRepo.NewCourierRepository()
	cs := courierSrv.NewCourierService(cr, adr, or, br, db, biteship)
	ch := courierHdl.NewCourierHandler(cs)

	return &Presenter{
		Auth:    ah.(*authHdl.AuthHandler),
		User:    uh.(*userHdl.UserHandler),
		Book:    bh.(*bookHdl.BookHandler),
		Order:   oh.(*orderHdl.OrderHandler),
		Address: adh.(*addressHdl.AddressHandler),
		Courier: ch.(*courierHdl.CourierHandler),
	}
}
