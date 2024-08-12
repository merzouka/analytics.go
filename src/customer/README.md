# Environment variables
__TRANSACTION_SERVICE__: url path of transactions service
__DB_URL__: Postgres database url
__CACHE_URL__: Redis cache url
__CACHE_PASSWORD__: Redis cache password
__LOGS_PATH__ [optional, default: './logs']: Path of logs file
__SERVICE_NAME__ [optional]: Name to be used as prefix in logs
__MODE__: (CACHE/DB) Data retrieval mode
- __DB__ query all data from database
- __CACHE__ use cache as intermediate to database

# Endpoints
- `customers/sorted`: the sorted list of customer by the number of transactions
- `customers/[id]/transactions/total`: the total number of transactions for the customer
- `customers/[id]/transactions`: transactions made by the customer

# How it works
- The endpoints for retrieving transactions and transaction total are proxies for the transactions service. 
- The sorted endpoint first fetches ids of customers sorted by the total amount of transactions the customer has made, 
customers that have made no transaction are sorted by their id
