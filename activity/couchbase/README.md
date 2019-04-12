# Couchbase
This activity allows you to Get, Insert, Update and Delete a document in couchbase database.

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/activity/couchbase
```

## Configuration

### Settings:
| Name              | Type   | Description
| :---              | :---   | :---
| Username          | string |     
| Password          | string |     
| BucketName        | string |     
| BucketPassword    | string |     
| Server            | string |     


### Input: 

| Name       | Type   | Description
| :---       | :---   | :---
| Key        | string |     
| Data       | string |     
| Method     | string |     
| Expiry     | int32  |     

### Output:

| Name       | Type   | Description
| :---       | :---   | :---
| Data       | object | 

## Example
The below example allows you to configure the activity to reply and set the output values to literals "name" (a string) and 2 (an integer).

```json
{
  "id": "flogo-mongodb",
  "name": "MongoDb",
  "description": "MongoDb Activity",
  "activity": {
    "ref": "github.com/project-flogo/datastore-contrib/activity/couchbase",
    "settings": {
      "server" : "http://localhost:8091",
      "username": "Administator",
      "password": "password",
      "bucketname" : "sample",
      "bucketpassword" : "",
    },
    "input" : {
        "key" : "test",
        "data" : "example",
        "method" : "Insert",
        "expiry" : 0
    }
  }
}
```