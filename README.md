# chiact

The goal of this project is to create a boilerplate for developing full-stack web applications using [Golang](https://golang.org) (Go) on the server side and [React](https://facebook.github.io/react/) on the client side.

chiact includes:
* Server side routing with [Chi](https://github.com/pressly/chi)
* Database connection
* Static directory serving
* API structure and starter methods
* Authentication and JWT
* User Storage and editing

# Please Note
I am working on this in order to learn React and Go. I am not seasoned with either language and so I am sure there are many things I could be doing differently/better. I welcome any pull requests to improve or add to this code.

### Instructions
1. cd $GOHOME/src/github.com/your-github-username
2. Clone Repo - git clone https://github.com/tutley/chiact myprojectname
3. Search and replace in the directory for tutley/chiact, replace with your-github-username/myprojectname
     1. (you will need to edit the name chiact in other places for display purposes)
4. glide install
5. npm install
6. npm start
7. http://localhost:3333/ - Create an account and login to test that it's working

To build for production: npm run production

### Go Structure
The server side, written with Go, uses glide for vendoring (dependency management) and can be run locally with "go run main.go". The glide.yaml file holds the dependencies, which can be installed with 'glide install'.

* The handlers directory contains the various local packages for handling routes.
* The helpers directory contains packages with functions to assist along the way.
* The models directory contains data models with methods included.
* The vendor directory contains the local copies of dependencies.

We have included mgo for storing data in MongoDB but if you prefer another database that could be changed by editing the models as well as the db init in main and helpers.

### React Structure
The client application is written in React, which requires a few files in the main directory for development purposes. The files that will be served for the client are in the client directory. The development is done on the files in the client/src directory, then webpack builds the final version.

These files are used with React:
* package.json - contains dependency configurations (use npm install to set it up)
* node_modules - this directory will be created by npm install to hold local dependencies
* webpack.config.js - this file is used to setup the dev environment for react
* client/js ; client/img ; client/css - these static directories serve assets
* client/src - react development directory
* client/index.html - this is generated by webpack when you build the client
* client/index_bundle.js - this is generated by webpack when you build the client

## Wishlist
* Refactor the local middleware functions so they are in a helper file
* A simple blog
* Navigation Bar
* Add the ability to serve the full HTML for a public route from the server side when the user enters the site (or it is crawled)

## Sources
I used the following sources for example code in making this project:
* [authentication-in-react-apps](https://github.com/vladimirponomarev/authentication-in-react-apps)
* [stack](https://github.com/bradialabs/stack)
