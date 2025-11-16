# Server-Side Pagination Plan

## Current State (as of Nov 16, 2025)

### What Works Now
- **ListUsers**: Returns up to 200 users with `LIMIT 200`
- **MonitorStudents**: Returns up to 200 students with `LIMIT 200`
- **Client-side pagination**: Frontend paginates results (10 items per page)
- **Fixed Issues**:
  - ✅ Removed LEFT JOIN with profile_submissions (caused row multiplication)
  - ✅ Used subqueries for profile data to avoid duplicates
  - ✅ Added COALESCE for NULL email/username fields

### Current Limitations
- Hard limit of 200 users per endpoint
- All data loaded into memory on frontend
- No filtering/sorting on server side for CreateStudents page
- MonitorStudents has server-side filtering but no pagination

---

## When to Implement Server-Side Pagination

**Trigger Point**: When active student count exceeds **100 students**

**Why 100?**
- Current limit (200) provides 2x buffer
- Performance impact becomes noticeable with 100+ rows of enriched data
- Gives time to implement before hitting hard limit

---

## Implementation Plan

### Phase 1: Add Pagination to ListUsers (Priority: HIGH)

**Endpoint**: `GET /api/admin/users`

**Changes**:
```go
// Add query parameters
page := c.DefaultQuery("page", "1")
limit := c.DefaultQuery("limit", "50")
offset := (page - 1) * limit

// Modify query
query := base + where + " ORDER BY last_name LIMIT $X OFFSET $Y"

// Return metadata
type ListUsersResponse struct {
    Data       []listUsersResp `json:"data"`
    Total      int             `json:"total"`
    Page       int             `json:"page"`
    Limit      int             `json:"limit"`
    TotalPages int             `json:"total_pages"`
}
```

**Frontend Changes** (CreateStudents.tsx):
```typescript
const [page, setPage] = useState(1);
const PAGE_SIZE = 50;

const { data } = useQuery({
  queryKey: ["admin", "users", page],
  queryFn: () => api(`/admin/users?page=${page}&limit=${PAGE_SIZE}`),
});

// data.data contains users
// data.total contains total count
// data.total_pages contains page count
```

**Benefits**:
- Reduces initial load time
- Scalable to 1000+ students
- Server-side filtering already in place (role, q)

---

### Phase 2: Optimize MonitorStudents (Priority: MEDIUM)

**Current State**:
- Already has filtering (program, department, cohort, advisor_id)
- LIMIT 200 hard-coded
- Returns enriched data (progress, advisors, deadlines)

**Changes**:
```go
// Add pagination parameters
page := c.DefaultQuery("page", "1")
limit := c.DefaultQuery("limit", "50")

// Keep existing filters
// Add OFFSET for pagination

// Return with metadata (same structure as ListUsers)
```

**Benefits**:
- Handles large student populations
- Filtering + pagination = fast queries
- Enriched data doesn't overwhelm frontend

---

### Phase 3: Database Indexing (Priority: MEDIUM)

**Current Indexes** (verify with `\d users` in psql):
- Primary key on `id`
- Probably index on `email` (for login)

**Add Indexes**:
```sql
-- For sorting by name
CREATE INDEX idx_users_last_name ON users(last_name, first_name) WHERE is_active = true;

-- For role filtering
CREATE INDEX idx_users_role_active ON users(role, is_active);

-- For MonitorStudents queries
CREATE INDEX idx_profile_submissions_user_id ON profile_submissions(user_id, created_at DESC);
```

**Benefits**:
- Faster ORDER BY last_name
- Faster role filtering
- Subqueries for profile data execute faster

---

### Phase 4: Caching Strategy (Priority: LOW)

**When Needed**: 500+ students with frequent reads

**Options**:
1. **Redis cache** for user lists (TTL: 5 minutes)
2. **Query result caching** in Go (memory)
3. **Frontend caching** with React Query (already implemented)

**Implementation**:
```go
// Example: Cache total count
cacheKey := fmt.Sprintf("users:count:role=%s", roleFilter)
if cached, found := h.cache.Get(cacheKey); found {
    total = cached.(int)
} else {
    // Query DB
    h.cache.Set(cacheKey, total, 5*time.Minute)
}
```

---

## Migration Checklist

**Before implementing pagination:**
- [ ] Current student count > 100
- [ ] Performance metrics show slow load times (>2s)
- [ ] User feedback about slow page loads

**Implementation steps:**
1. [ ] Add pagination to ListUsers backend
2. [ ] Update CreateStudents frontend to use pagination
3. [ ] Test with 200+ mock students
4. [ ] Add pagination to MonitorStudents
5. [ ] Create database indexes
6. [ ] Monitor query performance
7. [ ] Implement caching if needed

**Rollback plan:**
- Keep `LIMIT 200` fallback if pagination breaks
- Feature flag for pagination (enable/disable)
- Monitor error rates in production

---

## Performance Targets

| Metric | Current | Target (with pagination) |
|--------|---------|--------------------------|
| Initial load time | ~500ms (38 users) | <1s (1000+ users) |
| Memory usage (frontend) | ~50KB | <100KB per page |
| Query time (backend) | ~25ms | <50ms per page |
| Time to first render | ~200ms | <300ms |

---

## Testing Strategy

**Load Testing**:
```bash
# Create 500 mock students
go run ./cmd/mock --students=500

# Test pagination endpoint
curl "http://localhost:8280/api/admin/users?page=1&limit=50"
curl "http://localhost:8280/api/admin/users?page=10&limit=50"

# Measure query time
psql -c "EXPLAIN ANALYZE SELECT ... LIMIT 50 OFFSET 450"
```

**Frontend Testing**:
- Navigate through all pages
- Test search with pagination
- Test sorting with pagination
- Verify total count accuracy

---

## Notes

- Current implementation (LIMIT 200) is sufficient for small-medium universities
- Pagination adds complexity - only implement when needed
- Monitor student growth rate to plan ahead
- Consider batch operations for 1000+ students (CSV export, bulk updates)
