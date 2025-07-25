# API Documentation Template

## Overview
Brief description of the API endpoint.

## Endpoint
```
METHOD /path/to/endpoint
```

## Parameters
| Name | Type | In | Required | Description |
|------|------|----|----------|-------------|
| param1 | string | query/body/path | Yes | Description of param1 |
| param2 | integer | query/body/path | No | Description of param2 |

## Request Example
```json
{
  "param1": "value1",
  "param2": 123
}
```

## Response
### Success (200)
```json
{
  "result": "success"
}
```

### Error Responses
| Status Code | Description |
|-------------|-------------|
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing or invalid authentication |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

## Authentication
Description of authentication requirements, if any.

## Authorization
Description of authorization requirements, if any.

## Rate Limiting
Description of rate limiting, if any.

## Examples
### cURL
```bash
curl -X METHOD \
  'http://localhost:8080/path/to/endpoint' \
  -H 'Content-Type: application/json' \
  -d '{
  "param1": "value1",
  "param2": 123
}'
```

### Go
```go
// Example Go code to call this endpoint
