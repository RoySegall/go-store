# setting up

You'll need a Drupal environment: PHP 5.6+(7 and above is recomended), MySQL,
Apache, Drush and composer.

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
./store
```
