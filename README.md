# BigQuery importer for Google Workspace

Go service forq importing data from Google Workspace into Bigquery.

##Usage

To use the service the following environment variable have to be set:
- **LOGGER_SERVICENAME**: Will add the ServiceContext to the log with the specified service name.
- **LOGGER_LEVEL**: The minimum enabled logging level. Recommended: Debug
- **LOGGER_DEVELOPMENT**: If the logger is set to development mode or not. Recommended: false
- **WORKSPACECLIENT_APIKEYSECRET**: The service assumes that the API key will be stored in GCP Secret Manager. This 
  variable must contain the full resource name of the resource that contains the API-KEY. Instructions on how to create 
  an API Key can be found [here][api-key]. The api key must contain the admin.directory.user and admin.directory.group
  scopes.
- **WORKSPACECLIENT_JWTCONFIG_SUBJECT**: The subject should be the email of an account that has privileges to make calls
  to the Admin API. 
- **BIGQUERYCLIENT_PROJECTID**: The project that contains the BigQuery Dataset
- **JOB_DATASET**: The name of the dataset that the data will be exported to.
- **JOB_ORG**: The organisation the data belongs to. 
- **WORK_DOMAIN**: The comain that the data should be fetched from. 
  

[api-key]:https://developers.google.com/admin-sdk/directory/v1/guides/delegation
