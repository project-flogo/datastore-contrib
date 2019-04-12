# MongoDb
This activity allows you to Get, Insert, Update and Delete a document in MongoDb database.

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/activity/mongodb
```

## Configuration

### Settings:
| Name     | Type   | Description
|:---      | :---   | :---    
| Uri      | string | 

### Input: 

| Name       | Type   | Description
| :---       | :---   | :---
| DbName     | string |     
| Collection | string |     
| Method     | string |     
| KeyName    | string |     
| KeyValue   | string |     
| Data       | object | 

## Example
The below example allows you to configure the activity to reply and set the output values to literals "name" (a string) and 2 (an integer).

```json
{
  "id": "flogo-mongodb",
  "name": "MongoDb",
  "description": "MongoDb Activity",
  "activity": {
    "ref": "github.com/project-flogo/datastore-contrib/activity/mongodb",
    "settings": {
      "uri" : "localhost:27017"
    },
    "input" : {
        "dbname" : "test",
        "collection" : "example",
        "method" : "INSERT",
        "keyname" : "foo",
        "keyvalue" : "bar"
    }
  }
}
```