package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "emberctl",
		Usage: "CLI client for EmberDB",
		Commands: []*cli.Command{
			{
				Name:      "file",
				Usage:     "Upload a file to EmberDB",
				ArgsUsage: "<key> <file-path>",
				Action:    uploadFile,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func uploadFile(c *cli.Context) error {
	if c.Args().Len() != 2 {
		return cli.Exit("usage: emberctl file <key> <file-path>", 2)
	}

	key := c.Args().Get(0)
	path := c.Args().Get(1)

	file, err := os.Open(path)
	if err != nil {
		return cli.Exit(fmt.Sprintf("cannot open file: %v", err), 1)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	if _, err := io.Copy(part, file); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	writer.Close()

	url := fmt.Sprintf("http://localhost:9182/upload/%s", key)

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return cli.Exit(fmt.Sprintf("request failed: %v", err), 1)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return cli.Exit(
			fmt.Sprintf("server error (%d): %s", resp.StatusCode, respBody),
			resp.StatusCode,
		)
	}

	fmt.Println(string(respBody))
	return nil
}
