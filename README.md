# BigQuery importer for Google Workspace

Go service for importing data from Google Workspace into Bigquery.

## Usage

To use the service the following environment variable have to be set:

- LOGGER_SERVICENAME: Will add the ServiceContext to the log with the specified service name.
- LOGGER_LEVEL: The minimum enabled logging level. Recommended: **debug**.
- LOGGER_DEVELOPMENT: If the logger is set to development mode or not. Recommended: **false**.
- WORKSPACECLIENT_APIKEYSECRET: The service assumes that the API key will be stored in GCP Secret Manager. This variable
  must contain the full resource name of the resource that contains the API-KEY. Instructions on how to create an API
  Key can be found [here][api-key]. The api key must contain the **admin.directory.user** and admin.directory.group**
  scopes.
- WORKSPACECLIENT_JWTCONFIG_SUBJECT: The subject should be the email of an account that has privileges to make calls to
  the Admin API.
- BIGQUERYCLIENT_PROJECTID: The id of the project where the tables will be created.
- JOB_DATASET: The name of the dataset where tables should be created.
- JOB_ORG: The organization the data belongs to.
- WORK_DOMAIN: The domain (including the prefix) that the data should be fetched from. Ex: example@**
  domain.domainprefix**

[api-key]:https://developers.google.com/admin-sdk/directory/v1/guides/delegation

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to add and update tests as appropriate.

Contributions should adhere to the [Conventional Commits][commits] specification.

## License

[MIT](https://choosealicense.com/licenses/mit/)


[commits]:https://www.conventionalcommits.org/en/v1.0.0/
