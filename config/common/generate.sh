#!/bin/bash

read -p "namespace> " namespace
image="vmerv/analytics-seeder-$namespace"
db="$namespace""db"
read -p "output directory> " dir
read -p "rows number> " rows

sed "s|{{ns}}|$namespace|g" < init.yaml | sed "s|{{seeder-image}}|$image|g" | sed "s|{{seeder-rows}}|$rows|g" > "$dir/init.yaml"
sed "s|{{ns}}|$namespace|g" < db.yaml | sed "s|{{db-name}}|$db|g" > "$dir/db.yaml"
sed "s|{{ns}}|$namespace|g" < cache.yaml > "$dir/cache.yaml"
