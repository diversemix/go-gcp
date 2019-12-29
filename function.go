// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"cloud.google.com/go/functions/metadata"
	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	visionapi "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// Global API clients used across function invocations.
var (
	storageClient *storage.Client
	visionClient  *vision.ImageAnnotatorClient
)

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Kind                    string                 `json:"kind"`
	ID                      string                 `json:"id"`
	SelfLink                string                 `json:"selfLink"`
	Name                    string                 `json:"name"`
	Bucket                  string                 `json:"bucket"`
	Generation              string                 `json:"generation"`
	Metageneration          string                 `json:"metageneration"`
	ContentType             string                 `json:"contentType"`
	TimeCreated             time.Time              `json:"timeCreated"`
	Updated                 time.Time              `json:"updated"`
	TemporaryHold           bool                   `json:"temporaryHold"`
	EventBasedHold          bool                   `json:"eventBasedHold"`
	RetentionExpirationTime time.Time              `json:"retentionExpirationTime"`
	StorageClass            string                 `json:"storageClass"`
	TimeStorageClassUpdated time.Time              `json:"timeStorageClassUpdated"`
	Size                    string                 `json:"size"`
	MD5Hash                 string                 `json:"md5Hash"`
	MediaLink               string                 `json:"mediaLink"`
	ContentEncoding         string                 `json:"contentEncoding"`
	ContentDisposition      string                 `json:"contentDisposition"`
	CacheControl            string                 `json:"cacheControl"`
	Metadata                map[string]interface{} `json:"metadata"`
	CRC32C                  string                 `json:"crc32c"`
	ComponentCount          int                    `json:"componentCount"`
	Etag                    string                 `json:"etag"`
	CustomerEncryption      struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	}
	KMSKeyName    string `json:"kmsKeyName"`
	ResourceState string `json:"resourceState"`
}

var _testing = false

// init - Intialise the function on load
func init() {
	if _testing {
		return
	}
	// Declare a separate err variable to avoid shadowing the client variables.
	var err error

	storageClient, err = storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}

	visionClient, err = vision.NewImageAnnotatorClient(context.Background())
	if err != nil {
		log.Fatalf("vision.NewAnnotatorClient: %v", err)
	}
}

// DetectAndCrop - Detect and crop the incoming image.
func DetectAndCrop(ctx context.Context, e GCSEvent) error {
	_, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}

	// check the name of the file is from the right folder
	if !strings.Contains(e.Name, "raw_images") {
		return nil
	}

	// open the image
	imgURI := fmt.Sprintf("gs://%s/%s", e.Bucket, e.Name)
	img := vision.NewImageFromURI(imgURI)
	log.Printf("Found file: %v\n", e.Name)

	annotations, err := visionClient.DetectFaces(ctx, img, nil, 1)

	if err != nil {
		log.Printf("Error processing image: %v\n", err)
		return fmt.Errorf("AnnotateImage: %v", err)
	}

	log.Printf("Number of faces found: %v\n", len(annotations))
	if len(annotations) == 1 {
		log.Printf("Looking for the bounding Rectangle...\n")
		minX, minY, maxX, maxY := findBoundingRect(annotations[0].GetBoundingPoly().GetVertices())
		log.Printf("Found bounds of: (%d, %d) - (%d, %d)\n", minX, minY, maxX, maxY)

		outputName := strings.Replace(e.Name, "raw_images", "processed_images", 1)
		log.Printf("Cropping images from %v to %v\n", e.Name, outputName)
		crop(ctx, e.Bucket, "outgoing-images", e.Name, outputName, minX, minY, maxX, maxY)
	} else {
		log.Printf("ERROR - could not find exactly one face!\n")
	}

	return nil
}

// crop - crops the image to the bounding rect
// inputName e.g-  raw_images/peter/2019-12-17T14:39:18.635Z.png
// outputName e.g. processed_images/peter/2019-12-17T14:39:18.635Z.png

func crop(ctx context.Context, inputBucket, outputBucket, inputName string, outputName string, minX int32, minY int32, maxX int32, maxY int32) error {
	inputBlob := storageClient.Bucket(inputBucket).Object(inputName)
	r, err := inputBlob.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("NewReader: %v", err)
	}

	outputBlob := storageClient.Bucket(outputBucket).Object(outputName)
	w := outputBlob.NewWriter(ctx)
	defer w.Close()

	width := maxX - minX
	height := maxY - minY

	cropArgs := fmt.Sprintf("%dx%d+%d+%d", width, height, minX, minY)

	// Use - as input and output to use stdin and stdout.
	cmd := exec.Command("convert", "-", "-crop", cropArgs, "-")
	cmd.Stdin = r
	cmd.Stdout = w

	if err := cmd.Run(); err != nil {
		log.Printf("ImageMagik failed: %v\n", err)
		return fmt.Errorf("cmd.Run: %v", err)
	}

	log.Printf("Cropped image uploaded to gs://%s/%s | %s\n", outputBlob.BucketName(), outputBlob.ObjectName(), cropArgs)

	return nil
}

func findBoundingRect(vertices []*visionapi.Vertex) (int32, int32, int32, int32) {
	firstVert := vertices[0]

	// found one face get bounding rect
	var (
		minX = firstVert.GetX()
		maxX = firstVert.GetX()
		minY = firstVert.GetY()
		maxY = firstVert.GetY()
	)

	for _, v := range vertices {
		x := v.GetX()
		y := v.GetY()

		if x > maxX {
			maxX = x
		} else if x < minX {
			minX = x
		}

		if y > maxY {
			maxY = y
		} else if y < minY {
			minY = y
		}
	}
	return minX, minY, maxX, maxY
}
