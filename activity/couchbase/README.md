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
| Username          | string | Cluster username    
| Password          | string | Cluster password    
| BucketName        | string | The bucket name    
| BucketPassword    | string | The bucket password if any   
| Server            | string | The Couchbase server (e.g. couchbase://127.0.0.1)    


### Input: 

| Name       | Type   | Description
| :---       | :---   | :---
| Key        | string | The document key identifier    
| Data       | string | The document data (when the method is get this field is ignored)    
| Method     | string | The method type (Insert, Upsert, Remove or Get)    
| Expiry     | int32  | The document expiry (default: 0)    

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