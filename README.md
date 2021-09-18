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

```json lines
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

### GET /v1/companies/search_by_category
#### Parameters

- category
- date_from
- date_to
- limit

#### Response

```json lines
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


### POST /v1/company
#### Body

```json lines
{
  "name": <string>,
  "building_id": <string>,
  "address": <string>,
  "phone_numbers": [<string>, ...],
  "categories": [<string>, ...]
}
```

#### Response

```json lines
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
