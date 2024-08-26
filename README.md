# go-weather
Simple http server with an http interface to fetch the current weather

# API
## GET /forecast

Fetches the current weather incuding the temperature, short foreast
and feels like temperature, `Cold`, `Hot`, or `Moderate`.
 | field     | type  | description         | required |
 |-----------|-------|---------------------|-----|
| latitude  | float | latitude of the location  | true|
| longitude | float | longitude of the location | true | 

### Example
```azure
curl -s -G http://localhost:8080/forecast -d latitude=29.760427 -d longitude=-95.369804 | jq .
{
  "short-term-forecast": "Mostly Cloudy",
  "temperature": "88F",
  "temperature-feels-like": "Hot"
}

```

## TODO
- Add robust logging
- Make bad request response formatted as JSON
- End to end testing
- Add caching layer to weather fetcher
