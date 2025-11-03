# BSS
This is a test submission code, and includes a deployable golang based plan subscription and management code.

# How to run
To run this code, you need to have docker install. Then just run the following command to starup the server and all it's dependencies.

```
./dev up -d
```
The above command will spin up all the services mentioned in docker/docker-dev.yml, and allow you to use the postman collection included in this codebase.

# How to test
To run the test cases simply run

```
./test.sh
```

This will spin up a postgres database via a docker container, with definitions described in docker/docker-dev.yml, and then run the test cases against it. Be forewarned, the code will stop the postgres server once the test are completed.

# Exposed API.
The codebase exposes 7 API.
1. Get Plans
2. Get plan
3. Update plan
4. Create Plan.
5. Get user subscriptions
6. Subscribe
7. Unsubscribe.

All of thes API are defined in the BSS.postman_collection.json file. You can inport this file into postman, and run the API calls against the server.

# Logs.
Once you start up the system using '''./dev up -d''' you can view the logs of the server using the following command
```
./dev logs -f bss-server
```