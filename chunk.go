package main

const ChunkSize = 1024 * 1024

type ChunkWriteReq struct {
	buffer  []byte
	n       int
	sha1dir string
	sha1sum [sha1.Size]byte
}

// Writer interface
func (c *Chunk) Write(p []byte) (int, error) {
	return 0, nil
}

func (c *Chunk) Close() error {
}

func WriteChunk(buffer []byte, n int, sha1dir string, sha1sum [sha1.Size]byte) error {
	err := os.MkdirAll(sha1dir, 0777)
	if err != nil {
		return err
	}
	chunkFile := fmt.Sprintf("%s/%x", sha1dir, sha1sum)
	chunkTemp := fmt.Sprintf("%s.%d", chunkFile, os.Getpid())
	fp0, err := os.Create(chunkTemp)
	if err != nil {
		return err
	}

	fp := gzip.NewWriter(fp0)

	ofs, remain := 0, n
	for remain > 0 {
		written, err := fp.Write(buffer[ofs:remain])
		if err != nil {
			fp.Close()
			return err
		}
		remain -= written
		ofs += written
	}

	fp.Close()
	err = os.Rename(chunkTemp, chunkFile)
	return err
}
