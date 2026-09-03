package config

import (
	"log"
	"os"

	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

var IK *imagekit.Client

func InitImageKit() {
	privateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("IMAGEKIT_PRIVATE_KEY is not set in .env")
	}

	client := imagekit.NewClient(
		option.WithPrivateKey(privateKey),
	)

	IK = &client
	log.Println("ImageKit Client Initialized Successfully")
}