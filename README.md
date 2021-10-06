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

### POST /v1/categories

Запрос для создания новой категории.
В теле запрос нужно передать `json` который описан ниже.
В ответ вернётся только что созданный объект `category`.

#### Body

```json
{
  "name": "string"
}
```

| Название поля  | Описание |
| ------------- | ------------- |
| `name`  | Название категории. |

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

Запрос на получение списка категорий.
Параметры используемые в запросе перечислены ниже.

**TODO: приоритеты**

#### Parameters

| Параметр  | Описание |
| ------------- | ------------- |
| `id`  | Уникальный идентификатор компании. |
| `search_by_name`  | Поиск категорий по названию. Их можно перечислять через пробел либо через точку. Нет чувствительности к регистру. |
| `from_date` | Поиск по полю `create_at` с условием `create_at >= from_date`. |
| `to_date` | Поиск по полю `create_at` с условием `create_at <= from_date`. |
| `offset` | Смещение результата на указанное количество объектов. Лучше использовать это поле вместе с `ascending`. Иначе повторяемость результата не гарантируется.  |
| `ascending` | Сортировка результата по полю `create_at`. Значение `true` сортирует по возрастанию, `false` по убыванию. Если значение не указано, результат возвращается с не определённой сортировкой. |
| `limit` | Ограничение количества объектов. |

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

Запрос для создания нового здания. 
В теле запрос нужно передать `json` который описан ниже. 
В ответ вернётся только что созданный объект `building`. 

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

| Название поля  | Описание |
| ------------- | ------------- |
| `address`  | Адрес нового здания. |
| `location` | Геогарфические координаты здания. |

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

Запрос на получение списка зданий. 
Параметры используемые в запросе перечислены ниже.

**TODO: приоритеты**

#### Parameters

| Параметр  | Описание |
| ------------- | ------------- |
| `id`  | Уникальный идентификатор компании. |
| `address`  | Поиск здания по адресу. |
| `from_date` | Поиск по полю `create_at` с условием `create_at >= from_date`. |
| `to_date` | Поиск по полю `create_at` с условием `create_at <= from_date`. |
| `offset` | Смещение результата на указанное количество объектов. Лучше использовать это поле вместе с `ascending`. Иначе повторяемость результата не гарантируется.  |
| `ascending` | Сортировка результата по полю `create_at`. Значение `true` сортирует по возрастанию, `false` по убыванию. Если значение не указано, результат возвращается с не определённой сортировкой. |
| `limit` | Ограничение количества объектов. |

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

Запрос для создания компании. 
В теле запроса нужно передать `json` который описан ниже.
В ответ вернётся только что созданный объект `company`.

#### Body

```json
{
  "name": "string",
  "building_id": "string",
  "phone_numbers": ["string", "..."],
  "category_ids": ["string", "..."]
}
```

| Название поля  | Описание |
| ------------- | ------------- |
| `name`  | Название для новой компании. |
| `building_id`  | Идентификатор здания в котором будет находиться комапания. |
| `phone_numbers` | Список номеров телефонов которые относятся к компании. |
| `category_ids` | Список идентификаторов категорий к которым относится компания. |

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

Запрос на получение списка компаний. 
Параметры используемые в запросе перечислены ниже.

**TODO: приоритеты**

#### Parameters

| Параметр  | Описание |
| ------------- | ------------- |
| `id`  | Уникальный идентификатор компании. |
| `building_id`  | Идентификатор здания в котором находится комапания. |
| `search_by_category` | Поиск компании по названию категории. Категории можно перечислять через пробел либо через точку. Нет чувствительности к регистру. |
| `from_date` | Поиск по полю `create_at` с условием `create_at >= from_date`. |
| `to_date` | Поиск по полю `create_at` с условием `create_at <= from_date`. |
| `offset` | Смещение результата на указанное количество объектов. Лучше использовать это поле вместе с `ascending`. Иначе повторяемость результата не гарантируется.  |
| `ascending` | Сортировка результата по полю `create_at`. Значение `true` сортирует по возрастанию, `false` по убыванию. Если значение не указано, результат возвращается с не определённой сортировкой. |
| `limit` | Ограничение количества объектов. |

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

