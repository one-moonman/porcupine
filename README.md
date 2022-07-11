# bug-free-octo-broccoli
A restful JWT-authentication service written in Go (Golang). The project is successor to a previous Node.js project of mine. The used stack is the Gin web framework, Redis for storing the tokens and MongoDB for the user details.
* [Login](#Login)
* [Register](#Register)
* [Logout](#Logout)
* [Me](#Me)

### Endpoints and the logic behind them
#### Login
```sh
curl -d '{"email":<example_email>, "password":<example_password>}' -H "Content-Type: application/json" -X POST http://<your_host>:<your_port>/api/v1/login
```
Json body expects `email` and `password` which are checked and then passed to [GenerateTokens](#GenerateTokens) middleware.


#### Register
```sh
curl -d '{"username":<example_username>,"email":<example_email>, "password":<example_password>}' -H "Content-Type: application/json" -X POST http://<your_host>:<your_port>/api/v1/login
```
`username`, `email` and `password` 


#### GenerateTokens