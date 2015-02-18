![docker](https://docker.com)
![dockerinstall](https://docs.docker.com/installation/#installation)


go-events-service
=================

A simple Go web service allowing storage of time-based events.


## Recording an event

```
POST /events
{
	"name": "test",
	"timestamp": "2015-02-11T15:01:00+00:00"
}
```


## Aggregating events

```
GET /events/count?from=2015-02-11T15:01:00+00:00&to=2015-02-11T15:01:59+00:00
{
	"test": 1
}
```


## Playing around

### Docker

The service has been containerised using [Docker][docker]. If you'd like to
play around, you need to [install docker][dockerinstall] first (it's worth it!). 

## Running

A makefile is provided to help with common tasks. To spin up a working instance
of the service, simply type `make run` at the repository root. Once the service 
has started, `docker ps` will show it's running containers:

```
$ docker ps
CONTAINER ID        IMAGE                             STATUS              PORTS                     NAMES
5ee27d857668        go-events-service:latest          Up 5 seconds        0.0.0.0:49161->5000/tcp   go-events-service-app
6ede52667e85        redis:latest                      Up 24 hours         6379/tcp                  go-events-service-redis
```

As you can see, the service comprises two Docker containers by default, one for
the Go web application and one for the backing redis datastore. The app container
exposes its functionality on port 5000 and this is bound to a high port on the
host machine (49161 in this example). In this example, to store an event via the 
service, we would make a POST request to `http://localhost:49161/events`. If you
are running on MacOSX using boot2docker, you will access the service via the 
boot2docker host IP rather than localhost. This can be obtained as follows:

```
$ boot2docker ip
192.168.59.103
``` 

## Testing

Similarly, a make target is provided to run the application's test suite. `make test`
spins up a Docker container with all the system-level and Go dependencies required to
run the tests, printing test results and coverage before exiting.
