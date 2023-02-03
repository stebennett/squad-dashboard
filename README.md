# squad-dashboard

This repository contains a number of commands which can be run to collect data from services and store in a local database for processing.

- Jira Issue Collector - Finds issues based on a Jira Query and adds to database, including information about all transitions
- Github PR Collector - Under development 🏗 - Collects Pull-Requests in an organisation and stores key information


## Getting started
1. Install and open Docker
2. `homebrew install bazel`
3. `make deps`
4. Create your environment file
   - Duplicate the `.env_template` file as `.env`
   - Update values in `.env` to include your GH token, etc
5. `ENV=.env make init db migrate metabase`


## Building and Running

Build all services and init the databases

````
 # ENV=.env make init
````

Running all tests

````
 # ENV=.env make test
````

Building and run a command in docker

````
 # ENV=.env make {command-name}
````

## Populating the database

1. run up the database
2. run the migrations
3. run the config loaders
4. run your collectors
5. run any calculator (if necessary)
6. run metabase to inspect and build charts

````
 # ENV=.env make db
 # ENV=.env make migrate
 # ENV=.env make {command-name}
 # ENV=.env make metabase
````

## Running JIRA collector and calculator 

1. ENV=.env make db
2. ENV=.env make migrate
3. ENV=.env make init
4. ENV=.env make jiracollector
5. ENV=.env make jiraissuecalculator
6. ENV=.env make metabase

## Running a report

````
  # ENV=.env make report
````

## Troubleshooting
- Error: `xcrun: error: invalid active developer path (/Library/Developer/CommandLineTools), missing xcrun at: /Library/Developer/CommandLineTools/usr/bin/xcrun`
  - Fix it with: `xcode-select --install`
- Error when running `make deps` : ` unable to find valid certification path to requested target`
  - Fix it by exiting zScaler