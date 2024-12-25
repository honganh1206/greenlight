# Filtering, Sorting, and Pagination

What we will do:

- Return the details of **multiple resources** in a single JSON response
- Apply filter parameters to narrow down the returned data set
- Implement **full-text** search on your database using PostgreSQL **inbuilt** functionality
- Apply sort parameters to change the order of results
- Develop a pattern to support pagination on large data sets

---
 ## Parsing Query String Parameters

- Our query parameters will look like this: `/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=-year`
- The `sort` parameters use the `-` character to denote **descending sorting order**, so `sort=-title` imples a descending sort.
- The `genres` parameter must accept **comma-separated** values
- The `page` and `page_size` parameters must accept numbers so we will *convert these string query values to Go int types* + Some validation for negative values.
- We need **default values** in case parameters are not provided.

### Helper functions

- We have different helper functions like `readInt()`, `readCSV()` and `readString()` to extract and convert query parameters into specific types.
- Embedded struct like `data.Filters` is when we *include a type name without a field name*, so it is automatically treated as a field name.

### Validating query string parameters

Some rules:
- `page` between 1 and 1,000,000
- `page_size` between 1 and 100
- `sort` contains a known and supported value for our movie table like "id" and "-id"

---
## Listing data

- We need a `GetAll()` method, and we need to add timeout context for that too!

---
## Filtering Lists

- We are going to build a [reductive filter] (https://ux.stackexchange.com/questions/88993/inclusive-additive-vs-exclusive-reductive-filtering-how-to-differentiate) (TLDR: We provessively remove unwanted elements while preserving the important features, like refining gold from ore)

 - Our filtering feature will be DYNAMIC! Some of the example cases:
 1. List movies where the title is a case-insensitive *exact match* for 'black panther': `/v1/movies?title=black+panther`
2. List movies where the genres *includes* 'adventure'`/v1/movies?genres=adventure`
3. List movies where the title is a case-insensitive *exact match* for 'moana' AND the genres *include* both 'animation' AND 'adventure': `
/v1/movies?title=moana&genres=animation,adventure`

### Dynamic filtering

- Challenge: We need to work with different cases: No filters/Filters on either `title` and `genres`/Either of them.
- One option is to *build up the SQL query dynamically during rumtime*, but this approach will *potentially lead to messy code* especially for queries for multiple filter options.
- Our approach is going to be a **fixed** SQL query, with the aim to design the filters to *behave like they are optional*. Here is an example of that:

```sql
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
-- Either a case-sensitive match or empty - We skip the filter condition if empty
WHERE (LOWER(title) = LOWER($1) OR $1 = '')
-- Either the genres contain the value or the returned value will be an empty array
AND (genres @> $2 OR $2 = '{}')
ORDER BY id
```

---
## Full-Text Search

- We add support for *partial matches* by leveraging PostgreSQL's *full-text search* functionality. In detail, we use the `to_tsvector()` function to *convert text to searchable format* and split the movie title into **lexemes** (basic unit of meaning in a language - like a "dictionary form" of a word)

```sql
-- Output showing the positions of the words: 'fast':4 'run':1,2,5
SELECT to_tsvector('English', 'The runner is running fast runs');
```

- We also use `plainto_tsquery()` function to *take a search value and turn it into a formatted query term*. An example is the search value "The Club" would result in `'the' & 'club'`

- We also use the `@@` operator to check *whether the generated query term matches the lexemes*

### Alternatives: `STRPOS` and `ILIKE`

- The `STRPOS` function allows us to *check for the existence of a substring in a particular database field*.

```sql
WHERE (STRPOS(LOWER(title), LOWER($1)) > 0 OR $1 = '')
```

- However, the downside is that the result might be **unintuitive**:
  - Client side: If we are searching for `title=the`, the result would be both "THE Breakfast Club" and "Black PanTHEr".
  - Server side: There is no effective way to index the `title` column to see if the `STRPOS()` condition is met, so the query might require a full-table scan each time it is executed.

- The `ILIKE` function allows us to *find rows which match a specific (case-insensitive) pattern*.

```sql
WHERE (title ILIKE $1 OR $1 = '')
 ```

- This approach is better for both the client side (able to use prefix/suffix for search term like wildcard character `%`) and the server side (able to create an index with `pg_trgm` extension and GIN)

---
## Sorting Lists

- Problem: The values for `ORDER BY` must be generated *during runtime* based on the query string values from the client => Solution: We use `fmt.Sprintf()` to **interpolate** the dynamic values into our query string.

- Important note: The order of returned rows is guaranteed *only* by the rules from the `ORDER BY` clause a.k.a we need to *explicitly set the rules* if we want the returned rows to be ordered.

---
## Paginating Lists

- We can use `LIMIT` (set the maximum number of records for a query) and `OFFSET` (skip a specific number of rows) clauses to our query to page results.

- We return pagination metadata (counts, page size, total records, etc.) by using [window functions](https://www.postgresql.org/docs/current/tutorial-window.html) of PostgreSQL.
