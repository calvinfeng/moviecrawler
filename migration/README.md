# Database
## Setup
```
psql -d postgres -f ./migrations/000_setup.sql
```

## Create Table
```
psql -d moviecrawler -f ./migrations/001_create_movies.sql
```