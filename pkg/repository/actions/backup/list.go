package backup

import (
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/files"
	"github.com/Miguel-Dorta/gkup-backend/pkg/utils"
	"io"
	"os"
	pathPkg "path"
	"path/filepath"
)

func listFiles(paths []string, errWriter io.Writer) []*file {
	list := make([]*file, 0, len(paths))

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			printError(errWriter, "error getting info from path \"%s\": %s", path, err)
			continue
		}

		if stat.IsDir() {
			list = append(list, listFilesRecursive(path, stat.Name(), errWriter)...)
		} else if stat.Mode().IsRegular() {
			list = append(list, &file{
				RealPath: path,
				File: files.File{
					RelativePath: stat.Name(),
					Hash:         nil,
					Size:         stat.Size(),
				},
			})
		} else {
			printError(errWriter, "file type of \"%s\" is not supported", path) //TODO
		}
	}

	return list
}

func listFilesRecursive(pathReal, pathRelative string, errWriter io.Writer) []*file {
	list, err := utils.ListDir(pathReal)
	if err != nil {
		printError(errWriter, "error listing directory \"%s\": %s", pathReal, err)
		return nil
	}

	fileList := make([]*file, 0, len(list))
	for _, f := range list {
		fPathReal := filepath.Join(pathReal, f.Name())
		fPathRelative := pathPkg.Join(pathRelative, f.Name())

		if f.IsDir() {
			fileList = append(fileList, listFilesRecursive(fPathReal, fPathRelative, errWriter)...)
		} else if f.Mode().IsRegular() {
			fileList = append(fileList, &file{
				RealPath: fPathReal,
				File: files.File{
					RelativePath: fPathRelative,
					Hash:         nil,
					Size:         f.Size(),
				},
			})
		} else {
			printError(errWriter, "file type of \"%s\" is not supported", fPathReal) //TODO
		}
	}

	return fileList
}
