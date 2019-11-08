<!--
title: MongoDB Trigger
weight: 4622
-->
# MongoDB Trigger
This activity allows you to listen for Create, Update and Delete Document events in a single collection or all collections in a MongoDb database.

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/mongodb/trigger/mongodb
```

## Configuration

### Settings:
| Name                     | Type       | Description
| :---                     | :---       | :---
| connection               | connection | Choose a MongoDB connection from the drop down  - ***REQUIRED***


### Handler Settings
| Name                | Type   | Description
| :---                | :---   | :---
| databaseName        | string | MongoDB Database name - ***REQUIRED***
| collectionName      | string | The collection to listen to for changes. If left blank, listens to all collections in a DB
| listenInsert        | bool   | Should the trigger listen to Insert events?
| listenUpdate        | bool   | Should the trigger listen to Update events?
| listenRemove        | bool   | Should the trigger listen to Remove events?


### Output:

| Name   | Type   | Description
| :---   | :---   | :---
| output | object | A JSON object of key value pairs containing the information about the MongoDB event that the trigger is configuerd to listen

## Example
The below example shows the values for MongoDB Trigger in the Flogo app JSON as well as the connection object referenced by the trigger. The MongoDB trigger listens to all events of a collection called test in a MongoDB database named sample

```json
"triggers": [
  {
    "ref": "github.com/project-flogo/datastore-contrib/mongodb/trigger/mongodb",
    "name": "mongodb-eventlistener",
    "description": "",
    "settings": {
      "connection": "conn://a7730ae0-0199-11ea-9e1b-1b6d6afda988"
    },
    "id": "MongoDBTrigger",
    "handlers": [
      {
        "description": "",
        "settings": {
          "databaseName": "sample",
          "collectionName": "test",
          "listenInsert": true,
          "listenUpdate": true,
          "listenRemove": true
        },
        "action": {
          "ref": "github.com/project-flogo/flow",
          "settings": {
            "flowURI": "res://flow:nov8tr"
          },
          "input": {
            "output": "=$.output"
          }
        },
        "schemas": {
          "output": {
            "output": {
              "type": "json",
              "value": "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"definitions\":{},\"properties\":{\"NameSpace\":{\"type\":\"object\"},\"OperationType\":{\"type\":\"string\"},\"ResultDocument\":{\"type\":\"object\"}}}",
              "fe_metadata": "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"definitions\":{},\"properties\":{\"NameSpace\":{\"type\":\"object\"},\"OperationType\":{\"type\":\"string\"},\"ResultDocument\":{\"type\":\"object\"}}}"
            }
          }
        }
      }
    ]
  }
]

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