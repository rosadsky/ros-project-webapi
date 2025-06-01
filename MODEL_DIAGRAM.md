# Hospital Space Management - Model Relationships

## Entity Relationship Diagram

```
┌─────────────────────────────────────┐
│              Ambulance              │
├─────────────────────────────────────┤
│ + id: UUID                          │
│ + name: string                      │
│ + location: string                  │
│ + status: string                    │
│ + type: string                      │
│ + created_at: time.Time             │
│ + updated_at: time.Time             │
└─────────────────────────────────────┘
                    │
                    │ 1:N (optional)
                    │ assignment
                    ▼
┌─────────────────────────────────────┐
│               Space                 │
├─────────────────────────────────────┤
│ + id: UUID                          │
│ + name: string                      │
│ + type: string                      │
│ + floor: int                        │
│ + capacity: int                     │
│ + status: string                    │
│ + assigned_to: *string              │ ◄── Can be ambulance name
│ + assigned_type: *string            │ ◄── "ambulance" or "department"  
│ + assigned_id: *UUID                │ ◄── References Ambulance.id or department ID
│ + created_at: time.Time             │
│ + updated_at: time.Time             │
└─────────────────────────────────────┘
                    │
                    │ 1:N (optional)
                    │ assignment
                    ▼
┌─────────────────────────────────────┐
│            Department               │
│         (Virtual Entity)            │
├─────────────────────────────────────┤
│ • Represented as string name        │
│ • No formal model structure         │
│ • Referenced by assigned_to field   │
│ • Can have UUID identifier          │
└─────────────────────────────────────┘
```

## Relationship Description

### 1. Ambulance ↔ Space (1:N Optional)
- **Type**: One-to-Many (Optional)
- **Description**: One ambulance can be assigned to multiple spaces, but each space can only be assigned to one ambulance at a time
- **Connection Fields**:
  - `Space.assigned_id` → `Ambulance.id`
  - `Space.assigned_to` → `Ambulance.name`
  - `Space.assigned_type` = "ambulance"

### 2. Department ↔ Space (1:N Optional)
- **Type**: One-to-Many (Optional)
- **Description**: One department can be assigned to multiple spaces, but each space can only be assigned to one department at a time
- **Connection Fields**:
  - `Space.assigned_to` → Department name (string)
  - `Space.assigned_type` = "department"
  - `Space.assigned_id` → Department UUID (optional)

## Status Flow

```
Space Status Flow:
available → occupied (when assigned)
occupied → available (when unassigned)
available → maintenance
maintenance → available

Ambulance Status Flow:
available → busy
busy → available
available → maintenance
maintenance → available
```

## Key Features

1. **Simple Assignment Model**: Spaces can be assigned to either ambulances or departments
2. **Flexible References**: Uses both name (string) and ID (UUID) for assignments
3. **Optional Relationships**: All assignments are optional (nullable fields)
4. **Status Tracking**: Both entities track their current status
5. **Audit Trail**: Created/Updated timestamps on all entities

## API Endpoints

### Space Management (4 Simple CRUD Operations)
- `POST /api/spaces` - CREATE new space
- `GET /api/spaces` - READ all spaces  
- `PUT /api/spaces/{id}` - UPDATE space assignment
- `DELETE /api/spaces/{id}` - DELETE space

### Ambulance Support
- `POST /api/ambulances` - Create ambulance (for assignments)
- `GET /api/ambulances` - List ambulances (for assignments) 