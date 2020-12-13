package main

import (
    "context"
    _ "github.com/akurczyk/golang_echo_blogging_platform/003_mongo_and_redis/app/src/docs"
    "github.com/go-playground/validator"
    "github.com/go-redis/redis/v8"
    "github.com/labstack/echo-contrib/prometheus"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/swaggo/echo-swagger"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readconcern"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    "go.mongodb.org/mongo-driver/mongo/writeconcern"
    "os"
    "strconv"
    "time"
)

type (
    CustomValidator struct {
        validator *validator.Validate
    }
)

var (
    redisCtx        context.Context
    redisClient     *redis.Client
    mongoCtx        context.Context
    mongoClient     *mongo.Client
    mongoDatabase   *mongo.Database
    usersCollection *mongo.Collection
    postsCollection *mongo.Collection
)

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

func setupMongo() {
    var err error

    connectionString := os.Getenv("BP_MONGO_CONNECTION_STRING")
    username := os.Getenv("BP_MONGO_USERNAME")
    password := os.Getenv("BP_MONGO_PASSWORD")
    database := os.Getenv("BP_MONGO_DATABASE")

    credential := options.Credential{Username: username, Password: password}
    clientOptions := options.Client().ApplyURI(connectionString).SetAuth(credential)

    // Create context - a timeout can be specified here
    mongoCtx = context.Background()

    // Connect
    mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        panic(err)
    }

    // Check connection
    err = mongoClient.Ping(context.TODO(), nil)
    if err != nil {
        panic(err)
    }

    // Create database handle
    mongoDatabase = mongoClient.Database(database)

    // Create users collection handle
    usersOptions := options.CollectionOptions{
        ReadConcern: readconcern.Local(),
        WriteConcern: writeconcern.New(
            writeconcern.WMajority(),
            writeconcern.J(true),
            writeconcern.WTimeout(time.Duration(30 * time.Second)),
        ),
        ReadPreference: readpref.Nearest(),
    }
    usersCollection = mongoDatabase.Collection("users", &usersOptions)

    // Create posts collection handle
    postsOptions := options.CollectionOptions{
        ReadConcern: readconcern.Local(),
        WriteConcern: writeconcern.New(
            writeconcern.W(1),
            writeconcern.J(false),
        ),
        ReadPreference: readpref.Nearest(),
    }
    postsCollection = mongoDatabase.Collection("posts", &postsOptions)

    // Setup users unique indexes
    _, _ = usersCollection.Indexes().CreateMany(
        mongoCtx,
        []mongo.IndexModel{
            {Keys: bson.D{{Key: "name",  Value: 1}}, Options: options.Index().SetUnique(true)},
            {Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
        },
    )
}

func setupRedis() {
    connectionString := os.Getenv("BP_REDIS_CONNECTION_STRING")
    password := os.Getenv("BP_REDIS_PASSWORD")
    db, _ := strconv.Atoi(os.Getenv("BP_REDIS_DATABASE"))

    redisCtx = context.Background()

    redisClient = redis.NewClient(&redis.Options{
        Addr:     connectionString,
        Password: password,
        DB:       db,
    })
}

// @title Simple blogging platform API based on MongoDB and Redis databases
// @version 1.0
// @description Simple blogging platform API created in Golang with the use of Echo framework, MongoDB database for
// @description storing objects, and Redis for storing session data with integrated Prometheus and Swagger, Dockerized.

// @contact.name Aleksander Kurczyk
// @contact.url http://github.com/akurczyk
// @contact.email akurczyk000@gmail.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
    setupMongo()
    setupRedis()

    e := echo.New()

    // Validator
    e.Validator = &CustomValidator{validator: validator.New()}

    // Prometheus
    p := prometheus.NewPrometheus("echo", nil)
    p.Use(e)

    // Swagger
    e.GET("/swagger/*", echoSwagger.WrapHandler)

    e.GET("/users", listUserAccounts)
    e.POST("/users", createUserAccount)
    e.GET("/users/:id", retrieveUserAccount)
    e.PUT("/users", updateUserAccount, middleware.KeyAuth(checkAuthToken))
    e.DELETE("/users", deleteUserAccount, middleware.KeyAuth(checkAuthToken))

    e.POST("/token", issueAuthToken)
    e.DELETE("/token", revokeAuthToken, middleware.KeyAuth(checkAuthToken))

    e.GET("/posts", listPosts)
    e.POST("/posts", createPost, middleware.KeyAuth(checkAuthToken))
    e.GET("/posts/:id", retrievePost)
    e.PUT("/posts/:id", updatePost, middleware.KeyAuth(checkAuthToken))
    e.DELETE("/posts/:id", deletePost, middleware.KeyAuth(checkAuthToken))

    e.POST("/posts/:post_id/comments", createComment, middleware.KeyAuth(checkAuthToken))
    e.PUT("/posts/:post_id/comments/:comment_id", updateComment, middleware.KeyAuth(checkAuthToken))
    e.DELETE("/posts/:post_id/comments/:comment_id", deleteComment, middleware.KeyAuth(checkAuthToken))

    e.Logger.Fatal(e.Start(":1323"))
}

// TODO: Add pagination
// TODO: Write tests
// TODO: Auto update of authors of posts and comments during update of users
// TODO: Add Sentinel support for Redis
// TODO: Add read from slave support for Redis
