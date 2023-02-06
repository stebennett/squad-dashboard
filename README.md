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
1. Error: `xcrun: error: invalid active developer path (/Library/Developer/CommandLineTools), missing xcrun at: /Library/Developer/CommandLineTools/usr/bin/xcrun`
   - Fix it with: `xcode-select --install`
2. Error when running `make deps` : ` unable to find valid certification path to requested target`
   - This is a VPN issue, fix it by exiting zScaler
3. Error about missing python: `line 119: python: command not found`
   - Most likely you have python3.9 installed via homebrew, you can fix it by running `echo 'export PATH=/opt/homebrew/opt/python@3.9/libexec/bin:$PATH' >> ~/.zprofile` and then `. ~/.zprofile`
4. `jiracollector_1             | 2023/02/06 13:48:10 Post "https://stashinvest.atlassian.net/rest/api/2/search": x509: certificate signed by unknown authority`
   - This is a VPN issue, fix it by exiting zScaler
5. `squad-dashboard-configloader-1  | 2023/02/06 13:54:14 failed to save Jira Work ToDo States. pq: relation "jira_work_states" does not exist`
   - Not sure about this one