# ${docGenStepName}

## ${docGenDescription}

## Prerequisites

This step will execute every unit test associated with a package belonging to the specified local repository on an ABAP system.
Learn more about gCTS [here](https://help.sap.com/viewer/4a368c163b08418890a406d413933ba7/201909.001/en-US/f319b168e87e42149e25e13c08d002b9.html).

## ${docGenParameters}

## ${docGenConfiguration}

## ${docJenkinsPluginDependencies}

## Example

Example configuration for the use in a Jenkinsfile.

```groovy
gctsCreateRepository(
  script: this,
  host: "https://abap.server.com:port",
  client: "000",
  credentialsId: 'ABAPUserPasswordCredentialsId',
  repository: "myrepo"
  )
```

Example configuration for the use in a yaml config file (such as `.pipeline/config.yaml`).

```yaml
steps:
  <...>
  gctsCreateRepository:
    host: "https://abap.server.com:port"
    client: "000"
    username: "ABAPUsername"
    password: "ABAPPassword"
    repository: "myrepo"
```
