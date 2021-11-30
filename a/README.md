### Request

```bash
curl -X GET \
  http://127.0.0.1:7070/api/order \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: aec8dd29-8ec2-42c8-b731-26cf801559f9' \
  -H 'cache-control: no-cache' \
  -d '{
    "id": 12,
    "price": 1000,
    "title": "burger"
}'
```

### Response

```json
{
  "result": "order has published successfully"
}
```