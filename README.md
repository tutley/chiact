# chiact

The goal of this project is to create a boilerplate for developing full-stack web applications using [Golang](https://golang.org) (Go) on the server side and [React](https://facebook.github.io/react/) on the client side.

chiact includes:
* Server side routing with [Chi](https://github.com/pressly/chi)
* Database connection
* Static directory serving
* API structure and starter methods
* Authentication and JWT
* User Storage and editing
* A simple blog

### Go Structure
The server side, written with Go, uses glide for vendoring (dependency management) and can be run locally with "go run main.go". The glide.yaml file holds the dependencies, which can be installed with 'glide install'.

* The handlers directory contains the various local packages for handling routes.
* The helpers directory contains packages with functions to assist along the way.
* The models directory contains data models with methods included.
* The vendor directory contains the local copies of dependencies.

We have included mgo for storing data in MongoDB but if you prefer another database that could be changed by editing the models as well as the db init in main and helpers.

### React Structure
The client application is written in React, which requires a few files in the main directory for development purposes.

# !WORK IN PROGRESS!
I'm not finished with v1 of this yet

In fact I haven't added any of the React stuff at this point

## Wishlist
* Refactor the local middleware functions so they are in a helper file
* Add the ability to serve the full HTML for a route from the server side when the user enters the site (or it is crawled)
