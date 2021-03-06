package main

import (
	"crypto/sha1"
	"github.com/codegangsta/cli"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var (
	chunkDir = "./chunks/"
)

func main() {
	app := cli.NewApp()
	app.Name = "gochunk"
	app.Usage = "split a large file into smaller files named by the sha1's and compressed"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "dir, d",
			Value:  chunkDir,
			Usage:  "directory to read/write chunks from",
			EnvVar: "GOCHUNK_DIR",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "serve",
			ShortName: "s",
			Action: func(c *cli.Context) {
				// fs := http.FileServer(http.Dir(chunkDir))
				// http.Handle("/", fs)
				// log.Fatal(http.ListenAndServe(":3000", nil))
				serveDir, err := filepath.Abs(chunkDir)
				if err != nil {
					panic(err)
				}
				http.Handle("/chunks/", http.StripPrefix("/chunks/",
					http.FileServer(http.Dir(serveDir))))
				err = http.ListenAndServe(":9999", nil)
				if nil != err {
					panic(err)
				}
			},
		},
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
					wq := make(chan *ChunkWriteReq, 5)
					done := make(chan bool)

					go func() {
						for {
							select {
							case _, cwr_ok := <-wq:
								if !cwr_ok {
									return // channel closed unexpectedly
								}
							case donesig, done_ok := <-done:
								if done_ok && donesig {
									log.Println("All done")
									return
								}
							}
						}
					}()

					for {
						n, err := fp.Read(buffer)
						if err != nil && err != io.EOF {
							panic(err)
						}
						if n == 0 {
							break
						}
						sha1sum := sha1.Sum(buffer[:n])
						sha1dir := path.Join(
							chunkDir,
							fmt.Sprintf("%02x/%02x/%02x", sha1sum[0], sha1sum[1], sha1sum[2]))
						_ = WriteChunk(buffer, n, sha1dir, sha1sum)

						wq <- &ChunkWriteReq{
							n:       n,
							buffer:  buffer,
							sha1sum: sha1sum,
							sha1dir: sha1dir,
						}

						// ChunkWriteReqHandler(wq, done, errout)

						fmt.Printf("%x\t%s\n", sha1sum, fileName)

						if err == io.EOF {
							break
						}
					}
					done <- true
				}
			},
		},
	}
	app.Run(os.Args)
}
