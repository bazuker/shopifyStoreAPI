# Code for Shopify Winter Internship challenge 2019.
shopifyStoreAPI is an application that provides basic RESTful online store functionality.

## Installation with Docker
compose yourself
```bash
$ git clone https://github.com/kisulken/shopifyStoreAPI
$ cd shopifyStoreAPI
$ docker-compose up
```
or run from the existing image
```bash
docker run -it kisulken/shopifystoreapi:v5
```

## API endpoints
| endpoint       | methods       |
| -------------  |:-------------:| 
| /store         | GET, POST     | 
| /store/:storeid| GET, PATCH, DELETE      |
| /store/:storeid/products  | GET, POST      |
| /store/:storeid/products/:productid  | GET, PATCH, DELETE      |
| /store/:storeid/products/:productid/items  | GET,  POST      |
| /store/:storeid/products/:productid/items/:itemid  | GET,  DELETE      |
| /store/:storeid/products/:productid/items/:itemid/order  | POST,  DELETE      |
| /store/:storeid/orders | GET,  POST      |
| /store/:storeid/orders/:orderid  | GET,  DELETE      |

## API responses
All API endpoints will respond with a content-type application/json, corresponding http codes and contain a valid json data.

__GET__ methods will respond with requested data i.e store information, product information etc. in case of success
```
[
    {
        "id": 1,
        "name": "My very cool store",
        "description": "Descriptive description"
    }
]
```
"not found" in case if the requested data was not located in the database
```json
{
    "status": "fail",
    "data": "not found"
}
```

__POST__ methods in case of successful insertion will respond with
```json
{
    "status": "ok",
    "data": id
}
```
where *id* is a decimal number representing a unique identifier of the object

All other methods in case of success will respond with 
```json
{
    "status": "ok"
}
```
and in case of failure
```json
{
    "status": "fail",
    "data": "error message"
}
```

## API doc
### 1. Store

- __/store__

GET - returns all available stores.

POST - create a new store.
```json
{
  "name": "My cool store",
  "description": "Store description"
}
```


- __/store/:storeid__

GET - get specified store with id.

PATCH - update a specific store's information.

DELETE - delete a specific store.

All the endpoints below follow the same pattern of GET, POST, PATCH, DELETE as the store endpoints described above.

### 2. Products

- __/store/:storeid/products__

GET

POST
```json
{
  "name": "Phone",
  "price": 999.9
}
```


- __/store/:storeid/products/:productid__

GET, PATCH, DELETE

### 3. Items

- __/store/:storeid/products/:productid/items__

GET

POST
```
json body is not required
```

- __/store/:storeid/products/:productid/items/:itemid__

GET, DELETE

- __/store/:storeid/products/:productid/items/:itemid/order?id=yourOrderId__

POST - Adds an item with :itemid to the existing order with a specified ID automatically adding a price of the product to the total

DELETE - Removes an item with :itemid from the existing order with a specified ID automatically deducting the cost of the item

### 4. Orders

- __/store/:storeid/orders__

GET - returns all non-empty orders.

POST - creates a new empty order.
```
json body is not required
```

- __/store/:storeid/orders/:orderid__

GET - get a specific order.

DELETE - delete a specific order (all the related items must be unattached/removed beforehand).
