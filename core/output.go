package core

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func GenerateEdge(kind string, startId string, endId string) Edge {

	start := Link{
		Value:   startId,
		MatchBy: "id",
	}

	end := Link{
		Value:   endId,
		MatchBy: "id",
	}

	edge := Edge{
		Kind:  kind,
		Start: start,
		End:   end,
	}

	return edge
}

func GenerateNodes[T AnsibleType](objects []T) (nodes []Node) {
	for _, object := range objects {
		nodes = append(nodes, object.ToBHNode())
	}
	return nodes
}

func CalculateName(objectType string) string {
	now := time.Now()
	epoch := now.Unix()
	return fmt.Sprintf("%d_%s.json", epoch, objectType)
}

func WriteToFile(content []byte, filePath string) error {

	log.Debug(fmt.Sprintf("Writing to file `%s`.", filePath))

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}
