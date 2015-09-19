package main

import (
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/gopkg.in/amz.v1/aws"
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/gopkg.in/amz.v1/s3"
)

var (
	auth = aws.Auth{
		AccessKey: ConS3Key,
		SecretKey: ConS3Secret,
	}
	useast     = aws.USEast
	connection = s3.New(auth, useast)
	mybucket   = connection.Bucket(ConS3Bucket)
)

func SaveToS3(image []byte, name string) error {
	err := mybucket.Put(name+".png", image, "image/png", s3.PublicRead)
	if err != nil {
		LogError(err.Error())
	}
	return err
}

func GetFromS3(path string) ([]byte, error) {
	return mybucket.Get(path)
}
