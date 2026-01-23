package opengraph

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/TheManticoreProject/gopengraph"
	"github.com/charmbracelet/log"
)

func calculateName(objectType string) string {
	now := time.Now()
	epoch := now.Unix()
	return fmt.Sprintf("%d_%s.json", epoch, objectType)
}

func writeToFile(content []byte, filePath string) error {

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

func Output(graph *gopengraph.OpenGraph, outdir string) {

	outputJson, err := graph.ExportJSON(false)
	err = writeToFile(
		[]byte(outputJson),
		path.Join(outdir, calculateName("output")))
	if err != nil {
		log.Fatal(err)
	}

}
