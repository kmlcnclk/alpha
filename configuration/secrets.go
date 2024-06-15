package configuration

var JWT_SECRET = "my_secret_key_that_is_much_longer_and_more_complex_1234567890"
var REFRESH_SECRET = "my_refresh_secret_key_that_is_also_much_longer_and_more_complex_1234567890"

var MONGO_URI = "mongodb+srv://123asd123:123asd123@a.psdrs.mongodb.net/alpha-db?retryWrites=true&w=majority&appName=A"

var MONGO_DB_NAME = "alpha-db"
var MONGO_USERS_DB_NAME = "users"
var MONGO_JWT_DB_NAME = "jwt"

var ACCESS_TOKEN_TIME = "7d"
var REFRESH_TOKEN_TIME = "21d"
