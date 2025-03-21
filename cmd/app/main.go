package main

import (
	"fmt"
	"os"

	"github.com/y-yoshida/gcp-lab/internal/app"
	"github.com/y-yoshida/gcp-lab/internal/app/gcs"
)

func main() {
	app := application()
	output, err := app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(output)
}

func application() *app.App {
	bucket := os.Getenv("GCS_BUCKET_FOR_UPLOAD")
	return app.New(gcs.NewSignedURLGenerator(bucket))
}
