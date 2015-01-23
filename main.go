package main

import (
	"compress/gzip"
	"crypto/sha1"
	"github.com/codegangsta/cli"

	"fmt"
	"io"
	"os"
)

const ChunkSize = 1024 * 1024

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

func main() {
	app := cli.NewApp()
	app.Name = "gochunk"
	app.Usage = "split a large file into smaller files named by the sha1's and compressed"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "dir, d",
			Value:  "./",
			Usage:  "directory to read/write chunks from",
			EnvVar: "GOCHUNK_DIR",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "chop",
			ShortName: "c",
			Action: func(c *cli.Context) {
				for _, fileName := range c.Args() {
					fp, err := os.Open(fileName)
					if err != nil {
						fmt.Println(err)
						return
					}
					defer fp.Close()

					buffer := make([]byte, ChunkSize)
					for {
						n, err := fp.Read(buffer)
						if err != nil && err != io.EOF {
							panic(err)
						}
						if n == 0 {
							break
						}
						sha1sum := sha1.Sum(buffer[:n])
						sha1dir := fmt.Sprintf("%02x/%02x/%02x", sha1sum[0], sha1sum[1], sha1sum[2])
						werr := WriteChunk(buffer, n, sha1dir, sha1sum)
						fmt.Printf("%x %s\n", sha1sum, fileName)
						if werr != nil {
							panic(werr)
						}

						if err == io.EOF {
							break
						}
					}
				}
			},
		},
	}
	app.Run(os.Args)
}
