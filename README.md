# Bucket

The aws s3 `sync` command does not push files with the Content-MD5 header set in some cases, thus it's possible to corrupt files in transit. The `bucket` tool sets the MD5 header as described in the [aws docs](https://aws.amazon.com/premiumsupport/knowledge-center/data-integrity-s3/) in order to prevent corrupted uploads. An `md5chksum` metadata key is also set to allow for download validation and change detection.

## Install the CLI from local source

    $ go install ./bucket

## Command Line

    $ bucket
    An AWS S3 bucket synchronization tool.

    Bucket mimics the functionality of the the 'aws s3 sync' command, but uses
    MD5 hashing to validate the results of file transfers. All uploaded files
    will have an 'md5chksum' metadata attribute set with the base64 encoded
    checksum.

    Credentials can be specified via the typical AWS environment variables or
    with ~/.aws/credentials.

    AWS_ACCESS_KEY_ID
    AWS_SECRET_ACCESS_KEY

    Usage:
    bucket [command]

    Examples:
    bucket upload /path/to/my-bucket s3://my-bucket
    bucket download s3://my-bucket /path/to/my-bucket

    Available Commands:
    upload      Uploads a local directory to a remote bucket.

    Flags:
        --debug                 debug
    -d, --delete                delete destination files that no longer exist in the source
        --dry-run               perform a simulation without actually making any changes
    -u, --endpoint-url string   aws endpoint url
        --path-style            force path style urls
    -p, --profile string        aws profile in ~/.aws/credentials and ~/.aws/config
    -r, --region string         aws region

    Use "bucket [command] --help" for more information about a command.



## AWS Sync Issues

https://github.com/aws/aws-cli/issues/599


## License

Licensed under either of

 * Apache License, Version 2.0, ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
 * MIT license ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT) at your option.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall be dual licensed as above, without any
additional terms or conditions.
