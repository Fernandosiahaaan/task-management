package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB_NAME          = "log-service-db"
	UserCollectionDB = "user-log"
	TaskCollectionDB = "task-log"
)

type RepoMongo struct {
	ctx            context.Context
	cancel         context.CancelFunc
	Conn           *mongo.Client
	userCollection *mongo.Collection
	taskCollection *mongo.Collection
}

func Init(ctx context.Context) (*RepoMongo, error) {
	mongoCtx, mongoCancel := context.WithCancel(ctx)
	urlMongoDB := os.Getenv("MONGODB_URL")
	if urlMongoDB == "" {
		return nil, fmt.Errorf("not found url mongo DB in .env")
	}

	fmt.Println("urlMongoDB = ", urlMongoDB)

	// Membuat client MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(urlMongoDB))
	err = client.Connect(mongoCtx)
	if err != nil {
		return nil, fmt.Errorf("failed connect to mongodb. err: %v", err)
	}

	// Memeriksa koneksi MongoDB
	err = client.Ping(mongoCtx, nil) // Menggunakan ctx yang sama
	if err != nil {
		return nil, fmt.Errorf("failed ping to mongodb. err: %v", err)
	}

	userCollection := client.Database(DB_NAME).Collection(UserCollectionDB)
	taskCollection := client.Database(DB_NAME).Collection(TaskCollectionDB)

	if err = addUserActionValidator(mongoCtx, client, DB_NAME, UserCollectionDB); err != nil {
		return nil, fmt.Errorf("failed to add validator to collection %s. err = %v", UserCollectionDB, err)
	}

	if err = addTaskCollectionValidator(mongoCtx, client, DB_NAME, TaskCollectionDB); err != nil {
		return nil, fmt.Errorf("failed to add validator to collection %s. err = %v", TaskCollectionDB, err)
	}

	return &RepoMongo{
		ctx:            mongoCtx,
		cancel:         mongoCancel,
		Conn:           client,
		userCollection: userCollection,
		taskCollection: taskCollection,
	}, nil
}

// Contoh insert dokumen ke MongoDB
func (r *RepoMongo) InsertExample() {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	// Data yang akan diinsert
	user := bson.D{
		{Key: "name", Value: "John Doe"},
		{Key: "age", Value: 30},
		{Key: "email", Value: "johndoe@example.com"},
	}

	// Insert ke koleksi
	result, err := r.userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dokumen berhasil diinsert dengan ID:", result.InsertedID)
}

// Contoh query sederhana untuk mencari dokumen
func (r *RepoMongo) FindExample() {
	collection := r.Conn.Database(DB_NAME).Collection("users")

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	// Mencari dokumen dengan nama "John Doe"
	filter := bson.D{{Key: "name", Value: "John Doe"}}

	// Menjalankan query find
	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dokumen ditemukan:", result)
}

// Contoh query sederhana untuk mencari dokumen
func (r *RepoMongo) InsertUserLog(input primitive.M) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	// Menjalankan query find
	return r.userCollection.InsertOne(ctx, input)
}

// Contoh query sederhana untuk mencari dokumen
func (r *RepoMongo) InsertTaskLog(input primitive.M) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	// Menjalankan query find
	return r.taskCollection.InsertOne(ctx, input)
}

func (r *RepoMongo) Close() {
	r.cancel()
	r.Conn.Disconnect(r.ctx)
}
