package main

/*

	Author: Daniil Furmanov
	Date: September 23, 2018
	Purpose: Code challenge for Shopify Winter Internship 2019

*/

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/lars"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	DbHost     = "db"
	DbUser     = "postgres-dev"
	DbPassword = "mysecretpassword"
	DbName     = "dev"

	OkResponse = "{\"status\":\"ok\"}"
)

var (
	ErrInvalidId        = errors.New("invalid id")
	ErrInvalidStoreId   = errors.New("invalid store id")
	ErrInvalidProductId = errors.New("invalid product id")
	ErrInvalidOrderId   = errors.New("invalid order id")
)

type Store struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Product struct {
	Id      int     `json:"id"`
	StoreId int     `json:"store_id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
}

type Item struct {
	Id        int     `json:"id"`
	ProductId int     `json:"product_id"`
	StoreId   int     `json:"store_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
}

type Order struct {
	Id      int     `json:"id"`
	Total   float64 `json:"total"`
	Created string  `json:"created"`
	StoreId int     `json:"store_id"`
	Items   []Item  `json:"items"`
}

type OrderMap map[int]*Order

type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"data"`
}

var db *sql.DB

// Helper functions

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func errorMessage(msg string) GenericResponse {
	return GenericResponse{"fail", msg}
}

func notFound(c lars.Context) {
	c.JSON(http.StatusNotFound, errorMessage("not found"))
}

func internalError(c lars.Context, err error) {
	c.JSON(http.StatusInternalServerError, errorMessage(err.Error()))
}

func invalidRequest(c lars.Context, err error) {
	c.JSON(http.StatusUnprocessableEntity, errorMessage(err.Error()))
}

func okResponse(c lars.Context) {
	r := c.Response()
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(http.StatusOK)
	r.WriteString(OkResponse)
}

func okResponseId(c lars.Context, id int) {
	c.JSON(http.StatusOK, GenericResponse{"ok", strconv.Itoa(id)})
}

// Stores handlers

func getStores(c lars.Context) {
	rows, err := db.Query("SELECT * FROM stores LIMIT 100")
	checkErr(err)

	stores := make([]Store, 0)

	for rows.Next() {
		var id int
		var name string
		var description string
		err = rows.Scan(&id, &name, &description)
		if err != nil {
			internalError(c, err)
			return
		}
		stores = append(stores, Store{id, name, description})
	}

	c.JSON(http.StatusOK, stores)
}

func getStore(c lars.Context) {
	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	const q = `SELECT * FROM stores WHERE id=$1`

	rows, err := db.Query(q, storeId)
	if err != nil {
		internalError(c, err)
		return
	}

	for rows.Next() {
		var id int
		var name string
		var description string
		err = rows.Scan(&id, &name, &description)
		if err != nil {
			internalError(c, err)
			return
		}
		c.JSON(http.StatusOK, Store{id, name, description})
		return
	}

	notFound(c)
}

func postStore(c lars.Context) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		invalidRequest(c, err)
		return
	}
	defer c.Request().Body.Close()

	var store Store

	err = json.Unmarshal(body, &store)
	if err != nil {
		internalError(c, err)
		return
	}

	const q = `INSERT INTO stores(name, description) VALUES ($1, $2) RETURNING id`

	var id int
	err = db.QueryRow(q, store.Name, store.Description).Scan(&id)
	if err != nil {
		internalError(c, err)
		return
	}

	okResponseId(c, id)
}

func deleteStore(c lars.Context) {
	id := c.Param("id")
	if len(id) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	const q = `DELETE FROM stores WHERE id=$1`

	result, err := db.Exec(q, id)
	if err != nil {
		internalError(c, err)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil && affected < 1 {
		notFound(c)
		return
	}

	okResponse(c)
}

func updateStore(c lars.Context) {
	id := c.Param("id")
	if len(id) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		invalidRequest(c, err)
		return
	}
	defer c.Request().Body.Close()

	var store Store

	err = json.Unmarshal(body, &store)
	if err != nil {
		internalError(c, err)
		return
	}

	const q = `UPDATE stores SET name = $2, description = $3 WHERE id = $1`

	result, err := db.Exec(q, id, store.Name, store.Description)
	if err != nil {
		internalError(c, err)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil && affected < 1 {
		notFound(c)
		return
	}

	okResponse(c)
}

// Products handlers

func getProducts(c lars.Context) {
	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	const q = `SELECT id, store_id, name, price FROM products WHERE store_id=$1 LIMIT 100`

	rows, err := db.Query(q, storeId)
	checkErr(err)

	products := make([]Product, 0)

	for rows.Next() {
		var id int
		var storeId int
		var name string
		var price float64
		err = rows.Scan(&id, &storeId, &name, &price)
		if err != nil {
			internalError(c, err)
			return
		}
		products = append(products, Product{id, storeId, name, price})
	}

	c.JSON(http.StatusOK, products)
}

func getProduct(c lars.Context) {
	id := c.Param("productid")
	if len(id) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	const q = `SELECT id, store_id, name, price FROM products WHERE id=$1 AND store_id=$2`

	rows, err := db.Query(q, id, storeId)
	if err != nil {
		internalError(c, err)
		return
	}

	for rows.Next() {
		var id int
		var storeId int
		var name string
		var price float64
		err = rows.Scan(&id, &storeId, &name, &price)
		if err != nil {
			internalError(c, err)
			return
		}
		c.JSON(http.StatusOK, Product{id, storeId, name, price})
		return
	}

	notFound(c)
}

func postProduct(c lars.Context) {
	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		invalidRequest(c, err)
		return
	}
	defer c.Request().Body.Close()

	var product Product

	err = json.Unmarshal(body, &product)
	if err != nil {
		internalError(c, err)
		return
	}

	const q = `INSERT INTO products(name, price, store_id) VALUES ($1, $2, $3) RETURNING id`

	var id int
	priceStr := strconv.FormatFloat(product.Price, 'f', 3, 64)
	err = db.QueryRow(q, product.Name, priceStr, storeId).Scan(&id)
	if err != nil {
		internalError(c, err)
		return
	}

	okResponseId(c, id)
}

func deleteProduct(c lars.Context) {
	id := c.Param("productid")
	if len(id) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	const q = `DELETE FROM products WHERE id=$1 AND store_id=$2`

	result, err := db.Exec(q, id, storeId)
	if err != nil {
		internalError(c, err)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil && affected < 1 {
		notFound(c)
		return
	}

	okResponse(c)
}

func updateProduct(c lars.Context) {
	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	productId := c.Param("productid")
	if len(productId) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		invalidRequest(c, err)
		return
	}
	defer c.Request().Body.Close()

	var product Product

	err = json.Unmarshal(body, &product)
	if err != nil {
		internalError(c, err)
		return
	}

	const q = `UPDATE products SET name = $3, price = $4 WHERE id=$1 AND store_id=$2`

	result, err := db.Exec(q, productId, storeId, product.Name, product.Price)
	if err != nil {
		internalError(c, err)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil && affected < 1 {
		notFound(c)
		return
	}

	okResponse(c)
}

// Items handlers

func getItems(c lars.Context) {
	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	productId := c.Param("productid")
	if len(productId) < 1 {
		invalidRequest(c, ErrInvalidProductId)
		return
	}

	const q = `SELECT items.id, items.product_id, items.store_id, products.name, products.price
 FROM items JOIN products ON items.product_id=products.id WHERE items.store_id = $1 AND items.product_id=$2 LIMIT 100`

	rows, err := db.Query(q, storeId, productId)
	checkErr(err)

	products := make([]Item, 0)

	for rows.Next() {
		var id int
		var productId int
		var storeId int
		var name string
		var price float64
		err = rows.Scan(&id, &productId, &storeId, &name, &price)
		if err != nil {
			internalError(c, err)
			return
		}
		products = append(products, Item{id, productId, storeId, name, price})
	}

	c.JSON(http.StatusOK, products)
}

func getItem(c lars.Context) {
	id := c.Param("itemid")
	if len(id) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	productId := c.Param("productid")
	if len(productId) < 1 {
		invalidRequest(c, ErrInvalidProductId)
		return
	}

	const q = `SELECT items.id, items.product_id, items.store_id, products.name, products.price 
