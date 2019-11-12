<!--
title: MongoDB Update Document
weight: 4622
-->
# MongoDB Update Document
This activity allows you to update one or more documents or replace a document from a collection in MongoDb database based on a matching criteria entered as a JSON

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/mongodb/activity/updateDocument
```

## Configuration

### Settings:
| Name                   | Type       | Description
| :---                   | :---       | :---    
| connection             | connection | Choose a MongoDB connection from the drop down  - ***REQUIRED***
| operation              | string     | The Query operation type (Update One Document or Update Many Documents or Replace One Document) - ***REQUIRED***
| databaseName           | string     | MongoDB databse to update documents - ***REQUIRED***
| collectionName         | string     | The collection within the MongoDB database to update documents - ***REQUIRED***  
| timeout                | int32      | Timeout in seconds for the activity's operations

### Input: 

| Name               | Type   | Description
| :---               | :---   | :---  
| criteria           | object | The matching criteria entered as a JSON - ***In case an empty JSON is provided and if operation is Update Many Documents, 
all the documents in the collection will be updated***
| updateData         | object | The data to be updated  with


### Output: 

| Name         | Type   | Description
| :---         | :---   | :---
| matchedCount | int64  | Returns the total documents that matched the criteria for update
| updatedCount | int64  | Returns the total documents that were updated

## Example
The below example shows the values for MongoDB Update Document activity in the Flogo app JSON as well as the connection object referenced by the activity. The MongoDB Update Document activity in this example updates all the documents in a collection that has a key called 'location' with value' Palo Alto' with the provided udpateData

```json
{
  "id": "MongoDBUpdateDocument",
  "name": "MongoDBUpdateDocument",
  "description": "Mongodb Update Document activity",
  "activity": {
    "ref": "#updateDocument",
    "settings": {
      "connection": "conn://a7730ae0-0199-11ea-9e1b-1b6d6afda988",
      "operation": "Update Many Documents",
      "databaseName": "sample",
      "collectionName": "test",
      "timeout": 0
    },
    "input": {
      "criteria": {
        "location": "Palo Alto"
      },
      "updateData": {
        "location": "San Francisco"
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