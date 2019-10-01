package cloudvision

import (
	"log"
	"testing"

	"github.com/pongsanti/cloudvision/serialize"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func TestNewDocument(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Print("Start test new document")
	var (
		err    error
		result []*pb.EntityAnnotation
	)
	result, err = serialize.Deserialize("example-kyaw-zin-oo.json")
	if err != nil {
		log.Print(err)
	}
	d := NewDocument(result)
	log.Printf("%+v", d.Parse())
}
