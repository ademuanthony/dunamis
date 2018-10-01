# API for dunamis church Seed of Destiny


The app provides the following 

* RESTful endpoints for getting seed of destiny for the day which scrap the content from https://www.dunamisgospel.org/index.php/component/k2/itemlist/category/3
* Standard CRUD for dunamis artists
 
The app uses the following Go packages 

* Routing framework: [ozzo-routing](https://github.com/go-ozzo/ozzo-routing)
* Database: [ozzo-dbx](https://github.com/go-ozzo/ozzo-dbx)
* Data validation: [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
* Logging: [logrus](https://github.com/Sirupsen/logrus)
* Configuration: [viper](https://github.com/spf13/viper)
* Dependency management: [dep](https://github.com/golang/dep)
* Testing: [testify](https://github.com/stretchr/testify)


## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://golang.org/doc/install) to
install Go on your computer. The kit requires Go 1.10 or above.

After installing Go, run the following commands to download and install this starter kit:

```shell
# install the starter kit
go get github.com/ademuanthony/dunamis

# install dep
$ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# fetch the dependent packages
cd $GOPATH/ademuanthony/dunamis
dep ensure
```

## Project Structure

This starter kit divides the whole project into four main packages:

* `models`: contains the data structures used for communication between different layers.
* `services`: contains the main business logic of the application.
* `daos`: contains the DAO (Data Access Object) layer that interacts with persistent storage.
* `apis`: contains the API layer that wires up the HTTP routes with the corresponding service APIs.

[Dependency inversion principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)
is followed to make these packages independent of each other and thus easier to test and maintain.

The rest of the packages in the kit are used globally:
 
* `app`: contains routing middlewares and application-level configurations
* `errors`: contains error representation and handling
* `util`: contains utility code

The main entry of the application is in the `server.go` file. It does the following work:

* load external configuration
* establish database connection
* instantiate components and inject dependencies
* start the HTTP server
