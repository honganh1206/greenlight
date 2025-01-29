# User Model Setup and Registration

Here is how our `users` table will look like:

```
| Column          | Type                          | Collation | Nullable | Default                             |
| --------------- | ----------------------------- | --------- | -------- | ----------------------------------- |
| `id`            | `bigint`                      |           | Not null | `nextval('users_id_seq'::regclass)` |
| `created_at`    | `timestamp(0) with time zone` |           | Not null | `now()`                             |
| `name`          | `text`                        |           | Not null |                                     |
| `email`         | `citext`                      |           | Not null |                                     |
| `password_hash` | `bytea`                       |           | Not null |                                     |
| `activated`     | `boolean`                     |           | Not null |                                     |
| `version`       | `integer`                     |           | Not null | `1`                                 |
```

Indexes:

1. `users_pkey` PRIMARY KEY, btree (id)
2. `users_email_key` UNIQUE CONSTRAINT, btree (email)

> [!IMPORTANT]
> When storing a password as plain text, we use `*string` pointer type instead of `string` type, as the former offers `nil` value which is clearer than the ambiguous empty value of `string`

> [!TIP]
> The hashed password must be complex enough but not too slow to be hashed. We need to strike the balance here

Our emails are case-insensitive, so we should always store the email address using the exact casing, and **we should send people emails using that exact casing too**

## User enumeration

If an attacker wants to know if a specific email address exists, all they need is to send a request with that address and the error will return.

There are risks of leaking this information:

- User privacy
- Searching the email address in leaked password tables

Two ways to deal with this:

- [>] 1. Ensure the response to the client is always exactly the same, irrespective of whether the user exists or not
- [>] 2. Ensure the **time** taken to send the response is _always the same_, irrespective of whether a user exists or not

> [!NOTE]
> I might have gone overkill with this one, but this is worth the learning effort I guess

For the 2nd approach, we will use **a background goroutine to handle the registration process**. Doing this allows us to:

1. Timing attack prevention: A background goroutine mask the operation time through the use of channels
2. Timeout control
3. Independent timing control: The main goroutine handles timing while the background one handles the operation
4. Better resource management: A background goroutine can handle multiple registrations concurrently, and the garbage collector can collect timed-out registrations
5. Future extensibility: We can integrate a rate limiter to the `consistentTimeHandler`
6. Easier testing + Cleaner error handling