FROM items JOIN products ON items.product_id=products.id 
WHERE items.id=$1 AND items.store_id=$2 AND items.product_id=$3`

	rows, err := db.Query(q, id, storeId, productId)
	if err != nil {
		internalError(c, err)
		return
	}

	for rows.Next() {
		var id int
		var productId int
		var storeId int
		var name string
		var price float64
		err = rows.Scan(&id, &productId, &storeId, &name, &price)
		if err != nil {
			internalError(c, err)
			return
		}
		c.JSON(http.StatusOK, Item{id, productId, storeId, name, price})
		return
	}

	notFound(c)
}

func postItem(c lars.Context) {
	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	productId := c.Param("productid")
	if len(productId) < 1 {
		invalidRequest(c, ErrInvalidProductId)
		return
	}

	const q = `INSERT INTO items(store_id, product_id) VALUES ($1, $2) RETURNING id`

	var id int
	err := db.QueryRow(q, storeId, productId).Scan(&id)
	if err != nil {
		internalError(c, err)
		return
	}

	okResponseId(c, id)
}

func orderItem(c lars.Context) {
	id := c.Param("itemid")
	if len(id) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	productId := c.Param("productid")
	if len(productId) < 1 {
		invalidRequest(c, ErrInvalidProductId)
		return
	}

	queryParams := c.QueryParams()
	orderId := queryParams.Get("id")
	if len(orderId) < 1 {
		invalidRequest(c, ErrInvalidOrderId)
		return
	}

	// if it's not a POST method, then we do not need to set a new order id
	// because an item should not be added to the order
	extendOrder := c.Request().Method == "POST"

	var result sql.Result
	var err error

	// update order id of the item
	if extendOrder {
		const q = `UPDATE items SET order_id = $4 WHERE id=$1 AND store_id=$2 AND product_id=$3 AND order_id IS NULL`
		result, err = db.Exec(q, id, storeId, productId, orderId)
	} else {
		const q = `UPDATE items SET order_id = NULL WHERE id=$1 AND store_id=$2 AND product_id=$3 AND order_id IS NOT NULL`
		result, err = db.Exec(q, id, storeId, productId)
	}

	if err != nil {
		internalError(c, err)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil && affected < 1 {
		c.JSON(http.StatusConflict, errorMessage("item is unavailable"))
		return
	}

	// update total price of the order
	// if we are not extending the order, then we are subtracting from the price
	var q2 string

	if extendOrder {
		q2 = `UPDATE orders
