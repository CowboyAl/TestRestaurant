package routes

import (
	//"auth/controllers"
	//"auth/utils/auth"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.blockrules.com/br/personal/TestRestaurant/controllers"
	"gitlab.blockrules.com/br/personal/TestRestaurant/utils/auth"
)

// Handlers - routing handlers
func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)

	r.HandleFunc("/", controllers.TestAPI).Methods("GET")
	r.HandleFunc("/api", controllers.TestAPI).Methods("GET")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/menu", controllers.ListMenu).Methods("GET")

	// Secure

	// Create Order
	r.HandleFunc("/customers/{userid}/orders/{itemid}", controllers.CreateOrder).Methods("POST")

	// List Customer Orders
	r.HandleFunc("/customers/{userid}/orders", controllers.ListCustomerOrders).Methods("Get")

	// List All Orders
	r.HandleFunc("/orders", controllers.ListOrders).Methods("GET")

	// GET orders/{orderid}/price - get total price for an order
	r.HandleFunc("/orders/{orderid}/price", controllers.GetOrderPrice).Methods("GET")

	// GET orders/{orderid}/time - get time of delivery for an order
	r.HandleFunc("/orders/{orderid}/time", controllers.GetOrderTime).Methods("GET")

	//GET customers/{userid}/orders - list existing orders

	//r.HandleFunc("/register", controllers.CreateUser).Methods("POST"
	r.HandleFunc("/users", controllers.FetchUsers).Methods("GET")

	r.Use(auth.JwtVerify)
	r.HandleFunc("/user", controllers.FetchUsers).Methods("GET")

	// Auth route
	//s := r.PathPrefix("/auth").Subrouter()
	//s.Use(auth.JwtVerify)
	r.HandleFunc("/user", controllers.FetchUsers).Methods("GET")
	r.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")
	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
