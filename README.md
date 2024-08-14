# About
The goal of this project is to test out the change in performance when using a cache instead of a direct database access

# Usage
run the back-end using the script `init.sh` in the _config_ folder and run the front-end under _src/frontend_ which would allow you to 
specify the number of queries to be run by a proxy for each type, and observer the results in real time.

# Seeding
you can change the number of rows using the _ROWS_NUMBER_ environment variable in the _init.yaml_ configuration files; to find out more information about 
what constraints do apply view the respective seeding documentation under _seeders_.

# Observations
Oddly enough using a cache has made performance worse than running straight database queries, this is probably a flaw in my implementation, or due to the fact
that no models are retrieved from the cache (all queries are fresh). 
What is use is a for loop to query each object from the cache, maybe there is a better way of accessing the data.

