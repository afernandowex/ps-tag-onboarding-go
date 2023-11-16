# ps-tag-onboarding-go
This is an implementation of the exercise described here - https://wexinc.atlassian.net/wiki/spaces/TGT/pages/153576505378/Developer+Onboarding+Exercise+-+Advanced

# How to run

`docker compose up` in the root folder where the docker-compose.yml exists. That will bring up MySQL and the go application.

Wait until this is shown in the console:

```
ps-tag-onboarding-go-app-1  |    ____    __
ps-tag-onboarding-go-app-1  |   / __/___/ /  ___
ps-tag-onboarding-go-app-1  |  / _// __/ _ \/ _ \
ps-tag-onboarding-go-app-1  | /___/\__/_//_/\___/ v4.11.2
ps-tag-onboarding-go-app-1  | High performance, minimalist Go web framework
ps-tag-onboarding-go-app-1  | https://echo.labstack.com
ps-tag-onboarding-go-app-1  | ____________________________________O/_______
ps-tag-onboarding-go-app-1  |                                     O\
ps-tag-onboarding-go-app-1  | 2023/11/01 03:12:56 Starting server on port 8080
ps-tag-onboarding-go-app-1  | ⇨ http server started on [::]:8080
```

# Post Requests

To return some test users run:

`curl http://localhost:8080/find/[1-3]` 

To add a user try running this:

`curl -X POST -H "Content-Type: application/json" -d '{"firstName":"WexFirst2","lastName":"WexLast2","email":"wexfirst.wexlast2@wexinc.com","age":20}' http://localhost:8080/save` 

Since the name already exists in the database an error message is received. The test users inserted when the app starts up are [here](https://github.com/afernandowex/ps-tag-onboarding-go/blob/main/internal/app/user-api/mysql/mysql.go#L56)

`{"error":"Name already exists"}`

Try changing the first and last name and retrieving the record using the `find` API.

# Components chosen for this build

### MySQL
The backend persistence store is MySQL. It's a relatively easy database that has broad mindshare in the community. Helps with developers onboarding quickly without too much of a learning curve. It also has a broad community of devs on SO to help with issues if any. There is also a fair amount of ubiquity in the cloud as first class managed MySQLs - i.e. AWS RDS for example.

### Echo minimalist web framework
This is a lite high performance minimalist Go web framework. Very little boilerplate to get something up and running. Since we are looking at 2 endpoints (find and save), this should suffice at the moment.

