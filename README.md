## Entities

### Building

| Field | Type | Description |
| ------------- | --- | ------------- |
| `id` | `string(uuid)` | Unique identification. It generating on server. |
| `create_at` | `int64(timestamp)` | Date creation of the object. Timestamp using in milliseconds. It generating on server. |
| `address` | `string` | Address of the building. |
| `location` | [`location`](#location) | Geographic coordinates for building. |

### Location

| Field | Type | Description |
| ------------- | --- | ------------- |
| `lat` | `float64` | Geographic latitude. Specified by degrees, from -90° to 90°. Example `-82.91218902757737`. |
| `lng` | `float64` | Geographic longitude. Specified by degrees, from -180° to 180°. Example `-175.83667565162875`. |

### Category

| Field | Type | Description |
| ------------- | --- | ------------- |
| `id` | `string(uuid)` | Unique identification. It generating on server. |
| `name` | `string(uuid)` | Name of the category. It consist from words separate of dots. Full name starting from root element and finish on top. Example: `Transport.Petrol.Cars`. |
| `create_at` | `int64(timestamp)` | Date creation of the object. Timestamp using in milliseconds. It generating on server. |

### Company

| Field | Type | Description |
| ------------- | --- | ------------- |
| `id` | `string(uuid)` | Unique identification. It generating on server. |
| `name` | `string` | Name of the company. |
| `create_at` | `int64(timestamp)` | Date creation of the object. Timestamp using in milliseconds. It generating on server. |
| `building_id` | `string(uuid)` | Unique building identification for company. |
| `address` | `string` | Address of the company. |
| `phone_numbers` | `[]string` | List of company phone numbers. |
| `categories` | `[]string` | List of company categories. |

## HTTP API

### Overview

### Errors

### Warnings

### POST /v1/categories

Request for create new category.

#### Body

```json
{
  "name": "string"
}
```

| Field  | Description |
| ------------- | ------------- |
| `name` | Name of category. More [here](#category). |

#### Response

```json
{
  "error": null,
  "warning": null,
  "data":
  {
    "id": "string",
    "create_at": "int64(timestamp)",
    "name": "string"
  }
}
```


### GET /v1/categories

Request for get categories.

#### Parameters

Each parameter have priority. It means that if there is field with high priority, then fields with priority below will be ignored.
Usually very high priority have `id` field.

| Field  | Description |
| ------------- | ------------- |
| `id`  | Unique category identification. High priority. |
| `search_by_name` | Search categories by name. They can be listed via dot or space. Insensitive case. Priority below than `id` field.|
| `from_date` | Search by `create_at` field with condition `create_at >= from_date`. |
| `to_date` | Search by `create_at` field with condition `create_at <= from_date`. |
| `offset` | Offset of result. Better use this field with `ascending` field. Otherwise, repeatability of the result is not guaranteed. |
| `ascending` | Sorting of result by `create_at` field. Value is `true` sorting by ascending, and `false` by descending. If value is not define, then result returning with undefined sorting. |
| `limit` | Limit of amount objects in result. |

#### Response

```json
{
  "error": null,
  "warning": null,
  "data": [
    {
      "id": "string",
      "create_at": "int64(timestamp)",
      "name": "string"
    }
  ]
}
```

### POST /v1/buildings

Request for create new building.

#### Body

```json
{
  "address": "string",
  "location": {
    "lat": "float64", 
    "lng": "float64"
  }
}
```

| Field | Description |
| ------------- | ------------- |
| `address`  | New building address. |
| `location` | Geographic coordinates for building. More [here](#location). |

#### Response

```json
{
  "error": null,
  "warning": null,
  "data": 
    {
      "id": "string",
      "create_at": "int64(timestamp)",
      "address": "string",
      "location": {
        "lat": "float64",
        "lng": "float64"
      }
    }
}
```

### GET /v1/buildings

Request for get buildings.

#### Parameters

Each parameter have priority. It means that if there is field with high priority, then fields with priority below will be ignored.
Usually very high priority have `id` field.


| Field  | Description |
| ------------- | ------------- |
| `id`  | Unique building identification. High priority.|
| `address`  | Search by address. Priority below than `id` field.|
| `from_date` | Search by `create_at` field with condition `create_at >= from_date`. |
| `to_date` | Search by `create_at` field with condition `create_at <= from_date`. |
| `offset` | Offset of result. Better use this field with `ascending` field. Otherwise, repeatability of the result is not guaranteed. |
| `ascending` | Sorting of result by `create_at` field. Value is `true` sorting by ascending, and `false` by descending. If value is not define, then result returning with undefined sorting. |
| `limit` | Limit of amount objects in result. |

#### Response

```json
{
  "error": null,
  "warning": null,
  "data": [
    {
      "id": "string",
      "create_at": "int64(timestamp)",
      "address": "string",
      "location": {
        "lat": "float64",
        "lng": "float64"
      }
    }
  ]
}
```

### POST /v1/companies

Request for create company.

#### Body

```json
{
  "name": "string",
  "building_id": "string",
  "phone_numbers": ["string", "..."],
  "category_ids": ["string", "..."]
}
```

| Field  | Description |
| ----------| ------------- |
| `name` | New company name. |
| `building_id`  | Unique building identification in which the company is located. |
| `phone_numbers` | List of company phone numbers. |
| `category_ids` | List of category ids for the company. |

#### Response

```json
{
  "error": null,
  "warning": null,
  "data": 
    {
      "id": "string",
      "name": "string",
      "create_at": "int64(timestamp)",
      "building_id": "string",
      "phone_numbers": ["string", "..."],
      "categories": ["string", "..."]
    }
}
```

### GET /v1/companies

Request for get companies.

#### Parameters

Each parameter have priority. It means that if there is field with high priority, then fields with priority below will be ignored.
Usually very high priority have `id` field.


| Field  | Description |
| ------------- | ------------- |
| `id`  | Unique company identification. High priority.|
| `building_id`  | Unique building identification in which the company is located. Priority below than `id`, `search_by_category` fields.|
| `search_by_category` | Search by category name. They can be listed via dot or space. Insensitive case. Priority below than `id` field.|
| `from_date` | Search by `create_at` field with condition `create_at >= from_date`. |
| `to_date` | Search by `create_at` field with condition `create_at <= from_date`. |
| `offset` | Offset of result. Better use this field with `ascending` field. Otherwise, repeatability of the result is not guaranteed. |
| `ascending` | Sorting of result by `create_at` field. Value is `true` sorting by ascending, and `false` by descending. If value is not define, then result returning with undefined sorting. |
| `limit` | Limit of amount objects in result. |


#### Response

```json
{
  "error": null,
  "warning": null,
  "data": [
    {
      "id": "string",
      "name": "string",
      "create_at": "int64(timestamp)",
      "building_id": "string",
      "phone_numbers": ["string", "..."],
      "categories": ["string", "..."]
    }
  ]
}
```

