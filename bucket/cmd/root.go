// Copyright 2017 bucket Developers
//
// Licensed under the Apache License, Version 2.0, <LICENSE-APACHE or
// http://apache.org/licenses/LICENSE-2.0> or the MIT license <LICENSE-MIT or
// http://opensource.org/licenses/MIT>, at your option. This file may not be
// copied, modified, or distributed except according to those terms.

package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cmdFlags struct {
	cfgFile     string
	debug       bool
	endpointURL string
	profile     string
	dryRun      bool
	delete      bool
	pathStyle   bool
	region      string
}

var flags cmdFlags

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "bucket",
	Long: `An AWS S3 bucket synchronization tool.
	
Bucket mimics the functionality of the the 'aws s3 sync' command, but uses 
MD5 hashing to validate the results of file transfers. All uploaded files
will have an 'md5chksum' metadata attribute set with the base64 encoded 
checksum.

Credentials can be specified via the typical AWS environment variables or
with ~/.aws/credentials.

AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
`,
	Example: `   bucket upload /path/to/my-bucket s3://my-bucket
   bucket download s3://my-bucket /path/to/my-bucket`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

// NewS3 creates S3 service instance from the command line flags
func NewS3() *s3.S3 {
	config := aws.NewConfig()
	config.S3ForcePathStyle = &flags.pathStyle
	if len(flags.region) > 0 {
		config = config.WithRegion(flags.region)
	}
	if len(flags.endpointURL) > 0 {
		config = config.WithEndpoint(flags.endpointURL)
	}
	if len(flags.profile) > 0 {
		config = config.WithCredentials(credentials.NewSharedCredentials("", flags.profile))
	}
	if flags.debug {
		config = config.WithLogLevel(aws.LogDebug)
	}

	sess := session.New()
	svc := s3.New(sess, config)
	return svc
}

func init() {
	cobra.OnInitialize(initConfig)

	// RootCmd.PersistentFlags().StringVar(&flags.cfgFile, "config", "", "config file (default is $HOME/.bucket.yaml)")
	RootCmd.PersistentFlags().BoolVar(&flags.dryRun, "dry-run", false, "perform a simulation without actually making any changes")
	RootCmd.PersistentFlags().StringVarP(&flags.region, "region", "r", "", "aws region")
	RootCmd.PersistentFlags().StringVarP(&flags.profile, "profile", "p", "", "aws profile in ~/.aws/credentials and ~/.aws/config")
	RootCmd.PersistentFlags().StringVarP(&flags.endpointURL, "endpoint-url", "u", "", "aws endpoint url")
	RootCmd.PersistentFlags().BoolVarP(&flags.delete, "delete", "d", false, "delete destination files that no longer exist in the source")
	RootCmd.PersistentFlags().BoolVar(&flags.debug, "debug", false, "debug")
	RootCmd.PersistentFlags().BoolVar(&flags.pathStyle, "path-style", false, "force path style urls")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if flags.cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(flags.cfgFile)
	}

	viper.SetConfigName(".bucket") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
