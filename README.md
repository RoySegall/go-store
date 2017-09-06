# Requirements

First, you'll need [RethinkDB](https://www.rethinkdb.com) running in the background. After installing RethinkDB just run
```
rethinkdb --http-port 8877
```

# OpenSSL keys.
```
openssl genrsa -out private.key 2048
openssl rsa -in private.key -pubout > public.key
```

**Clone the project to `store` and not `go-store`**

## Install

```bash
./install
./migrate
```

## Firing up the store
```bash
./store
```

## Data structure
There are two primary entities:
* users - Will store the users in the store. The user objects contain three sub data models:
  * Access token - the access token of the user, refresh token and when it's expired.
  * Cart - Contains the current cart which the user added items to
  * Past carts - List of the past carts has purchased.
* Items - The items in the store.

## How to interact with the backend

### items
All the items are lists under `API/items` and have the next structure:
```json
{
    Id: "72f53a6e-260f-459f-9252-75ef298f08c7",
    Title: "Overcorrecting for pat failures",
    Description: "Working emotional volatility into your decision tree",
    Price: 24.5,
    Image: "images/book13.jpg"
}
```

### User
Due to security issues, we won't have a list of the users. Instead, you could ask for the user info in `API/user`. To get the access to the current user, an access token is needed in the header(Will be elaborated below). The user object looks like that:
```json
{
    "data": {
        "Id": "49770a56-9284-4ce9-8c13-8acd76f9f82f",
        "Username": "adminadmin",
        "Password": "$2a$14$u2JhQ8bxoAtc1Rc6ncye.uwudc0jPuNk7Mz4Oj66NiqK5P2Vxzva.",
        "Email": "",
        "Image": "",
        "Role": {
            "Title": "Member"
        },
        "Token": {
            "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtb2RlIjoibWFpbl90b2tlbiIsIm5hbWUiOiJhZG1pbmFkbWluIiwidGltZSI6MTUyNTMxODMzMn0.8OqATlIVgjqCUBVdeiE88j27aV2RusNLJHEguZhHum4",
            "Expire": 1503606679,
            "RefreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtb2RlIjoicmVmcmVzaF90b2tlbiIsIm5hbWUiOiJhZG1pbmFkbWluIiwidGltZSI6MTUyNTM4MTk5MH0.wxhFTOKq0UCWc0QnVs1t_pGAt_jDz68yI-aNxUCr5MU"
        },
        "Cart": {
            "Items": []
        },
        "PastCarts": []
    }
}
```
In order to register a user you'll to do a post method with the following deailts:
* username
* password

For now, the other details are not handled.

### Access token
To get an access token, you'll need to do a post request to `api/user/login` with the payload of
* username
* password

Access token will be valid for only 24 hours. If you want to get a new one, you'll need to do a `POST` request to `/API/user/token_refresh` with the payload of `refresh_token` as the only value. You'll get a new refresh token.

### Cart
Adding items to cart can be done with a `POST` request to `/API/cart/items` when the `item_id` is the only value. Remove that Item can be done with a `DELETE` method with the same payload. To archive the cart you'll need to send a `DELETE` request to `/API/cart`.
