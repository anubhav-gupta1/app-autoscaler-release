# App-AutoScaler [![Build Status](https://github.com/cloudfoundry/app-autoscaler/actions/workflows/postgres.yaml/badge.svg)](https://github.com/cloudfoundry/app-autoscaler/actions/workflows/postgres.yaml) [![Build Status](https://github.com/cloudfoundry/app-autoscaler/actions/workflows/mysql.yaml/badge.svg)](https://github.com/cloudfoundry/app-autoscaler/actions/workflows/mysql.yaml) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=cloudfoundry_app-autoscaler-release&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=cloudfoundry_app-autoscaler-release)

The `App-AutoScaler` provides the capability to adjust the computation resources for Cloud Foundry applications through

  * Dynamic scaling based on application performance metrics
  * Scheduled scaling based on time

The `App-AutoScaler` has the following components:

* `api` : provides public APIs to manage scaling policy
* `servicebroker`: implements the [Cloud Foundry service broker API][k]
* `metricsgateway` : collects and filters loggregator events via loggregator v2  API
* `metricsserver`: transforms loggregator events to app-autoscaler performance metrics ( metricsgateway + metricsserver is a replacement of metricscollector)
* `metricsforwarder`: receives and forwards custom metrics to loggregator via v2 ingress API
* `eventgenerator`: aggregates memory metrics, evaluates scaling rules and triggers events for dynamic scaling
* `scheduler`: manages the schedules in scaling policy and trigger events for scheduled scaling
* `scalingengine`: takes the scaling actions based on dynamic scaling rules or schedules

## Development

### System requirements

* Java 11 or above
* Docker
* [Apache Maven][b] 3
* [Cloud Foundry cf command line][f] 7 or 8
* Go 1.17 or above

### Database requirement

The `App-AutoScaler` supports Postgres and MySQL. It uses Postgres as the default backend
data store. These are run up locally with docker images so ensure that docker is working on
your system before running up the tests.

### Setup
**Note:** all of the setup is encapsulated in the makefile targets. So you can run the test targets (test|integration) directly
and it will setup and start the tests.

To set up the development, firstly clone this project

```shell
$ git clone https://github.com/cloudfoundry/app-autoscaler.git
```

