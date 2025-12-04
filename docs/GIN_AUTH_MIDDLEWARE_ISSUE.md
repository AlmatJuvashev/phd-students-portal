# Gin Auth Middleware Issue: User Context Not Available in Handlers

## Problem Description

When using Gin framework with JWT authentication, handlers may receive `401 Unauthorized` errors even when the JWT token is valid. The root cause is a subtle bug in middleware chaining where `c.Next()` is called prematurely.

### Symptoms

- API returns `{"error":"unauthorized"}` despite valid JWT token
- Logs show `hasToken: true` but handler can't access user context
- `c.Get("userID")` returns `nil` or `false` in handlers
- Works inconsistently - sometimes works after login, fails after page refresh

## Root Cause

### The Bug Pattern

```go
// ❌ BROKEN: AuthRequired middleware
func AuthRequired(db *gorm.DB, redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Validate JWT token
        claims, err := validateToken(c)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        
        // 2. Set claims in context
        c.Set("claims", claims)
        
        // 3. Call another middleware to hydrate user
        HydrateUserFromClaims(db, redis)(c)  // ⚠️ This middleware has c.Next()!
        
        // 4. Continue to handler
        c.Next()  // ⚠️ This is called AFTER handler already executed!
    }
}

// The HydrateUserFromClaims middleware
func HydrateUserFromClaims(db *gorm.DB, redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... fetch user from DB/cache ...
        c.Set("userID", user.ID)
        c.Set("current_user", &user)
        
        c.Next()  // ⚠️ BUG: This passes control to the HANDLER, not back to AuthRequired!
    }
}
```

### Why This Happens

In Gin, `c.Next()` doesn't return control to the calling function. Instead, it immediately executes the next handler in the chain. So when `HydrateUserFromClaims` calls `c.Next()`:

1. Control passes directly to the route handler
2. Handler executes **before** `c.Set("userID", ...)` completes
3. Handler tries to get `userID` from context → gets `nil`
4. Handler returns 401

### Visual Flow

```
❌ BROKEN FLOW:
AuthRequired() 
  → validates JWT ✓
  → calls HydrateUserFromClaims()
      → starts fetching user...
      → c.Next() ← JUMPS TO HANDLER IMMEDIATELY!
          → Handler executes (userID not set yet!)
          → Returns 401
      → c.Set("userID", ...) ← TOO LATE!
```

## Solution

### Option 1: Separate Validation Function (Recommended)

Create a validation function that does NOT call `c.Next()`:

```go
// ✅ FIXED: Separate validation function (no c.Next())
func validateJWT(c *gin.Context, jwtSecret string) (jwt.MapClaims, error) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        return nil, errors.New("missing or invalid authorization header")
    }
    
    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    
    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(jwtSecret), nil
    })
    
    if err != nil || !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("invalid claims")
    }
    
    return claims, nil
}

// ✅ FIXED: AuthMiddleware that orchestrates the flow
func AuthMiddleware(db *gorm.DB, redis *redis.Client, jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Validate JWT (no c.Next() called)
        claims, err := validateJWT(c, jwtSecret)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        
        // 2. Set claims
        c.Set("claims", claims)
        
        // 3. Hydrate user (inline, no c.Next())
        sub, _ := claims["sub"].(string)
        var user User
        if err := db.Where("id = ?", sub).First(&user).Error; err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "user not found"})
            return
        }
        
        // 4. Set user context
        c.Set("userID", user.ID)
        c.Set("current_user", &user)
        
        // 5. Verify everything is set
        if _, exists := c.Get("userID"); !exists {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        
        // 6. NOW call Next() - after all context is set
        c.Next()
    }
}
```

### Option 2: Inline User Hydration

If you want to keep middleware separate, inline the hydration logic:

