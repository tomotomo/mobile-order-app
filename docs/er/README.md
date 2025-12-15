# Database Schema

## Tables

### `system_admins`
| Column | Type | Notes |
|---|---|---|
| `id` | uint | PK |
| `name` | string | |
| `email` | string | Unique, Indexed |
| `password_hash` | string | |

### `restaurant_staffs`
| Column | Type | Notes |
|---|---|---|
| `id` | uint | PK |
| `restaurant_id` | uint | FK -> `restaurants.id` (NOT NULL) |
| `name` | string | |
| `email` | string | Unique, Indexed |
| `password_hash` | string | |
| `role` | string | `manager` or `staff` |

### `restaurants`
| Column | Type | Notes |
|---|---|---|
| `id` | uint | PK |
| `name` | string | |
| `current_guest_session` | string | Session ID for Locking |
| `guest_session_expires` | datetime | Lock expiration |

### `menu_items`
| Column | Type | Notes |
|---|---|---|
| `id` | uint | PK |
| `restaurant_id` | uint | FK -> `restaurants.id` |
| `name` | string | |
| `price` | int | |
| `stock` | int | |
| `is_sold_out` | boolean | |

## Relationships
- **Restaurant** has many **RestaurantStaffs**.
- **Restaurant** has many **MenuItems**.
- **SystemAdmin** is standalone (manages System).
- **RestaurantStaff** belongs to **Restaurant**.
