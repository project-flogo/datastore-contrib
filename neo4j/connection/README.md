<!--
title: Neo4j Connection
weight: 4622
-->
# Neo4j Connection
This connection allows you to configure properties necessary to establish a connection with a Neo4j Graph DB. A Neo4j Connection is necessary to work with the activities and trigger under neo4j contribution.

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/datastore-contrib/neo4j/connection
```

## Configuration

### Settings:
| Name             | Type       | Description
| :---             | :---       | :---    
| name             | string     | A name for the connection  - ***REQUIRED***
| description      | string     | A short description for the connection
| connectionURI    | string     | Neo4j instance connection URI - ***REQUIRED***
| credType         | string     | Credential Type e.g None, BasicAuth
| username         | string     | Username of the Neo4j instance
| password         | string     | Password of the Neo4j instance

## Example
A sample Neo4j connection

```json
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

## Testing

Launch Neo4j docker container using below command
```bash
docker run  --publish=7474:7474 --publish=7687:7687  --volume=$HOME/neo4j/data:/data  neo4j
```
Open "http://localhost:7474/browser" and verify connectivity. You can optionally load sample data by running ":play movie-graph" command and follow the steps.

