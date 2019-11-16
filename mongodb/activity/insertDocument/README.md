<!--
title: MongoDB Insert Document
weight: 4622
-->
# MongoDB Insert Document
This activity allows you to Insert one or more documents into a collection in a MongoDb database

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/mongodb/activity/insertDocument
```

## Configuration

### Settings:
| Name                   | Type       | Description
| :---                   | :---       | :---    
| connection             | connection | Choose a MongoDB connection from the drop down  - ***REQUIRED***
| operation              | string     | The Insert operation type (Insert One Document or Insert Many Documents) - ***REQUIRED***
| databaseName           | string     | MongoDB databse to Insert - ***REQUIRED***
| collectionName         | string     | The collection within the MongoDB database to Insert - ***REQUIRED***  
| timeout                | int32      | Timeout in seconds for the activity's operations
| continueOnErr          | bool       | For Insert Many Documents operation, should the activity continue to insert documents when the previous insertDocument failed?

### Input: 

| Name               | Type   | Description
| :---               | :---   | :---  
| data               | object | The The data to be inserted. In case of Insert Many Documents, pass an array of JSONs


### Output: 

| Name         | Type   | Description
| :---         | :---   | :---
| insertedId   | string | InsertedId of inserted document. In case of Insert Many Documents, a list of IDs is returned
| totalCount   | int    | Applicable for Insert Many Documents only. The total numner of Documents that were attempted for Insert.
| successCount | int    | Applicable for Insert Many Documents only. The total number of successful Document Insertions.
| failureCount | int    | Applicable for Insert Many Documents only. The total number of Document insertions that failed.

## Example
The below example shows the values for MongoDB Insert Document activity in the Flogo app JSON as well as the connection object referenced by the activity. The MongoDB Insert Document activity in this example attempts to insert many documents with continueOnError property set to true

```json
{
  "id": "MongoDBInsertDocument",
  "name": "MongoDBInsertDocument",
  "description": "Mongodb Insert Document activity",
  "activity": {
    "ref": "#insertDocument",
    "settings": {
      "connection": "conn://a7730ae0-0199-11ea-9e1b-1b6d6afda988",
      "operation": "Insert Many Documents",
      "databaseName": "sample",
      "collectionName": "test",
      "timeout": 23,
      "continueOnErr": true
    },
    "input": {
      "data": [
        {
        "name": "Mongo User1",
        "age": "22",
        "location": "Palo Alto"
        },
        {
        "name": "Mongo User2",
        "age": "33",
        "location": "Seattle"
        }
      ]
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