# Environment variables
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
- `customers/[id]`: the customer with id
