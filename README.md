# theScore "the Rush" Interview Challenge
At theScore, we are always looking for intelligent, resourceful, full-stack developers to join our growing team. To help us evaluate new talent, we have created this take-home interview question. This question should take you no more than a few hours.

**All candidates must complete this before the possibility of an in-person interview. During the in-person interview, your submitted project will be used as the base for further extensions.**

### Why a take-home challenge?
In-person coding interviews can be stressful and can hide some people's full potential. A take-home gives you a chance work in a less stressful environment and showcase your talent.

We want you to be at your best and most comfortable.

### A bit about our tech stack
As outlined in our job description, you will come across technologies which include a server-side web framework (like Elixir/Phoenix, Ruby on Rails or a modern Javascript framework) and a front-end Javascript framework (like ReactJS)

### Challenge Background
We have sets of records representing football players' rushing statistics. All records have the following attributes:
* `Player` (Player's name)
* `Team` (Player's team abbreviation)
* `Pos` (Player's postion)
* `Att/G` (Rushing Attempts Per Game Average)
* `Att` (Rushing Attempts)
* `Yds` (Total Rushing Yards)
* `Avg` (Rushing Average Yards Per Attempt)
* `Yds/G` (Rushing Yards Per Game)
* `TD` (Total Rushing Touchdowns)
* `Lng` (Longest Rush -- a `T` represents a touchdown occurred)
* `1st` (Rushing First Downs)
* `1st%` (Rushing First Down Percentage)
* `20+` (Rushing 20+ Yards Each)
* `40+` (Rushing 40+ Yards Each)
* `FUM` (Rushing Fumbles)

In this repo is a sample data file [`rushing.json`](/rushing.json).

##### Challenge Requirements
1. Create a web app. This must be able to do the following steps
    1. Create a webpage which displays a table with the contents of [`rushing.json`](/rushing.json)
    2. The user should be able to sort the players by _Total Rushing Yards_, _Longest Rush_ and _Total Rushing Touchdowns_
    3. The user should be able to filter by the player's name
    4. The user should be able to download the sorted data as a CSV, as well as a filtered subset
    
2. The system should be able to potentially support larger sets of data on the order of 10k records.

3. Update the section `Installation and running this solution` in the README file explaining how to run your code

### Submitting a solution
1. Download this repo
2. Complete the problem outlined in the `Requirements` section
3. In your personal public GitHub repo, create a new public repo with this implementation
4. Provide this link to your contact at theScore

We will evaluate you on your ability to solve the problem defined in the requirements section as well as your choice of frameworks, and general coding style.

### Help
If you have any questions regarding requirements, do not hesitate to email your contact at theScore for clarification.

### Installation and running this solution

Dependencies: `docker`, `mysql-client-core-8.0`, `go`. On an Ubuntu 20.04 machine, you can install these via:
- `sudo snap install docker`
- `sudo apt install mysql-client-core-8.0`
- Instructions for downloading and installing Go: https://golang.org/doc/install
- setup `$GOPATH` env variable and add the Go binary to `$PATH` env variable

Preparation
- Create the necessary directory structure: `mkdir -p $GOPATH/src/github.com/geekgunda/`
- Clone the git repo in the directory above: `git clone git@github.com:geekgunda/nfl-rushing.git`

Installation (using Makefile):
- `make setup`: download and setup mysql-server 8.0 docker image (using docker-compose)
- `make build`: compile go executable binary for this app (using go binary)
- `make run`  : start the app 
- `make clean`: stop docker containers and remove them
- `make test` : resets DB and runs automated test cases

In a browser, head to http://127.0.0.1:8081/rushingstats to access the web app.

##### Configuration and options:
- By default `main.go` will always import data.  
  To avoid duplication of data in DB, make sure you comment `shouldImport = true` line in `main.go`, before restarting the app.
- To import stats from a different file, just over-write `rushing.json` with the new file. Ensure filename is exactly same.  
  Alternatively, change the file name in `importer.go` under `statsFile` constant. Then reset DB, and start from scratch.

##### Points to note:
- Ensure nothing is running on ports 3306 and 8081 on local machine
- Ensure the app base directory is: `$GOPATH/src/github.com/geekgunda/nfl-rushing`
- `make setup` might need to run with `sudo`, if docker is running as `root` user
- `make setup` might fail as docker container might not have started by then. Retry the command in that case

##### Design
The problem is solved entirely using Go, instead of a separate front-end framework.  
So the UI is pretty bare-bones and nothing flashy. It just gets the job done.  
For rendering UI, Go's awesome `template` library has been used.

APIs:
- A single endpoint handles everything
- GET request will render webpage with all available stats
- POST request with filter params render the same webpage with matching stats
- POST request with "Download" button triggers a csv export and download with the requested filters

Database:
- a MySQL server is used as the datastore
- JSON data type is used to dump the records as is
- Only the fields, that require filtering are parsed and populated as separate columns with indexes on them

Scale:
- Since data is dumped into a DB, the solution is scalable at a high level
- Pagination support is missing, but with existing APIs and small changes in the template, it's achievable

##### Pending updates (compromises due to time constraints)
- Only a few invalid input types are handled during import. (Refer to `cleanStats` fn inside `importer.go` for reference)
- Pagination is not added within the UI template. Currently standard DB query limits max records to 1000
- Automated tests for http handlers have not been added yet
- Currently the filters in UI can be independently updated, and a download requested.  
  This will lead to downloaded file, giving data with new filters, while the UI showing data from previous request.  
  This can be a bit confusing as a user.

##### Possible improvements (future updates for production rollout)
- The APIs can be split into separate endpoints for responding with data and file download.  
  This will simplify response format and make them easier to test too.
- All fields can be parsed and populated into separate columns in MySQL DB.  
  This way, if a filter or ordering on a new column is required, it'll be just a matter of adding another index in DB.
