# GymShop

# REST API

The REST API to the example app is described below.

## Create an order

### Request

`POST /orders/`

    curl -X POST localhost:9003/orders -H "Content-Type: application/json" -d '{"customerName": "TestCustomer", "total": 12001}'
}'


### Response

    HTTP/1.1 200 OK

## Get orders list

### Request

`GET /orders/`

    curl localhost:9003/orders
}'


### Response

    HTTP/1.1 200 OK
    [{"id":1,"customerName":"TestCustomer","total":12001,"order_packs":null}]

## TODO
Currently, the pack sizes and their quantity are printed in the console, next step is to add them to the response when an order is created.