SET total=total+subquery.price
FROM (SELECT products.id, products.price FROM products WHERE products.id=$1) AS subquery
WHERE orders.id=$2`
	} else {
		q2 = `UPDATE orders
SET total=total-subquery.price
FROM (SELECT products.id, products.price FROM products WHERE products.id=$1) AS subquery
WHERE orders.id=$2`
	}

	result, err = db.Exec(q2, productId, orderId)
	if err != nil {
		internalError(c, err)
		return
	}
	affected, err = result.RowsAffected()
	if err != nil && affected < 1 {
		notFound(c)
		return
	}

	okResponse(c)
}

func deleteItem(c lars.Context) {
	id := c.Param("itemid")
	if len(id) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	productId := c.Param("productid")
	if len(productId) < 1 {
		invalidRequest(c, ErrInvalidProductId)
		return
	}

	const q = `DELETE FROM items WHERE id=$1 AND store_id=$2 AND product_id=$3`

	result, err := db.Exec(q, id, storeId, productId)
	if err != nil {
		internalError(c, err)
		return
	}
	affected, err := result.RowsAffected()
	if err != nil && affected < 1 {
		notFound(c)
		return
	}

	okResponse(c)
}

// Orders handlers

func getOrders(c lars.Context) {
	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	const q = `
SELECT orders.id, orders.total, orders.created, orders.store_id, items.id, items.product_id, products.name, products.price
FROM orders 
JOIN items ON items.order_id=orders.id 
JOIN products ON products.id=items.product_id
WHERE orders.store_id=$1 LIMIT 100`

	rows, err := db.Query(q, storeId)
	checkErr(err)

	orders := make(OrderMap)

	for rows.Next() {
		var id int
		var total float64
		var created string
		var storeId int
		var itemId int
		var productId int
		var productName string
		var productsPrice float64

		err = rows.Scan(&id, &total, &created, &storeId, &itemId, &productId, &productName, &productsPrice)
		if err != nil {
			internalError(c, err)
			return
		}

		if order, ok := orders[id]; ok {
			order.Items = append(order.Items, Item{itemId, productId, storeId, productName, productsPrice})
		} else {
			items := []Item{{itemId, productId, storeId, productName, productsPrice}}
			orders[id] = &Order{
				Id:      id,
				Total:   total,
				Created: created,
				StoreId: storeId,
				Items:   items,
			}
		}
	}

	c.JSON(http.StatusOK, orders)
}

func getOrder(c lars.Context) {
	id := c.Param("id")
	if len(id) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	const q = `
