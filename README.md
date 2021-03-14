## Description

*eshop* service implements api [specification](./api/openapi.yaml)

It supports two operations:
1. Create a cart with given line items. This operation will evaluate the cart with current promotions in place too.
2. Fetch a cart with given cart identifier. This operation will evaluate the cart with current promotions in place too.  


This *eshop* service implementation uses cockroach-db as database. The database stores following information in given tables.
1. current inventory stock in table eshop.inventory
2. current active promotions in table eshop.promotions
3. user/shopper carts in table eshop.carts

On setup, the database is pre-populated with state as advised in the test description with help of migration scripts at [here](./scripts/db/migrations)

## Unit Tests

Logic of a cart evaluation with current promotions is unit tested with test cases in [cart_test.go](./internal/eshop/cart_test.go)

```
make test
```

## Running on local machine

*Requirements*
* Mac/Linux based operating system
* docker & cocker-compose installed locally
* Go version >= 1.16
* git command line

### docker/docker-compose
Easiest & cleanest way to run this application locally is with help of *docker-compose*.

This will need following ports to be available locally
* 8080  - swagger-ui container
* 9080  - eshop api container
* 26257 - cockroachdb container

```
make docker-serve-all
```

>> a container named 'migrate' will run-to-complete and it will run the db migration scripts.

Browse to [swagger-ui](http://localhost:8080) to view the open-api specification and test the *eshop* service using easy to use interface. The api specification is documented with request and response examples.


### Running deps in docker network and api on host machine

To run *eshop* service dependencies in docker network, and run the service on host machine (for quick/continuous development).

```
make docker-deps

make serve
```

Once service is running browse to [swagger-ui](http://localhost:8080) to view the open-api specification and test the *eshop* service.

### Cleanup

```
make docker-clean
``` 





