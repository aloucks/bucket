// Copyright 2017 bucket Developers
//
// Licensed under the Apache License, Version 2.0, <LICENSE-APACHE or
// http://apache.org/licenses/LICENSE-2.0> or the MIT license <LICENSE-MIT or
// http://opensource.org/licenses/MIT>, at your option. This file may not be
// copied, modified, or distributed except according to those terms.

package cmd

import (
	"os"

	"fmt"
	"strings"

	"github.com/aloucks/bucket"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload DIRECTORY BUCKET",
	Short: "Uploads a local directory to a remote bucket.",
	Long: `Recursively uploads the source DIRECTORY into the destination BUCKET.
	
Uploads new or modified files to an S3 bucket. Modification is detected via a
change in the MD5 checksum. The Content-MD5 header is set on the upload request
and the cheksum is stored in the 'md5chksum' metadata key.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Help()
			os.Exit(-1)
		} else if !strings.HasPrefix(args[1], "s3://") {
			fmt.Fprintln(os.Stderr, "invalid bucket url:", args[1])
			os.Exit(-1)
		} else {
			args[1] = strings.Replace(args[1], "s3://", "", 1)
		}

		svc := NewS3()
		bucket.Upload(svc, args[0], args[1], flags.delete, flags.dryRun)
	},
	Example: "  bucket upload /path/to/my-bucket s3://my-bucket",
}

func init() {
	RootCmd.AddCommand(uploadCmd)
}
