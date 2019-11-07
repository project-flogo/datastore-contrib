<!--
title: MongoDB Query Document
weight: 4622
-->
# MongoDB Query Document
This activity allows you to query for one or more documents from a collection in MongoDb database based on a matching criteria entered as a JSON

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/mongodb/activity/queryDocument
```

## Configuration

### Settings:
| Name                   | Type       | Description
| :---                   | :---       | :---    
| connection             | connection | Choose a MongoDB connection from the drop down  - ***REQUIRED***
| operation              | string     | The Query operation type (Find One Document or Find Many Documents) - ***REQUIRED***
| databaseName           | string     | MongoDB databse to query - ***REQUIRED***
| collectionName         | string     | The collection within the MongoDB database to query - ***REQUIRED***  
| timeout                | int32      | Timeout in seconds for the activity's operations

### Input: 

| Name               | Type   | Description
| :---               | :---   | :---  
| criteria           | object | The matching criteria entered as a JSON


### Output: 

| Name   | Type | Description
| :---   | :--- | :---
| output | any  | Returns a single JSON object for Find One Document operation or returns an array of JSON objects in case of Find Many Documents operation

## Example
The below example shows the values for MongoDB Query Document activity in the Flogo app JSON as well as the connection object referenced by the activity. The MongoDB Query Document activity in this example searches for and returns all the documents in a collection that has a key called 'location' with value' Palo Alto'

```json
{
  "id": "MongoDBQueryDocument",
  "name": "MongoDBQueryDocument",
  "description": "Mongodb Query Document activity",
  "activity": {
    "ref": "github.com/project-flogo/datastore-contrib/mongodb/activity/queryDocument",
    "settings": {
      "connection": "conn://a7730ae0-0199-11ea-9e1b-1b6d6afda988",
      "operation": "Find One Document",
      "databaseName": "sample",
      "collectionName": "test",
      "timeout": 0
    },
    "input": {
      "criteria": {
        "location": "Palo Alto"
      }
    }
  }
}

"connections": {
  "a7730ae0-0199-11ea-9e1b-1b6d6afda988": {
    "id": "a7730ae0-0199-11ea-9e1b-1b6d6afda988",
    "name": "mc2",
    "ref": "github.com/project-flogo/datastore-contrib/mongodb/connection",
    "settings": {
      "name": "mc2",
      "description": "",
      "connectionURI": "<Connection URI Here>",
      "credType": "<One of None or SCRAM-SHA-1 or SCRAM-SHA-256>",
      "username": "<Enter Username here in case of SCRAM-SHA-1 or SCRAM-SHA-256>",
      "password": "",
      "ssl": false,
      "x509": false,
      "trustCert": "",
      "clientCert": "",
      "clientKey": "",
      "keyPassword": "",
      }
  }
}
```