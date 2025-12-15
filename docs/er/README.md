# Database Schema

## Tables

### `users`
| Column | Type | Notes |
|---|---|---|
| `id` | uint | PK |
| `name` | string | |
| `email` | string | Unique, Indexed |
| `password_hash` | string | |
| `role` | string | `admin`, `manager`, `staff` |
| `restaurant_id` | uint | FK -> `restaurants.id` (Nullable for admin) |

### `restaurants`
| Column | Type | Notes |
|---|---|---|
| `id` | uint | PK |
| `name` | string | |
| `user_id` | uint | ID of the Creator/Manager |
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
- **Restaurant** has many **Users** (Manager, Staff).
- **Restaurant** has many **MenuItems**.
- **User** belongs to **Restaurant** (if not admin).
