# bug-free-octo-broccoli
A restful JWT-authentication service written in Go (Golang). The project is successor to a previous Node.js project of mine. The used stack is the Gin web framework, Redis for storing the tokens and MongoDB for the user details.
* [Login](#Login)
* [Register](#Register)
* [Logout](#Logout)
* [Me](#Me)
#### Login
```sh
curl -d '{"email":<example_email>, "password":<example_password>}' -H "Content-Type: application/json" -X POST http://<your_host>:<your_port>/api/v1/login
```
`email` and `password` are expected as input data from the *POST* requests JSON body. The `password` is compared with the hash of the record with the corresponding `email` if it exists at all, then a new pair of tokens with a common `pairId` is generated. The refresh-token is stored in the secondary Redis database for verification purposes with the following JSON: 
```json
<userId>_<pairId>: {
	"expiresAt": <token_expiration_unix_time>,
	"refreshToken": <refresh_token>
}
```
Then the handler responds successfully with the token pair. 
#### Register
```sh
curl -d '{"username":<example_username>,"email":<example_email>, "password":<example_password>}' -H "Content-Type: application/json" -X POST http://<your_host>:<your_port>/api/v1/login
``` 
`username`, `email` and `password` are expected as input data from the *POST* requests JSON body. If a record with the corresponding `username` or `email` does not exist, a new one is created with encrypted `password` as hash. The primary database and the model can be replaced with appropriate ones for a relational database as the [upper/db](https://github.com/upper/db) package allows. The flag `bson` should be replaced with the `db` one. Here are some [instructions ](https://tour.upper.io/welcome/01) for the package use and our model that you may want to change:
```go
type User struct {
	ID 			bson.ObjectId 	`bson:"_id" json:"_id"`
	Username 	string  		`bson:"username" json:"username" binding:"required"`
	Email 		string  		`bson:"email" json:"email" binding:"required"`
	Hash 		string  		`bson:"hash" json:"hash" binding:"required"`
}
```