# Reading and writing to the request context

Every `http.Request` has a context embedded in it. We store key/value pairs during the lifetime of the request

Values stored in the context has the type `any`, so we need type assertion before using it

Remember to use your **own custom type** for the request context key to avoid confusion between your code and 3rd-party packages
