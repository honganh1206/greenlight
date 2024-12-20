# Parsing Query String Parameters

- Our query parameters will look like this: `/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=-year`
- The `sort` parameters use the `-` character to denote **descending sorting order**, so `sort=-title` imples a descending sort.
- The `genres` parameter must accept **comma-separated** values
- The `page` and `page_size` parameters must accept numbers so we will *convert these string query values to Go int types* + Some validation for negative values.
- We need **default values** in case parameters are not provided.

## Helper functions
