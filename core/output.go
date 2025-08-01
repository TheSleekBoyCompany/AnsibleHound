package core

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func OutputBH_Edge(kind string, startId string, endId string) Edge {

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

// TODO: Must be a cleaner way than this big switch case.
// PRs WELCOME *wink* *wink*

func OutputBH_Node[T AnsibleTypeList](objectLists T) []Node {

	nodes := []Node{}

	switch objects := any(objectLists).(type) {
	case []User:
		for i, user := range objects {
			user.UUID = uuid.NewString()
			objects[i] = user
			node := user.ToBHNode()
			nodes = append(nodes, node)
		}

	case []Organization:
		for i, org := range objects {
			org.UUID = uuid.NewString()
			objects[i] = org
			node := org.ToBHNode()
			nodes = append(nodes, node)
		}

	case []JobTemplate:
		for i, jobTemplate := range objects {
			jobTemplate.UUID = uuid.NewString()
			objects[i] = jobTemplate
			node := jobTemplate.ToBHNode()
			nodes = append(nodes, node)
		}

	case []Job:
		for i, job := range objects {
			job.UUID = uuid.NewString()
			objects[i] = job
			node := job.ToBHNode()
			nodes = append(nodes, node)
		}

	case []Project:
		for i, project := range objects {
			project.UUID = uuid.NewString()
			objects[i] = project
			node := project.ToBHNode()
			nodes = append(nodes, node)
		}

	case []Credential:
		for i, credential := range objects {
			credential.UUID = uuid.NewString()
			objects[i] = credential
			node := credential.ToBHNode()
			nodes = append(nodes, node)
		}

	case []Inventory:
		for i, inventory := range objects {
			inventory.UUID = uuid.NewString()
			objects[i] = inventory
			node := inventory.ToBHNode()
			nodes = append(nodes, node)
		}

	case []Host:
		for i, host := range objects {
			host.UUID = uuid.NewString()
			objects[i] = host
			node := host.ToBHNode()
			nodes = append(nodes, node)
		}

	case []Team:
		for i, team := range objects {
			team.UUID = uuid.NewString()
			objects[i] = team
			node := team.ToBHNode()
			nodes = append(nodes, node)
		}
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
