shopifyStoreAPI is an application that provides basic RESTful online store functionality.
It was made as a challenge for Shopify Winter Internship 2019.

API
1. Store

/store
GET - returns all available stores
POST - create a new store

/store/:storeid
GET - get specified store with id
PATCH - update a specific store's information
DELETE - delete a specific store

All the endpoints below follow the same pattern of GET, POST, PATCH, DELETE as the store endpoints described above.

2. Products

/store/:storeid/products
GET, POST

/store/:storeid/products/:productid
GET, PATCH, DELETE

3. Items

/store/:storeid/products/:productid/items
GET, POST

/store/:storeid/products/:productid/items/:itemid
GET, DELETE

/store/:storeid/products/:productid/items/:itemid/order?id=yourOrderId
POST - Adds an item with :itemid to the existing order with a specified ID
DELETE - Removes an item with :itemid from the existing order with a specified ID

4. Orders

/store/:storeid/orders
GET - returns all non-empty orders
POST - creates a new empty order

/store/:storeid/orders/:orderid
GET - get a specific order
DELETE - delete a specific order (all the related items must be unattached/removed beforehand)