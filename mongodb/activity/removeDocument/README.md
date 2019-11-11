<!--
title: MongoDB Remove Document
weight: 4622
-->
# MongoDB Remove Document
This activity allows you to remove one or more documents from a collection in MongoDb database based on a matching criteria entered as a JSON

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/mongodb/activity/removeDocument
```

## Configuration

### Settings:
| Name                   | Type       | Description
| :---                   | :---       | :---    
| connection             | connection | Choose a MongoDB connection from the drop down  - ***REQUIRED***
| operation              | string     | The Query operation type (Remove One Document or Remove Many Documents) - ***REQUIRED***
| databaseName           | string     | MongoDB databse to remove documents - ***REQUIRED***
| collectionName         | string     | The collection within the MongoDB database to remove documents - ***REQUIRED***  
| timeout                | int32      | Timeout in seconds for the activity's operations

### Input: 

| Name               | Type   | Description
| :---               | :---   | :---  
| criteria           | object | The matching criteria entered as a JSON - ***In case an empty JSON is provided and if operation is Remove Many Documents, all the documents in the collection will be removed***


### Output: 

| Name         | Type   | Description
| :---         | :---   | :---
| deletedCount | int64  | Returns the total deleted documents 

## Example
The below example shows the values for MongoDB Remove Document activity in the Flogo app JSON as well as the connection object referenced by the activity. The MongoDB Remove Document activity in this example deletes all the documents in a collection that has a key called 'location' with value' Palo Alto'

```json
{
  "id": "MongoDBRemoveDocument1",
  "name": "MongoDBRemoveDocument1",
  "description": "Mongodb Remove Document activity",
  "activity": {
    "ref": "#removeDocument",
    "settings": {
      "connection": "conn://a7730ae0-0199-11ea-9e1b-1b6d6afda988",
      "operation": "Remove Many Documents",
      "databaseName": "sample",
      "collectionName": "deletetesting",
      "timeout": "20"
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