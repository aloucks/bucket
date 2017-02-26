// Copyright 2017 bucket Developers
//
// Licensed under the Apache License, Version 2.0, <LICENSE-APACHE or
// http://apache.org/licenses/LICENSE-2.0> or the MIT license <LICENSE-MIT or
// http://opensource.org/licenses/MIT>, at your option. This file may not be
// copied, modified, or distributed except according to those terms.

package bucket

import (
	"crypto/md5"

	"encoding/base64"
	"encoding/hex"

	"fmt"
	"io"
	"os"
	"strings"

	"path"
	"path/filepath"

	"mime"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func isNoSuchKey(err error) bool {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return true
			case "NotFound":
				return true
			}
			check(err)
		}
	}
	return false
}

type maybeUpload struct {
	absFilePath string
	md5sum      []byte
	fileSize    int64
	key         string
	contentType string
}

func (mu *maybeUpload) Md5hex() string {
	return hex.EncodeToString(mu.md5sum)
}

func (mu *maybeUpload) Md5base64() string {
	return base64.StdEncoding.EncodeToString(mu.md5sum)
}

func contentType(filename string) string {
	ext := path.Ext(filename)
	contentType := mime.TypeByExtension(ext)
	if len(contentType) == 0 {
		contentType = "binary/octet-stream"
	}
	return contentType
}

// Upload a directory to an S3 bucket. If delete is set to true, objects no longer
// in the source directory will be deleted from the destination bucket.
func Upload(svc *s3.S3, sourceDir string, destBucket string, delete bool, dryRun bool) {
	absSourceDir, err := filepath.Abs(sourceDir)
	check(err)

	workerCount := 4
	maybeUploads := make(chan maybeUpload, workerCount)
	done := make(chan bool, 1)

	buffer := make([]byte, 8192)
	walk := func(filePath string, info os.FileInfo, err error) error {
		check(err)
		if !info.IsDir() {
			absFilePath, err := filepath.Abs(filePath)
			check(err)
			digest := md5.New()
			file, err := os.Open(filePath)
			check(err)
			defer file.Close()
			for {
				read, err := file.Read(buffer)
				if err == io.EOF {
					break
				} else {
					check(err)
					digest.Write(buffer[0:read])
				}
			}
			md5sum := make([]byte, 0)
			md5sum = digest.Sum(md5sum)
			fileSize := info.Size()
			contentType := contentType(absFilePath)
			key, err := filepath.Rel(absSourceDir, absFilePath)
			check(err)
			key = strings.Replace(key, "\\", "/", -1)
			maybeUploads <- maybeUpload{absFilePath, md5sum, fileSize, key, contentType}
		}
		return err
	}

	upload := func(maybeUploads <-chan maybeUpload) {
		for maybeUpload := range maybeUploads {
			fileMd5chksum := maybeUpload.Md5base64()
			objMd5chksum := (*string)(nil)
			res, err := svc.HeadObject(&s3.HeadObjectInput{
				Bucket: &destBucket,
				Key:    &maybeUpload.key})
			if !isNoSuchKey(err) {
				objMd5chksum = res.Metadata["Md5chksum"]
			}
			if objMd5chksum == nil || *objMd5chksum != fileMd5chksum {
				file, err := os.Open(maybeUpload.absFilePath)
				check(err)
				defer file.Close()
				req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
					Bucket:        &destBucket,
					Key:           &maybeUpload.key,
					ContentLength: &maybeUpload.fileSize,
					ContentType:   &maybeUpload.contentType,
					Body:          file})
				req.HTTPRequest.Header.Add("x-amz-meta-md5chksum", fileMd5chksum)
				req.HTTPRequest.Header.Add("content-md5", fileMd5chksum)

				if !dryRun {
					err := req.Send()
					check(err)
				}

				fmt.Println("+", maybeUpload.key)
			} else {
				fmt.Println(" ", maybeUpload.key)
			}
		}
		done <- true
	}

	go func() {
		err := filepath.Walk(absSourceDir, walk)
		close(maybeUploads)
		check(err)
	}()

	for w := 0; w < workerCount; w++ {
		go upload(maybeUploads)
	}

	for w := 0; w < workerCount; w++ {
		<-done
	}

}
