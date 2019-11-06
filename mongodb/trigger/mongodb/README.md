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
| mongodbConnection        | connection | hoose a MongoDB connection from the drop down  - ***REQUIRED***


### Handler Settings
| Name                | Type   | Description
| :---                | :---   | :---
| collectionName      | string | The collection to listen to for changes. If left blank, listens to all collections in a DB
| listenInsert        | bool   | Should the trigger listen to Insert events?
| listenUpdate        | bool   | Should the trigger listen to Update events?
| listenRemove        | bool   | Should the trigger listen to Remove events?


### Output:

| Name   | Type   | Description
| :---   | :---   | :---
| output | object | A JSON object of key value pairs containing the information about the MongoDB event that the trigger is configuerd to listen

## Example

```json
"triggers": [
  {
    "ref": "#mongodb",
    "name": "mongodb-eventlistener",
    "description": "",
    "settings": {
      "mongodbConnection": "conn://b5002f80-ffe6-11e9-9e1b-1b6d6afda988"
    },
    "id": "MongoDBTrigger",
    "handlers": [
      {
        "description": "",
        "settings": {
          "collectionName": "test",
          "listenInsert": true,
          "listenUpdate": true,
          "listenRemove": true
        },
        "action": {
          "ref": "github.com/project-flogo/flow",
          "settings": {
            "flowURI": "res://flow:trigger1"
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