```go
func AuthMiddleware(db *gorm.DB, redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Validate JWT...
        claims, err := validateToken(c)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        c.Set("claims", claims)
        
        // Inline hydration (don't call another middleware!)
        sub := claims["sub"].(string)
        
        // Try cache first
        var user User
        cached, err := redis.Get(c, "user:"+sub).Result()
        if err == nil {
            json.Unmarshal([]byte(cached), &user)
        } else {
            db.Where("id = ?", sub).First(&user)
            // Cache for next time
            data, _ := json.Marshal(user)
            redis.Set(c, "user:"+sub, data, 5*time.Minute)
        }
        
        c.Set("userID", user.ID)
        c.Set("current_user", &user)
        
        c.Next()  // Only ONE c.Next() at the very end
    }
}
```

## How to Scan Your Codebase for This Issue

### 1. Search for Middleware Calling Other Middleware

```bash
# Look for middleware functions that call other middleware
grep -rn "gin.HandlerFunc" --include="*.go" | head -20

# Look for patterns where one middleware calls another
grep -rn ")(c)" --include="*.go" | grep -v "test"
```

### 2. Check for Multiple c.Next() Calls

```bash
# Find files with multiple c.Next() that might be chained
grep -rn "c.Next()" --include="*.go" -A 2 -B 2
```

### 3. Audit Middleware Chain

Look for patterns like:

```go
// ⚠️ SUSPICIOUS: Middleware calling middleware
SomeMiddleware(deps)(c)  // If SomeMiddleware has c.Next(), this is a bug!
c.Next()
```

### 4. Add Debug Logging

Add these logs to identify the issue:

```go
// In your auth middleware
log.Printf("[AuthMiddleware] Starting for path=%s", c.Request.URL.Path)

// After setting user context
log.Printf("[AuthMiddleware] userID set: %v", c.Get("userID"))

// In your handler
log.Printf("[Handler] userID from context: %v", c.Get("userID"))
```

## Quick Checklist

- [ ] Each middleware has **exactly one** `c.Next()` call
- [ ] `c.Next()` is called **at the end** of the middleware, after all context is set
- [ ] Middlewares don't call other middlewares that have `c.Next()`
- [ ] If you need to compose middlewares, extract the logic into functions without `c.Next()`
- [ ] Add debug logging to trace the flow during development

## Example: Before and After

### Before (Broken)

```go
// api.go
api := r.Group("/api")
api.Use(AuthRequired(db, redis))  // Calls HydrateUserFromClaims internally

// auth.go
func AuthRequired(...) gin.HandlerFunc {
    return func(c *gin.Context) {
        validateJWT(c)
        HydrateUserFromClaims(db, redis)(c)  // ❌ Has c.Next()!
        c.Next()
    }
}

func HydrateUserFromClaims(...) gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... set user ...
        c.Next()  // ❌ Jumps to handler!
    }
}
```

### After (Fixed)

```go
// api.go  
api := r.Group("/api")
api.Use(AuthMiddleware(db, redis, cfg.JWTSecret))

// auth.go
func validateJWT(c *gin.Context, secret string) (jwt.MapClaims, error) {
    // Pure function, no c.Next()
    return claims, nil
}

func hydrateUser(c *gin.Context, db *gorm.DB, redis *redis.Client, sub string) error {
    // Pure function, no c.Next()
    // Sets c.Set("userID", ...) and c.Set("current_user", ...)
    return nil
}

func AuthMiddleware(...) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, err := validateJWT(c, secret)  // ✅ No c.Next()
        if err != nil { c.AbortWithStatusJSON(401, ...); return }
        
        c.Set("claims", claims)
        
        if err := hydrateUser(c, db, redis, claims["sub"].(string)); err != nil {
            c.AbortWithStatusJSON(401, ...); return
        }
        
        c.Next()  // ✅ Only one c.Next(), at the very end
    }
}

// DEPRECATED - do not use
func AuthRequired(...) gin.HandlerFunc { ... }
```

## Related Issues

- Gin issue: Handler executes before middleware completes
- JWT claims available but user context missing
- Intermittent 401 errors with valid tokens
- Context values not persisting across middleware chain

## References

- [Gin Middleware Documentation](https://gin-gonic.com/docs/examples/custom-middleware/)
- [Understanding c.Next() in Gin](https://gin-gonic.com/docs/examples/using-middleware/)
