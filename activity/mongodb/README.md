<!--
title: MongoDB
weight: 4622
-->
# MongoDb
This activity allows you to Get, Insert, Update and Delete a document in MongoDb database.

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/activity/mongodb
```

## Configuration

### Settings:
| Name       | Type   | Description
|:---        | :---   | :---    
| uri        | string | The MongoDB connection URI  - ***REQUIRED***
| method     | string | The method type (GET, INSERT, UPDATE or DELETE) - ***REQUIRED***
| dbName     | string | The name of the database - ***REQUIRED***
| collection | string | The collection to work on - ***REQUIRED***  
| username   | string | The username of the client
| password   | string | The password of the client   

### Input: 

| Name       | Type   | Description
| :---       | :---   | :---
| keyName    | string | The name of the key to use when looking up an object (used in GET, UPDATE and DELETE)    
| keyValue   | string | The value of the key to use when looking up an object (used in GET, UPDATE, and DELETE)    
| data       | object | The bson document to insert in mongodb


### Output: 

| Name   | Type | Description
| :---   | :--- | :---
| output | any  | The result of the MongoDB method invocation

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
      "uri" : "localhost:27017",
      "method" : "INSERT",
      "dbName" : "test",
      "collection" : "example"
    },
    "input" : {
        "keyName" : "foo",
        "keyValue" : "bar"
    }
  }
}
```