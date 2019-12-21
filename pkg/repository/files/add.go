package files

import (
	"encoding/hex"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"path/filepath"
	"strconv"
)

func Add(filesDirPath, filePath string, f *File, copyBuf []byte) error {
	hashStr := hex.EncodeToString(f.Hash)
	destination := filepath.Join(filesDirPath, hashStr[:2], hashStr + "-" + strconv.FormatInt(f.Size, 10))

	exists, err := utils.FileExist(destination)
	if err != nil {
		return fmt.Errorf("error checking file existence: %s", err)
	}

	if exists {
		return nil
	}

	if err := utils.CopyFile(filePath, destination, copyBuf); err != nil {
		return fmt.Errorf("error copying file (%s) to repository: %s", filePath, err)
	}
	return nil
}
