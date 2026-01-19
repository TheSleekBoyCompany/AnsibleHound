package opengraph

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

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
