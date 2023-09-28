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

	// rents
	rentHdl "book-nest/internal/handlers/rent"
	rentRepo "book-nest/internal/repositories/rent"
	rentSrv "book-nest/internal/services/rent"

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
	Rent    *rentHdl.RentHandler
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

	// rents
	rr := rentRepo.NewRentRepository()
	rs := rentSrv.NewRentService(rr, br, ur, db, gomail, midtrans)
	rh := rentHdl.NewRentHandler(rs)

	// address
	adr := addressRepo.NewAddressRepository()
	ads := addressSrv.NewAddressService(adr, db)
	adh := addressHdl.NewAddressHandler(ads)

	// couriers
	cr := courierRepo.NewCourierRepository()
	cs := courierSrv.NewCourierService(cr, adr, rr, db, biteship)
	ch := courierHdl.NewCourierHandler(cs)

	return &Presenter{
		Auth:    ah.(*authHdl.AuthHandler),
		User:    uh.(*userHdl.UserHandler),
		Book:    bh.(*bookHdl.BookHandler),
		Rent:    rh.(*rentHdl.RentHandler),
		Address: adh.(*addressHdl.AddressHandler),
		Courier: ch.(*courierHdl.CourierHandler),
	}
}
