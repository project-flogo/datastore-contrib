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
| mongoConnection        | connection | Choose a MongoDB connection from the drop down  - ***REQUIRED***
| operation              | string     | The Query operation type (Find One Document or Find Many Documents) - ***REQUIRED***
| collectionName         | string     | The collection to work on - ***REQUIRED***  

### Input: 

| Name               | Type   | Description
| :---               | :---   | :---  
| jsonDocument       | object | The matching critiera entered as a JSON


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
    "ref": "#queryDocument",
    "settings": {
      "mongoConnection": "conn://b5002f80-ffe6-11e9-9e1b-1b6d6afda988",
      "operation": "Find Many Documents",
      "collectionName": "test"
    },
    "input": {
      "jsonDocument": {
        "location": "Palo Alto"
      }
    }
  }
}

"connections": {
  "b5002f80-ffe6-11e9-9e1b-1b6d6afda988": {
    "id": "b5002f80-ffe6-11e9-9e1b-1b6d6afda988",
    "name": "mc2",
    "ref": "#connection",
    "settings": {
      "Name": "mc2",
      "Description": "",
      "ConnectionURI": "<Connection URI Here>",
      "Database": "<DB Name here>",
      "CredType": "<One of None or SCRAM-SHA-1 or SCRAM-SHA-256>",
      "UserName": "<Enter Username here in case of SCRAM-SHA-1 or SCRAM-SHA-256>",
      "Password": "",
      "Ssl": false,
      "X509": false,
      "TrustCert": "",
      "ClientCert": "",
      "ClientKey": "",
      "KeyPass": "",
      }
  }
}
```