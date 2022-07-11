# Token Generator API

This microservice is exposing RESTful APIs to generate,validate,visualize invite tokens

## API Documentation

API documentation can be found in `./doc/api/openapi.yaml`.

Use [Swagger Editor](https://editor.swagger.io/) to view the file.

## Code Quality Check

`SonarQube` is used for checking code quality. Much of the needed sonarqube configurations
are in the `sonar-project.properties` file.

To run the code quality check `sonar-scanner` is needed.
Download `sonar-scanner` from [here](https://docs.sonarqube.org/display/SCAN/Analyzing+with+SonarQube+Scanner).

Also an authentication key is needed to access the `SonarQube` web service.
Follow [these](https://docs.sonarqube.org/7.4/user-guide/user-token/) instructions to setup an authentication key.

Use following command to run code quality check.
```bash
/<path-to-sonar-scanner-location-in-filesystem>/bin/sonar-scanner \
  -Dsonar.login=<authentication-key>
```
