# About
The source code for the seeder of the customer service database

# Environment variables
- __OUTPUT_MEDIUM__: where to write the output to
    - FILE (default): write output as raw SQL to file
    - DB: write output to database
- __OUTPUT_DESTINATION__: depending on the medium, the value is either a connection string (_DB_) or a 
file system path (_FILE_)
- ROWS_NUMBER: a list of comma separated key-value pairs (`customers:n1,transactions:n2`) used to set the number of rows 
to insert for the two tables, the number can be omitted to use the default (1 million), omitting a table will opt it out of seeding. 
This parameter __needs to be specified__.

