# Permission-based Authorization

Only activated users can access `/v1/movies**`

Fine-grained control over which users can access which endpoints

## Relationship between permissions and users

The relationship between `permissions` and `users` is **many-to-many**: A user can have multiple permissions

```sql
-- Get the permissions corresponding to an user
SELECT permissions.code
FROM permissions
INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
INNER JOIN users ON users_permissions.user_id = users.id
WHERE users.id = $1
```
