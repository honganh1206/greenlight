# Parsing Query String Parameters

- Our query parameters will look like this: `/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=-year`
- The `sort` parameters use the `-` character to denote **descending sorting order**, so `sort=-title` imples a descending sort.
- The `genres` parameter must accept **comma-separated** values
- The `page` and `page_size` parameters must accept numbers so we will *convert these string query values to Go int types* + Some validation for negative values.
- We need **default values** in case parameters are not provided.

## Helper functions

- We have different helper functions like `readInt()`, `readCSV()` and `readString()` to extract and convert query parameters into specific types.
- Embedded struct like `data.Filters` is when we *include a type name without a field name*, so it is automatically treated as a field name.

## Validating query string parameters

Some rules:
- `page` between 1 and 1,000,000
- `page_size` between 1 and 100
- `sort` contains a known and supported value for our movie table like "id" and "-id"

## Listing data

- We need a `GetAll()` method, and we need to add timeout context for that too!
