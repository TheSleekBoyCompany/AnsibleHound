package core

import (
	"ansible-hound/core/ansible"
	"fmt"
	"os"
	"time"

	"github.com/TheManticoreProject/gopengraph/edge"
	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/charmbracelet/log"
)

func GenerateEdge(edgeKind string, startId string, endId string, startKind ...string) (e *edge.Edge) {

	e, err := edge.NewEdge(startId, endId, edgeKind, nil)
	if err != nil {
		log.Error(err)
	}

	return e
}

func GenerateNodes[T ansible.AnsibleType](objects map[int]T) (nodes []*node.Node) {
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
