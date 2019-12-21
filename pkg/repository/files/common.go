package files

const FolderName = "files"

type File struct {
	RelativePath string
	Hash []byte
	Size int64
}
