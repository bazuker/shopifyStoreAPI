# Code for Shopify Winter Internship challenge 2019.
shopifyStoreAPI is an application that provides basic RESTful online store functionality.

## API doc
### 1. Store

__/store__

GET - returns all available stores.

POST - create a new store.
```json
{
  "name": "My cool store",
  "description": "Store description"
}
```

__/store/:storeid__

GET - get specified store with id.
PATCH - update a specific store's information.
DELETE - delete a specific store.

All the endpoints below follow the same pattern of GET, POST, PATCH, DELETE as the store endpoints described above.

### 2. Products

__/store/:storeid/products__

GET

POST
```json
{
  "name": "Phone",
  "price": 999.9
}
```


__/store/:storeid/products/:productid__

GET, PATCH, DELETE

### 3. Items

__/store/:storeid/products/:productid/items__

GET

POST
```
json body is not required
```

__/store/:storeid/products/:productid/items/:itemid__

GET, DELETE

__/store/:storeid/products/:productid/items/:itemid/order?id=yourOrderId__

POST - Adds an item with :itemid to the existing order with a specified ID.
DELETE - Removes an item with :itemid from the existing order with a specified ID.

### 4. Orders

__/store/:storeid/orders__

GET - returns all non-empty orders.

POST - creates a new empty order.
```
json body is not required
```

__/store/:storeid/orders/:orderid__

GET - get a specific order.
DELETE - delete a specific order (all the related items must be unattached/removed beforehand).
