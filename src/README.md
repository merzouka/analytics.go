# What it does
The basic idea is there are three main services:
- Transaction
- Customer
- Product
The API has a couple of aggregate queries than involve all three services. Queries for the data are either run directly on the database, or 
through the medium of a cache. Performance is measured in each case, and the results are to be compared.

The cache in all three holds string representations of the object (transaction/customer/product) that have been accessed before 
special cases are documented in the read me of the specific service.
