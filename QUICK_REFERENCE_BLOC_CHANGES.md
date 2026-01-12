# Quick Reference - Backend Bloc Support Changes

## Files Modified

### 1. `attendance.sql` - Database Migration
```sql
-- Added at the end of the file:
ALTER TABLE events
ADD COLUMN IF NOT EXISTS tagged_blocs TEXT;
```

### 2. `models/event_model.go` - Model Definition
Added to Event struct:
```go
TaggedBlocsCSV string   `json:"-" gorm:"type:text;column:tagged_blocs"`
TaggedBlocs    []string `json:"tagged_blocs,omitempty" gorm:"-"`
```

Updated EventRequest struct:
```go
type EventRequest struct {
    // ... existing fields ...
    TaggedBlocs []string `json:"tagged_blocs,omitempty"`
}
```

### 3. `services/event_service.go` - Business Logic
Added functions:
- `setTaggedBlocsFromRequest()` - Normalize and store blocs
- `parseTaggedBlocsCSV()` - Convert CSV to array

Updated functions:
- `CreateEvent()` - Added bloc handling
- `UpdateEvent()` / `applyOtherUpdates()` - Added bloc handling
- `GetEvent()` - Added bloc parsing
- `GetAllEvents()` - Added bloc parsing
- `populateTaggedCoursesAndAllowed()` - Added bloc parsing

### 4. `controller/event_controller.go` - HTTP Handler
Added function:
- `parseTaggedBlocs()` - Extract blocs from request

Updated handlers:
- `CreateEvent()` - Call parseTaggedBlocs()
- `UpdateEvent()` - Call parseTaggedBlocs()

## How It Works

### Database
- Blocs stored as CSV: "Block 1,Block 2,Block 3"
- Column name: `tagged_blocs`

### API Request
```json
{
  "tagged_blocs": ["Block 1", "Block 2", "Block 3"]
}
```

### API Response
```json
{
  "tagged_blocs": ["Block 1", "Block 2", "Block 3"]
}
```

## Deployment Steps

1. Run database migration in `attendance.sql`
2. Rebuild Go binary: `go build`
3. Deploy new binary
4. Frontend can now send/receive `tagged_blocs` in event creation

## Verification

After deployment, test with:
```bash
# Create event with blocs
curl -X POST http://localhost:8000/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Event",
    "event_date": "2026-02-15",
    "start_time": "09:00",
    "end_time": "11:00",
    "tagged_blocs": ["Block 1", "Block 2"]
  }'

# Response should include:
{
  "tagged_blocs": ["Block 1", "Block 2"],
  ...
}
```

## Backward Compatibility

✅ All existing events without blocs continue to work  
✅ No breaking changes to existing endpoints  
✅ Optional field - not required for event creation  
✅ Null/empty blocs handled gracefully  

## Key Patterns Followed

Same patterns used for courses:
- Form & JSON parsing
- CSV storage in database
- Array response in API
- Normalization function
- CSV parsing function
- Updated in create/update flows
- Parsed in retrieval flows

## Summary

The backend is now ready to support the frontend's bloc selection feature. All changes follow existing patterns in the codebase for consistency and maintainability.
