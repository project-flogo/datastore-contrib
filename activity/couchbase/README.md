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
| userName          | string | Cluster username    
| password          | string | Cluster password    
| bucketName        | string | The bucket name    
| bucketPassword    | string | The bucket password if any   
| server            | string | The Couchbase server (e.g. couchbase://127.0.0.1)    


### Input: 

| Name       | Type   | Description
| :---       | :---   | :---
| key        | string | The document key identifier    
| data       | string | The document data (when the method is get this field is ignored)    
| method     | string | The method type (Insert, Upsert, Remove or Get)    
| expiry     | int32  | The document expiry (default: 0)    

### Output:

| Name       | Type   | Description
| :---       | :---   | :---
| data       | object | 

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
      "userName": "Administator",
      "password": "password",
      "bucketName" : "sample",
      "bucketPassword" : "",
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