package http

import (
	"book-nest/clients/biteship"
	"book-nest/clients/gomail"
	"book-nest/clients/midtrans"

	"gorm.io/gorm"

	i "book-nest/internal/interfaces"

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
	// initiate variable
	var (
		// users
		userRepository i.UserRepository
		userService    i.UserService
		userHandler    i.UserHandler

		// auth
		authRepository i.AuthRepository
		authService    i.AuthService
		authHandler    i.AuthHandler

		// books
		bookRepository i.BookRepository
		bookService    i.BookService
		bookHandler    i.BookHandler

		// orders
		orderRepository i.OrderRepository
		orderService    i.OrderService
		orderHandler    i.OrderHandler

		// address
		addressRepository i.AddressRepository
		addressService    i.AddressService
		addressHandler    i.AddressHandler

		// couriers
		courierRepository i.CourierRepository
		courierService    i.CourierService
		courierHandler    i.CourierHandler
	)
	// clients
	gomail := gomail.NewGomailClient()
	midtrans := midtrans.NewMidtransClient()
	biteship := biteship.NewBiteshipClient()

	// users
	userRepository = userRepo.NewUserRepository()
	userService = userSrv.NewUserService(userRepository, db, gomail)
	userHandler = userHdl.NewUserHandler(userService)

	// auth
	authRepository = authRepo.NewAuthRepository()
	authService = authSrv.NewAuthService(authRepository, userRepository, db)
	authHandler = authHdl.NewAuthHandler(authService)

	// books
	bookRepository = bookRepo.NewBookRepository()
	bookService = bookSrv.NewBookService(bookRepository, orderRepository, db)
	bookHandler = bookHdl.NewBookHandler(bookService)

	// orders
	orderRepository = orderRepo.NewOrderRepository()
	orderService = orderSrv.NewOrderService(orderRepository, bookRepository, userRepository, db, gomail, midtrans)
	orderHandler = orderHdl.NewOrderHandler(orderService)

	// address
	addressRepository = addressRepo.NewAddressRepository()
	addressService = addressSrv.NewAddressService(addressRepository, db)
	addressHandler = addressHdl.NewAddressHandler(addressService)

	// couriers
	courierRepository = courierRepo.NewCourierRepository()
	courierService = courierSrv.NewCourierService(courierRepository, addressRepository, orderRepository, bookRepository, db, biteship)
	courierHandler = courierHdl.NewCourierHandler(courierService)

	return &Presenter{
		Auth:    authHandler.(*authHdl.AuthHandler),
		User:    userHandler.(*userHdl.UserHandler),
		Book:    bookHandler.(*bookHdl.BookHandler),
		Order:   orderHandler.(*orderHdl.OrderHandler),
		Address: addressHandler.(*addressHdl.AddressHandler),
		Courier: courierHandler.(*courierHdl.CourierHandler),
	}
}
