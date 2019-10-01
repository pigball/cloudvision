package serialize

import (
	"bytes"
	"encoding/json"
	"os"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func Deserialize(file string) ([]*pb.EntityAnnotation, error) {
	var (
		result []*pb.EntityAnnotation
	)

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(f)

	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return nil, err
	}

	return result, nil
}
