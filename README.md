# Restaurant Example

To start this app, perform the following step in order

1. Clone this repo to your machine
1. Confirm postgres is running
1. Create a postgres database called "restaurant"
1. Update the .env file in the root of the project with the proper passwords
1. cd into the project folder
1. Enter `go run main.go` to start server

## REQUIREMENTS

Restaurant role with features:

- Sign in/out as a restaurant.
- Get a list of available orders.
- Print invoice as a pdf file.

Customer role with feature:

- Sign in/out as a customer
- Order a food item from list of available items (populate some dummy items in the db)
- Get total fare of the food ordered which can be calculated as per following rule-
  - Base (Item) cost = Rate \* Quantity
  - Taxes - 5 % on base cost
  - Delivery charge - Rs 1 \* distance b/w customer and restaurant.
- Get estimated time which can be calculated as per following rule-
  - Food prep time = Unit prep time \* quantity
  - Delivery time - Assume driver to drive at rate of 40 kmph.
- Cancel an order.

Driver role with features -

- Sign in/out as driver
- Get an order information with the restaurant and delivery address.
- Update an order as picked up.
- Update an order as delivered.

## ASSUMPTIONS

I have made the following assumptions and simplifications:

1.  Everyone is a user. Currently I don't distinguish between customer, driver, and restaurant staff
1.  A customer can only order one item. Otherwise, I would have to manage complex order lists.
1.  I've harcoded the secret for the JWT bearer token. This should come from the environment or a database.
1.  I've populated the database with 1 user, bob, with a password of bobsecret
1.  Prices are in dollars. Travel cost is $.10/mile

## API

// POST login - returns a bearer token
// GET menu - list menu items

// POST customers/{userid}/orders/{itemid} - place an order
// GET customers/{userid}/orders - list existing orders

// GET orders/{orderid}/price - get total price for an order
// GET orders/{orderid}/time - get time of delivery for an order

// GET orders - list all orders
// PUT orders/{orderid}/pickedup - mark an order as picked up
// PUT orders/{orderid}/delivered - mark an order as delivered
// PUT orders/{orderid}/delete
// GET orders/{orderid}/print - print order to pdf

// GET users - list all users

## DATABASE

// Customers
// CustomerID, user name, password, address, distance

// MenuItems
// ItemID, Description, Price, prep time

// Orders
// OrderID, UserID, ItemID, PickedUp, Delivered

## ARCHITECTURE

API <=> Routers <=> Authentication <=> Controllers <=> Database
