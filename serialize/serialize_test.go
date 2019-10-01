package serialize

import (
	"log"
	"testing"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func TestSerialize(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Print("Start test serialize")
	var (
		err    error
		result []*pb.EntityAnnotation
	)
	result, err = Deserialize("../example-text-santi-1280.json")
	if err != nil {
		log.Print(err)
	}
	log.Printf("%+v", result)
}
