package controllers

import (
	//"auth/models"
	//"auth/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gitlab.blockrules.com/br/personal/TestRestaurant/models"
	"gitlab.blockrules.com/br/personal/TestRestaurant/utils"
	"golang.org/x/crypto/bcrypt"
)

// ErrorResponse - return for errors
type ErrorResponse struct {
	Err string
}

type error interface {
	Error() string
}

var db = utils.ConnectDB()

// TestAPI - root test
func TestAPI(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Running"))
}

// Login - User login function
func Login(w http.ResponseWriter, r *http.Request) {
	//decoder := json.NewDecoder(r.Body)
	//var data models.LoginRequest
	//err := decoder.Decode(&data)
	//if err != nil {
	//	panic(err)
	//}

	//username := data.Username
	//password := data.Password

	err := r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	//user := &models.LoginRequest{}
	//body, err := ioutil.ReadAll(r.Body)
	//fmt.Println("Bodyd:", string(body))
	//err := json.NewDecoder(r.Body).Decode(user)
	//user := models.LoginRequest{}
	//if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//fmt.Println("got", user.Username+" "+user.Password)
	//resp := FindOne(user.Username, user.Password)
	resp := FindOne(username, password)
	json.NewEncoder(w).Encode(resp)
}

// FindOne - looks up a user in the database
func FindOne(username, password string) map[string]interface{} {
	user := &models.User{}

	if err := db.Where("Username = ?", username).First(user).Error; err != nil {
		var resp = map[string]interface{}{"status": false, "message": "User not found"}
		return resp
	}
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	//errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	//if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
	if password != user.Password {
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	}

	tk := &models.Token{
		UserID:   user.ID,
		Username: user.Username,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user
	return resp
}

//CreateUser function -- create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		err := ErrorResponse{
			Err: "Password Encryption  failed",
		}
		json.NewEncoder(w).Encode(err)
	}

	user.Password = string(pass)

	createdUser := db.Create(user)
	var errMessage = createdUser.Error

	if createdUser.Error != nil {
		fmt.Println(errMessage)
	}
	json.NewEncoder(w).Encode(createdUser)
}

// FetchUsers - lists all users
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.Preload("auths").Find(&users)

	json.NewEncoder(w).Encode(users)
}

// UpdateUser - updates user data
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	params := mux.Vars(r)
	var id = params["id"]
	db.First(&user, id)
	json.NewDecoder(r.Body).Decode(user)
	db.Save(&user)
	json.NewEncoder(w).Encode(&user)
}

// DeleteUser - deletes a user from the database
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	db.First(&user, id)
	db.Delete(&user)
	json.NewEncoder(w).Encode("User deleted")
}

// GetUser - gets user data
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	db.First(&user, id)
	json.NewEncoder(w).Encode(&user)
}

// ListMenu - lists all menu items
func ListMenu(w http.ResponseWriter, r *http.Request) {
	var items []models.MenuItem
	db.Preload("menu_items").Find(&items)

	json.NewEncoder(w).Encode(items)
}

// ListOrders - lists all menu items
func ListOrders(w http.ResponseWriter, r *http.Request) {
	var orders []models.Order
	//db.Preload("order").Find(&orders)
	db.Find(&orders)

	json.NewEncoder(w).Encode(orders)
}

// CreateOrder - lists all menu items
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	// POST customers/{userid}/orders/{itemid} - place an order
	vars := mux.Vars(r)
	//fmt.Println("before")
	//fmt.Fprintf(w, "user: %s, item %s", vars["userid"], vars["itemid"])
	//fmt.Fprintf(w, "item: %s", vars["itemid"])

	var order models.Order
	order.UserID = vars["userid"]

	itemID, err := strconv.Atoi(vars["itemid"])
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid menu item"}
		json.NewEncoder(w).Encode(resp)
	}

	order.ItemID = itemID

	createdOrder := db.Create(&order)
	var errMessage = createdOrder.Error

	if createdOrder.Error != nil {
		fmt.Println(errMessage)
	}
	json.NewEncoder(w).Encode(createdOrder)
}

// GetOrderPrice - Gets the price for an order
// GET orders/{orderid}/price - get total price for an orde
//- Base (Item) cost = Rate \* Quantity
//- Taxes - 5 % on base cost
//- Delivery charge - Rs 1 \* distance b/w customer and restaurant.
func GetOrderPrice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var order models.Order
	err := db.Where("id = ?", vars["orderid"]).First(&order).Error
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Order not found"}
		json.NewEncoder(w).Encode(resp)
	}

	//base := order.

	var resp = map[string]interface{}{"status": true, "message": "Invalid menu item"}
	json.NewEncoder(w).Encode(resp)
}

// GetOrderTime - Gets the time for an order
// GET orders/{orderid}/time - get time of delivery for an order
func GetOrderTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var order models.Order
	err := db.Where("id = ?", vars["orderid"]).First(&order)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Order not found"}
		json.NewEncoder(w).Encode(resp)
	}

	json.NewEncoder(w).Encode(order)
}
