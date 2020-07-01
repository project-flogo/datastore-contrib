<!--
title: Neo4j Execute Cypher Query
weight: 4622
-->
# Neo4j Execute Cypher Query
This activity allows you to query Neo4j Graph DB using Cypher Query Language

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/neo4j/activity/executeCypherQuery
```

## Configuration

### Settings:
| Name                   | Type       | Description
| :---                   | :---       | :---    
| connection             | connection | Choose a Neo4j connection from the drop down  - ***REQUIRED***

### Input: 

| Name               | Type   | Description
| :---               | :---   | :---  
| cypherQuery        | string | The Cypher Query to execute


### Output: 

| Name   | Type | Description
| :---   | :--- | :---
| output | any  | Returns cypher query execution response

## Example


```json
{
            "id": "executeCypherQuery_2",
            "name": "Neo4j Execute Cypher Query",
            "description": "Neo4j Execute Cypher Query activity",
            "activity": {
              "ref": "#executeCypherQuery",
              "input": {
                "cypherQuery": "MATCH (n) RETURN n LIMIT 25"
              },
              "settings": {
                "accessMode": "Read",
                "databaseName": "neo4j",
                "connection": "conn://neo4jcon"
              }
            }
}

"connections": {
  "neo4jcon": {
    "ref": "github.com/project-flogo/datastore-contrib/neo4j/connection",
    "settings": {
				"name": "neo4jcon",
				"description": "",
				"connectionURI": "bolt://localhost:7687",
				"credType": "None",
				"username": "",
				"password": ""
	}
  }
}

```