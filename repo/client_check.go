package repo

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/toxyl/devbox/utils"
)

func (c *Client) CheckIfFileIsUpToDate(fileNameLocal, fileNameRemote, storagePath string) (bool, error) {
	fileNameBase := filepath.Base(fileNameLocal)
	fileNameLocal = filepath.Join(storagePath, fileNameBase)
	fileHash := utils.FileToSha256(fileNameLocal)
	fmt.Fprintln(c.conn, "CHECK", filepath.Base(fileNameRemote), fileHash, c.user, c.password)
	return c.isUpToDate()
}

func (c *Client) isUpToDate() (bool, error) {
	// Print server response
	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return false, err
	}
	response = strings.TrimSpace(response)
	if response == "OUTDATED" {
		return false, nil
	}
	if response == "UP-TO-DATE" {
		return true, nil
	}
	if strings.HasPrefix(response, "ERROR") {
		return false, fmt.Errorf(response)
	}
	return true, nil
}
