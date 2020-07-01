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
  
}
```