SELECT orders.id, orders.total, orders.created, orders.store_id, items.id, items.product_id, products.name, products.price
FROM orders 
JOIN items ON items.order_id=orders.id 
JOIN products ON products.id=items.product_id
WHERE orders.store_id=$2 AND orders.id=$1`

	rows, err := db.Query(q, id, storeId)
	if err != nil {
		internalError(c, err)
		return
	}

	var order *Order

	for rows.Next() {
		var id int
		var total float64
		var created string
		var storeId int
		var itemId int
		var productId int
		var productName string
		var productsPrice float64

		err = rows.Scan(&id, &total, &created, &storeId, &itemId, &productId, &productName, &productsPrice)
		if err != nil {
			internalError(c, err)
			return
		}

		if order != nil {
			order.Items = append(order.Items, Item{itemId, productId, storeId, productName, productsPrice})
		} else {
			items := []Item{{itemId, productId, storeId, productName, productsPrice}}
			order = &Order{
				Id:      id,
				Total:   total,
				Created: created,
				StoreId: storeId,
				Items:   items,
			}
		}
	}

	if order == nil {
		notFound(c)
		return
	}

	c.JSON(http.StatusOK, order)
}

func postOrder(c lars.Context) {
	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	const q = `INSERT INTO orders(total, store_id) VALUES ($1,$2) RETURNING id`

	var id int
	err := db.QueryRow(q, 0, storeId).Scan(&id)
	if err != nil {
		internalError(c, err)
		return
	}

	okResponseId(c, id)
}

func deleteOrder(c lars.Context) {
	id := c.Param("id")
	if len(id) < 1 {
		invalidRequest(c, ErrInvalidId)
		return
	}

	storeId := c.Param("storeid")
	if len(storeId) < 1 {
		invalidRequest(c, ErrInvalidStoreId)
		return
	}

	const q = `DELETE FROM orders WHERE id=$1 AND store_id=$2`

	result, err := db.Exec(q, id, storeId)
	if err != nil {
		internalError(c, err)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil && affected < 1 {
		notFound(c)
		return
	}

	okResponse(c)
}

func main() {
	var err error

	router := lars.New()

	// stores
	stores := router.Group("/stores")
	stores.Get("/", getStores)
	stores.Get("", getStores)
	stores.Post("/", postStore)
	stores.Post("", postStore)
	oneStore := stores.Group("/:storeid")
	oneStore.Get("/", getStore)
	oneStore.Get("", getStore)
	oneStore.Delete("/", deleteStore)
	oneStore.Delete("", deleteStore)
	oneStore.Patch("/", updateStore)
	oneStore.Patch("", updateStore)
	// products
	products := oneStore.Group("/products")
	products.Get("/", getProducts)
	products.Get("", getProducts)
	products.Post("/", postProduct)
	products.Post("", postProduct)
	oneProduct := products.Group("/:productid")
	oneProduct.Get("/", getProduct)
	oneProduct.Get("", getProduct)
	oneProduct.Delete("/", deleteProduct)
	oneProduct.Delete("", deleteProduct)
	oneProduct.Patch("/", updateProduct)
	oneProduct.Patch("", updateProduct)
	// items
	items := oneProduct.Group("/items")
	items.Get("/", getItems)
	items.Get("", getItems)
	items.Post("/", postItem)
	items.Post("", postItem)
	oneItem := items.Group("/:itemid")
	oneItem.Post("/order", orderItem)
	oneItem.Delete("/order", orderItem)
	oneItem.Get("/", getItem)
	oneItem.Get("", getItem)
	oneItem.Delete("/", deleteItem)
	oneItem.Delete("", deleteItem)
	// orders
	orders := oneStore.Group("/orders")
	orders.Get("/", getOrders)
	orders.Get("", getOrders)
	orders.Post("/", postOrder)
	orders.Post("", postOrder)
	oneOrder := orders.Group("/:id")
	oneOrder.Get("/", getOrder)
	oneOrder.Get("", getOrder)
	oneOrder.Delete("/", deleteOrder)
	oneOrder.Delete("", deleteOrder)

	server := &http.Server{Addr: ":8080", Handler: router.Serve()}

	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		DbHost, DbUser, DbPassword, DbName)
	db, err = sql.Open("postgres", dbinfo)
	checkErr(err)

	defer db.Close()

	log.Println("running the server...")

	defer server.Close()
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Startup failed")
	}
}
