[![Functional Tests](https://github.com/elastic/apm-server/actions/workflows/functional-tests.yml/badge.svg)](https://github.com/elastic/apm-server/actions/workflows/functional-tests.yml)

# Functional Tests

The functional tests test that APM Server works as expected after version upgrades.

## Running the Tests

To run the tests, you will first need to set the `EC_API_KEY` environment variable, which can be obtained by following
[this guide](https://www.elastic.co/guide/en/cloud/current/ec-api-authentication.html).

Then, from the current directory, simply run:
```sh
go test -v -timeout=30m -cleanup-on-failure=false -target="pro" ./
```

You can also specify a specific test you want to run, for example:
```sh
go test -run=TestUpgrade_8_18_to_9_0 -v -timeout=30m -cleanup-on-failure=false -target="pro" ./
```

Note: Before running tests, make sure to delete the Terraforms by running `rm -r tf-*`.

### Debugging the Tests

If you get some errors after running the test, you can try heading to the [Elastic Cloud console](https://cloud.elastic.co/home)
in order to access the Kibana instance. From there, you can use Dev Tools to check the data streams etc.

Note: If the tests failed due to deployment, you may need to access the Elastic Cloud admin console instead to check the
deployment errors.

## Code Structure

The following is the simplified directory structure of functional tests.
```
- functionaltests/
   |- infra/
   |- internal/
   |- main_test.go
   |- x_y_test.go
```

All the functional tests are written in the current directory.

The `internal/` directory contains helper packages used in the tests, e.g. Elasticsearch, Kibana client wrapper etc.

The `infra/` directory contains infrastructure related code. In our case, we use Terraform for deploying the stack in
Elastic Cloud. The Terraform files are located in `infra/terraform`, and are copied into `tf-<test_name>/` e.g.
`tf-TestUpgrade_8_19_to_9_0/`, at the start of each test (since Terraform saves state in the directory it is initialized
in).

### Upgrade Tests

We suggest each upgrade test to be named in the format of `TestUpgrade_<from_version>_to_<to_version_1>[_to_<to_version_N>]*[_<suffix>]?`.
This means that the test will start from `from_version`, and be upgraded to `to_version_1`, then subsequently to
`to_version_2` etc. all the way to `to_version_N`.

The upgrade tests are implemented in each version test file. The test file is named after the last version of the upgrade
chain. For example, `TestUpgrade_8_15_to_8_16` will be in `8_16_test.go`.

### Standalone-to-Managed Tests

If the standalone-to-managed test includes an upgrade, simply add `Standalone_to_Managed` at the end of the test name,
e.g. `TestUpgrade_7_17_to_8_x_Standalone_to_Managed`. Otherwise, if there is no upgrade simply omit the `Upgrade`
prefix, e.g. `Test_7_17_Standalone_to_Managed`.

The standalone-to-managed tests are implemented in `standalone_test.go`.