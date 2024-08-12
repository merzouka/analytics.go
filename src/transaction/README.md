# Environment variables
- __PRODUCT_SERVICE__: path to products service
- __DB_URL__: Postgres database url
- __CACHE_URL__: Redis cache url
- __CACHE_PASSWORD__: Redis cache password
- __MODE__: (CACHE/DB) Data retrieval mode
    - __DB__ query all data from database
    - __CACHE__ use cache as intermediate to database

# Endpoints
- `transactions/:id/products`: get products related to transaction
- `transactions/:id`: get transaction with id
- `transactions?ids`: get transactions with ids
- `transactions/total?ids`: calculate the sum of value for the ids
