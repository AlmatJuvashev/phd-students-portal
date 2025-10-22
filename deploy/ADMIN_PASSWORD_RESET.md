# Admin Password Reset Guide

## Overview

Administrators can reset passwords for any user (except superadmin) from the user management interface.

## API Endpoint

```http
POST /api/admin/users/:id/reset-password
Authorization: Bearer <admin_jwt_token>
```

### Request

No body required. The endpoint automatically generates a new secure temporary password.

### Response

```json
{
  "username": "john.doe",
  "temp_password": "Abc123XyZ"
}
```

### Security Features

1. **Admin/Superadmin Only**: Requires admin or superadmin role
2. **Cannot Reset Superadmin**: Superadmin passwords cannot be reset by admins
3. **Auto-Generated Password**: Secure random password (12+ chars, mixed case, numbers)
4. **Returns Credentials**: Admin receives username + temp password to share with user

## Usage Flow

### From Frontend Admin Panel

1. Admin navigates to User Management
2. Finds user in the list
3. Clicks "Reset Password" button
4. System displays dialog with new credentials:
   ```
   Username: john.doe
   Temporary Password: Abc123XyZ
   
   Share these credentials with the user securely.
   ```
5. Admin copies credentials and shares with user via secure channel

### User First Login

1. User logs in with temporary password
2. System prompts to change password immediately
3. User sets new permanent password

## Frontend Integration

### Admin Component Example

```typescript
async function handleResetPassword(userId: string) {
  try {
    const response = await api.post(`/admin/users/${userId}/reset-password`);
    const { username, temp_password } = response.data;
    
    // Show modal with credentials
    alert(`New credentials:\nUsername: ${username}\nPassword: ${temp_password}`);
    
    // Or copy to clipboard
    await navigator.clipboard.writeText(
      `Username: ${username}\nPassword: ${temp_password}`
    );
  } catch (error) {
    console.error('Failed to reset password:', error);
  }
}
```

### User List Table Example

```tsx
<Table>
  <TableBody>
    {users.map(user => (
      <TableRow key={user.id}>
        <TableCell>{user.name}</TableCell>
        <TableCell>{user.email}</TableCell>
        <TableCell>{user.role}</TableCell>
        <TableCell>
          {user.role !== 'superadmin' && (
            <Button 
              onClick={() => handleResetPassword(user.id)}
              variant="outline"
              size="sm"
            >
              Reset Password
            </Button>
          )}
        </TableCell>
      </TableRow>
    ))}
  </TableBody>
</Table>
```

## Database Changes

The password reset updates:

```sql
UPDATE users 
SET password_hash = $1, updated_at = NOW() 
WHERE id = $2
```

No audit trail is stored by default. Consider adding audit logging if needed.

## Security Considerations

### ✅ Implemented

- Role-based access control (admin/superadmin only)
- Superadmin protection (cannot be reset)
- Secure password generation
- HTTPS transport (in production)
- JWT authentication required

### ⚠️ Recommendations for Production

1. **Add Audit Logging**: Log all password resets
   ```sql
   INSERT INTO admin_actions (admin_id, action, target_user_id, created_at)
   VALUES ($1, 'password_reset', $2, NOW())
   ```

2. **Notification**: Email user when password is reset
   ```go
   mailer.Send(user.Email, "Password Reset", 
     "Your password was reset by an administrator...")
   ```

3. **Force Password Change**: Mark user to change password on next login
   ```sql
   ALTER TABLE users ADD COLUMN must_change_password BOOLEAN DEFAULT FALSE;
   ```

4. **Rate Limiting**: Prevent abuse
   ```go
   // Max 5 resets per admin per hour
   ```

## Related Endpoints

- `POST /api/admin/users` - Create new user (generates temp password)
- `PUT /api/admin/users/:id` - Update user details
- `PATCH /api/admin/users/:id/active` - Activate/deactivate user
- `PATCH /api/me/password` - User changes own password

## Error Codes

| Status | Error | Reason |
|--------|-------|--------|
| 401 | `unauthorized` | Missing or invalid JWT token |
| 403 | `cannot reset superadmin password` | Target user is superadmin |
| 404 | `user not found` | Invalid user ID |
| 500 | `update failed` | Database error |

## Testing

### Manual Test

```bash
# 1. Login as admin
curl -X POST http://localhost:8280/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Save token
TOKEN="<jwt_token_from_response>"

# 2. Reset user password
curl -X POST http://localhost:8280/api/admin/users/<user_id>/reset-password \
  -H "Authorization: Bearer $TOKEN"

# Response:
# {
#   "username": "john.doe",
#   "temp_password": "Abc123XyZ"
# }
```

### Unit Test Example

```go
func TestResetPasswordForUser(t *testing.T) {
	// Setup
	db := setupTestDB()
	handler := NewUsersHandler(db, config.AppConfig{})
	
	// Create test user
	userID := createTestUser(db, "student")
	
	// Test reset
	req := httptest.NewRequest("POST", "/admin/users/"+userID+"/reset-password", nil)
	w := httptest.NewRecorder()
	handler.ResetPasswordForUser(w, req)
	
	// Assert
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "temp_password")
}
```

## Migration History

- **0008_remove_password_reset**: Removed email-based password reset (not used)
- Password reset now admin-only via `/admin/users/:id/reset-password`

---

**Note**: This replaces the old email-based password reset flow. Users cannot reset passwords themselves - only admins can do this.