Generate [scheduler test certs](https://github.com/cloudfoundry/app-autoscaler/blob/main/scheduler/README.md#generate-certificates)


#### Initialize the Database

**Note:** The makefile will init the database if it has not already been run before running the tests.

* **Postgres**

  ```shell
  make init-db
  ```

* **MySQL**

  ```shell
  make init-db db_type=mysql
  ```

#### Generate TLS Certificates

Create the certificates.

**Note**:

  * on macos it will install `certstrap` automatically but on other OS's it needs to be pre-installed
  * The makefile will create the certificates if it has not already been run before running the tests.

```shell
make test-certs
```

### Unit tests
The default database is postgres

  * **Postgres**:

  ```shell
  make test
  ```

To use a specific postgres version:

```shell
make clean #Only if you're changing versions to refresh the running docker image.
make test POSTGRES_TAG=x.y
```

where:

  * x is the major version
  * y is the minor version ( this can be left out to get the most recent patch)
  * **MySQL**:

    ```shell
    make test db_type=mysql
    ```

To use a specific MySQL version:

```shell
make clean #Only if you're changing versions to refresh the running docker image.
make test db_type=mysql MYSQL_TAG=x.y
```

where:

  * x is the major version
  * y is the minor version ( this can be left out to get the most recent patch)


### Integration tests
The default database is postgres

  * **Postgres**:

  ```shell
  make integration
  ```

To use a specific postgres version:

```shell
make clean #Only if you're changing versions to refresh the running docker image.
make integration POSTGRES_TAG=x.y
```

where:

  * x is the major version
  * y is the minor version ( this can be left out to get the most recent patch)
  * **MySQL**:

    ```shell
    make integration db_type=mysql
    ```

To use a specific MySQL version:

```shell
make clean #Only if you're changing versions to refresh the running docker image.
make integration db_type=mysql MYSQL_TAG=x.y
```

where:

  * x is the major version
  * y is the minor version ( this can be left out to get the most recent patch)

### Build App-AutoScaler

```shell
make build
```

### Clean up

You can use the  `make clean` to remove:

  * database ( postgres or mysql)
  * autoscaler build artifacts

### Coding Standards

Autoscaler uses Golangci and Checkstyle for its code base. Refer to [style-guide](style-guide/README.md)

## Bosh Release for app-autoscaler service

## Purpose

The purpose of this bosh release is to deploy and setup the [app-autoscaler](https://github.com/cloudfoundry-incubator/app-autoscaler) service.

## Usage

### Bosh Lite Deployment

* Install [Bosh-cli-v2](https://bosh.io/docs/cli-v2.html#install)
* Install and start [BOSH-Deployment](https://github.com/cloudfoundry/bosh-deployment), following its [README](https://github.com/cloudfoundry/bosh-deployment/blob/master/README.md).
* Install [CF-deployment](https://github.com/cloudfoundry/cf-deployment#deploying-cf)
* Create a new autoscaler client
  UAA CLI is required to here to create a new UAA client id.
  * Install the UAA CLI, `uaac`.

    ```sh
    gem install cf-uaac
    ```

  * Obtain `uaa_admin_client_secret`

    ```sh
    bosh interpolate --path /uaa_admin_client_secret /path/to/cf-deployment/deployment-vars.yml
    ```

  * Use the `uaac target uaa.YOUR-DOMAIN` command to target your UAA server and obtain an access token for the admin client.

    ```sh
    uaac target uaa.bosh-lite.com --skip-ssl-validation
    uaac token client get admin -s <uaa_admin_client_secret>
    ```

  * Create a new autoscaler client

    ```sh
    uaac client add "autoscaler_client_id" \
        --authorized_grant_types "client_credentials" \
        --authorities "cloud_controller.read,cloud_controller.admin,uaa.resource" \
        --secret <AUTOSCALE_CLIENT_SECRET>
    ```

* Create and upload App-Autoscaler release

  ```sh
  git clone https://github.com/cloudfoundry/app-autoscaler-release
  cd app-autoscaler-release
  make mod-tidy vendor db scheduler
  bosh create-release
  bosh -e YOUR_ENV upload-release
  ```

* Deploy app-autoscaler with the newly created autoscaler client

  In the latest App-Autoscaler v2.0 release , App-Autoscaler retrieves application's metrics with [loggregator V2 API](https://github.com/cloudfoundry/loggregator-api/blob/master/README.md) via gRPC over mutual TLS connection.

  So the valid TLS certification to access Loggregator Reverse Log Proxy is required here.   When deploying in bosh-lite, the most easy way is to provide loggregator certificates generated by `cf-deployments`.

  ```sh
  bosh -e YOUR_ENV -d app-autoscaler \
      deploy templates/app-autoscaler-deployment.yml \
      --vars-store=bosh-lite/deployments/vars/autoscaler-deployment-vars.yml \
      -l <PATH_TO_CF_DEPLOYMENT_VAR_FILES> \
      -v system_domain=bosh-lite.com \
      -v cf_client_id=autoscaler_client_id \
      -v cf_client_secret=<AUTOSCALE_CLIENT_SECRET> \
      -v skip_ssl_validation=true
  ```

* Deploy autoscaler with cf deployment mysql database

  **Notes**: It is blocked by the pull request [cf-deployment #881](https://github.com/cloudfoundry/cf-deployment/pull/881) temporarily. If you would like to use the cf mysql, please apply the `set-autoscaler-db.yml` in the pull request when deploy cf deployment.

  The lastest Autoscaler release add the support for mysql database, Autoscaler can connect the same mysql database with cf deployment. Use the operation file `example/operation/cf-mysql-db.yml` which including the cf database host , password and tls.ca cert.

  ```sh
  bosh -e YOUR_ENV -d app-autoscaler \
      deploy templates/app-autoscaler-deployment.yml \
      --vars-store=bosh-lite/deployments/vars/autoscaler-deployment-vars.yml \
      -l <PATH_TO_CF_DEPLOYMENT_VAR_FILES> \
      -v system_domain=bosh-lite.com \
      -v cf_client_id=autoscaler_client_id \
      -v cf_client_secret=<AUTOSCALE_CLIENT_SECRET> \
      -v skip_ssl_validation=true \
      -o example/operation/cf-mysql-db.yml
  ```

* Deploy autoscaler with external postgres database and mysql database

  ```sh
  bosh -e YOUR_ENV -d app-autoscaler \
      deploy templates/app-autoscaler-deployment.yml \
      --vars-store=bosh-lite/deployments/vars/autoscaler-deployment-vars.yml \
      -l <PATH_TO_CF_DEPLOYMENT_VAR_FILE> \
      -l <PATH_TO_DATABASE_VAR_FILE> \
      -v system_domain=bosh-lite.com \
      -v cf_client_id=autoscaler_client_id \
      -v cf_client_secret=<AUTOSCALE_CLIENT_SECRET> \
      -v skip_ssl_validation=true \
      -o example/operation/external-db.yml
  ```

>** The DATABASE_VAR_FILE should look like as below

  ```sh
  database:
    name: <database_name>
    host: <database_host>
    port: <database_port>
    scheme: <database_scheme>
    username: <database_username>
    password: <database_password>
    sslmode: <database_sslmode>
    tls:
      ca: |
        -----BEGIN CERTIFICATE-----

        -----END CERTIFICATE-----
  ```

The table below shows the description of all the variables:

Property | Description
-------- | -------------
database.name | The database name.
database.host | The database server ip address or hostname.
database.port | The database server port.
database.scheme | The database scheme. Currently Autoscaler supports "postgres" and "mysql".
database.username | The username of the database specified above in "database.name".
database.password | The password of the user specified above in "database.username".
database.sslmode | There are 6 values allowed for "postgres": disable, allow, prefer, require, verify-ca and verify-full. Please refer to [Postgres SSL definition](https://www.postgresql.org/docs/current/libpq-ssl.html) when define `database_sslmode`.  For "mysql", there are 7 values allowed: false, true, skip-verify, preferred, verify-ca, verify_identity.Please refer to [Mysql SSL definition(Golang)](https://github.com/go-sql-driver/mysql#tls) and [Mysql Connector SSL](https://dev.mysql.com/doc/connector-j/8.0/en/connector-j-reference-using-ssl.html)
database.tls.ca | PEM-encoded certification authority for secure TLS communication. Only required when sslmode is verify-ca or verify-full(postgres) or verify_identity(mysql) and can be omitted for other sslmode.

## Register service

Log in to Cloud Foundry with admin user, and use the following commands to register `app-autoscaler` service

```sh
cf create-service-broker autoscaler <brokerUserName> <brokerPassword> <brokerURL>
```

* `brokerUserName`: the user name to authenticate with service broker. It's default value is `autoscaler_service_broker_user`.
* `brokerPassword`: the password to authenticate with service broker. It will be stored in the file passed to the --vars-store flag (bosh-lite/deployments/vars/autoscaler-deployment-vars.yml in the example). You can find them by searching for `autoscaler_service_broker_password`.
* `brokerURL`: the URL of the service broker

All these parameters are configured in the bosh deployment. If you are using default values of deployment manifest, register the service with the commands below.

```sh
cf create-service-broker autoscaler autoscaler_service_broker_user `bosh int ./bosh-lite/deployments/vars/autoscaler-deployment-vars.yml --path /autoscaler_service_broker_password` https://autoscalerservicebroker.bosh-lite.com
```

## Acceptance test

Refer to [AutoScaler UAT guide](src/acceptance/README.md) to run acceptance test.

## Use service

To use the service to auto-scale your applications, log in to Cloud Foundry with admin user, and use the following command to enable service access to all or specific orgs.

```sh
cf enable-service-access autoscaler [-o ORG]
```

The following commands don't require admin rights, but user needs to be Space Developer. Create the service instance, and then bind your application to the service instance with the policy as parameter.

```sh
cf create-service autoscaler  autoscaler-free-plan  <service_instance_name>
cf bind-service <app_name> <service_instance_name> -c <policy>
```

## Remove the service

Log in to Cloud Foundry with admin user, and use the following commands to remove all the service instances and the service broker of `app-autoscaler` from Cloud Foundry.

```sh
cf purge-service-offering autoscaler
cf delete-service-broker autoscaler
```

## Monitoring the service

The app-autoscaler provides a number of health endpoints that are available externally that can be used to check the state of each component. Each health endpoint is protected with basic auth (apart from the api server), the usernames are listed in the table below, but the passwords are available in credhub.

Component | Health URL | Username | Password Key |
--------- | -----------| ---------| -------------|
eventgenerator|https://autoscaler-eventgenerator.((system_domain))/health|eventgenerator|/autoscaler_eventgenerator_health_password|
metricsforwarder|https://autoscaler-metricsforwarder.((system_domain))/health|metricsforwarder|/autoscaler_metricsforwarder_health_password|
metricsgateway|https://autoscaler-metricsgateway.((system_domain))/health|metricsgateway|/autoscaler_metricsgateway_health_password|
metricsserver|https://autoscaler-metricsserver.((system_domain))/health|metricsserver|/autoscaler_metricsserver_health_password|
scalingengine|https://autoscaler-scalingengine.((system_domain))/health|scalingengine|/autoscaler_scalingengine_health_password|
operator|https://autoscaler-operator.((system_domain))/health|operator|/autoscaler_operator_health_password|
scheduler|https://autoscaler-scheduler.((system_domain))/health|scheduler|/autoscaler_scheduler_health_password|

These endpoints can be disabled by using the ops file `example/operations/disable-basicauth-on-health-endpoints.yml`

You can follow the development progress on [Pivotal Tracker][t].

## Deploy and offer Autoscaler as a service

Go to [app-autoscaler-release][r] project for how to BOSH deploy `App-AutoScaler`

## Use Autoscaler service

Refer to [user guide][u] for the details of how to use the Auto-Scaler service, including policy definition, supported metrics, public API specification and command line tool.

## License

This project is released under version 2.0 of the [Apache License][l].


[b]: https://maven.apache.org/
[c]: http://couchdb.apache.org/
[d]: http://www.eclipse.org/m2e/
[e]: http://www.cloudant.com
[f]: https://github.com/cloudfoundry/cli/releases
[k]: http://docs.cloudfoundry.org/services/api.html
[l]: LICENSE
[t]: https://www.pivotaltracker.com/projects/1566795
[p]: https://www.postgresql.org/
[r]: https://github.com/cloudfoundry/app-autoscaler-release
[u]: docs/Readme.md
[m]: https://www.mysql.com/
