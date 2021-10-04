## API

### GET /v1/companies
#### Parameters

- id
- building_id
- category
- date_from
- date_to
- limit

#### Response

```json
{
  "error": null,
  "warning": null,
  "data": [
    {
      "id": <string>,
      "name": <string>,
      "create_at": <int64(timestamp)>,
      "building_id": <string>,
      "phone_numbers": [<string>, ...],
      "categories": [<string>, ...],
    }, ...
  ]
}
```

#### Example

`GET http://api-ugc.2gis.ru/v1/companies?category=food milk`

### POST /v1/company
#### Body

```json
{
  "name": <string>,
  "building_id": <string>,
  "address": <string>,
  "phone_numbers": [<string>, ...],
  "categories": [<string>, ...]
}
```

#### Response

```json
{
  "error": null,
  "warning": null,
  "data": [
    {
      "id": <string>,
      "name": <string>,
      "create_at": <int64(timestamp)>,
      "building_id": <string>,
      "phone_numbers": [<string>, ...],
      "categories": [<string>, ...],
    }, ...
  ]
}
```
