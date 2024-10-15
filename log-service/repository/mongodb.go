package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepoMongo struct {
	ctx    context.Context
	cancel context.CancelFunc
	Conn   *mongo.Client
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

	return &RepoMongo{
		ctx:    mongoCtx,
		cancel: mongoCancel,
		Conn:   client,
	}, nil
}

// Contoh insert dokumen ke MongoDB
func (r *RepoMongo) InsertExample() {
	collection := r.Conn.Database("mydb").Collection("users")

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	// Data yang akan diinsert
	user := bson.D{
		{Key: "name", Value: "John Doe"},
		{Key: "age", Value: 30},
		{Key: "email", Value: "johndoe@example.com"},
	}

	// Insert ke koleksi
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dokumen berhasil diinsert dengan ID:", result.InsertedID)
}

// Contoh query sederhana untuk mencari dokumen
func (r *RepoMongo) FindExample() {
	collection := r.Conn.Database("mydb").Collection("users")

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

func (r *RepoMongo) Close() {
	r.cancel()
	r.Conn.Disconnect(r.ctx)
}
