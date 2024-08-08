# Endpoints
- `products?ids`: product with ids

# Environment variables
- __DB_URL__: postgres url
- __CACHE_URL__: redis url (host:port)
- __CACHE_PASSWORD__: redis password
- __MODE__: the mode to use for the api, values:
    - __CACHE__: utilize the cache to improve performance
        - variables: DB_URL, CACHE_URL, CACHE_PASSWORD
    - __DB__: run all queries against the database
        - variables: DB_URL
