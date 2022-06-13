# Clean Architecture in Go
### In this basic application it is intended to apply the principles of clean architecture
The project is structured based on the dependency rule which means that the dependencies of the source code must always point inwards.
![alt text](https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg)
[The Clean Code Blog](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

# Domain
In this layer the models folder has been created which encapsulates the commercial rules and here the custom validations and errors are also defined
# Infrastructure
For this layer, the infrastructure folder has been created, which is related to frameworks, databases, among others. Changes made to this layer should not affect the inner layers.
# Interface
For this layer, the interface folder has been created that serves as a translator between the infrastructure layer and the domain since it packs the input and output data as required.
# Use Cases
The app folder has been created for this layer, in which the commercial rules of the application are defined.

### Installation:
To run our application we execute the following command `docker-compose up --build`, you must have previously installed docker.
# END-POINTS
### Sign up
```
curl -X POST http://localhost:8080/api/v1/auth/signup -H 'Content-Type: application/json' -d '{ "name": "test", "lastname": "test", "email": "test@gmail.com","password": "test123232"}'
```
### Sign in
```
curl -X POST http://localhost:8080/api/v1/auth/login -H 'Content-Type: application/json' -d '{ "email": "test@gmail.com", "password": "123456789" }'
```
### With this project we have reviewed the following:
1. Clean Architecture.
2. Unit tests with testify and go-sqlmock.
3. Using interfaces in golang.
4. Using gin-gonic for http routing, middleware, logger (zerolog) and recover
5. The GORM ORM for postgres database
6. For authentication we use jwt
7. We have used docker and docker-compose to manage the containers of our application.