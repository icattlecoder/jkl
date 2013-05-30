**jkl** is a static site generator written in [Go](http://www.golang.org),
based on [Jekyll](https://github.com/mojombo/jekyll)

[![Build Status](https://drone.io/github.com/why404/jkl/status.png)](https://drone.io/github.com/why404/jkl/latest)

Notable similarities between jkl and Jekyll:

* Directory structure
* Use of YAML front matter in Pages and Posts
* Availability of `site`, `content`, `page` and `posts` variables in templates
* Copies all static files into destination directory

Notable differences between jkl and Jekyll:

* Uses [Go templates](http://www.golang.org/pkg/text/template)
* Only supports YAML front matter in markup files
* No plugin support

Additional features:

* Deploy to [Qiniu Cloud Storage](http://www.qiniu.com/)
* Deploy to [Amazon S3](http://aws.amazon.com/s3/)

Sites built with jkl:

* Drone.io Documentation: <http://docs.drone.io>
* Drone.io Blog: <http://blog.drone.io>
* Qiniu Cloud Storage Documentation: <http://docs.qiniu.com>

--------------------------------------------------------------------------------

### Installation

If you want to install from source you should [install the Go tools](http://golang.org/doc/install) first.

In order to compile with `go build` you will first need to download
the following dependencies(be sure you have already set the `GOPATH` environment variable):

```
go get github.com/qiniu/bytes
go get github.com/qiniu/rpc
go get github.com/qiniu/api
go get github.com/russross/blackfriday
go get github.com/howeyc/fsnotify
go get launchpad.net/goyaml
go get launchpad.net/goamz/aws
go get launchpad.net/goamz/s3

go build
```
Once you have compiled `jkl` you can install with the following command:

```
sudo mv ./jkl /usr/local/bin
sudo chmod +x /usr/local/bin/jkl
```

### Usage

```
Usage: jkl [OPTION]... [SOURCE]

      --auto                re-generates the site when files are modified
      --base-url            serve website from a given base URL
      --source              changes the dir where Jekyll will look to transform files
      --destination         changes the dir where Jekyll will write files to
      --server              starts a server that will host your _site directory
      --port                changes the port that the Jekyll server will run on
      --s3                  copies the _site directory to s3
      --s3-key              aws access key use for s3 authentication
      --s3-secret           aws secret key use for s3 authentication
      --s3-bucket           name of the s3 bucket
      --qiniu               copies the _site directory to Qiniu Cloud Storage
      --qiniu-key           access key use for qiniu authentication
      --qiniu-secret        secret key use for qiniu authentication
      --qiniu-bucket        name of the qiniu bucket

  -v, --verbose             runs Jekyll with verbose output
  -h, --help                display this help and exit

Examples:
  jkl                       generates site from current working directory
  jkl /path/to/site         generates site from source dir /path/to/site
  jkl --server              generates site and serves at localhost:4000
  jkl --server --port=:4567 generates site and serves at localhost:4567
  jkl --s3 --verbose        copies the _site directory to s3
  jkl --qiniu --verbose     copies the _site directory to Qiniu Cloud Storage
```

### Auto Generation

If you are running the website in server mode, with the `--server` flag, you can
also instruct `jkl` to auto-recompile you website by adding the `--auto` flag.

### Deploy to Qiniu Cloud Storage

In order to deploy to [Qiniu Cloud Storage](http://www.qiniu.com/) you must include a `_jekyll_qiniu.yml` file in your
site's root directory that specifies your Qiniu key, secret and bucket:

```
access_key: YOUR_QINIU_ACCESS_KEY
secret_key: YOUR_QINIU_SECRET_KEY
bucket: YOUR_QINIU_BUCKET
```

Run 

```
jkl --qiniu
```

### Deploy to Amazon S3

If you want to deploy to [Amazon S3](http://aws.amazon.com/s3/) you must include a `_jekyll_s3.yml` file in your
site's root directory that specifies your AWS key, secret and bucket:

```
s3_id: YOUR_AWS_S3_ACCESS_KEY_ID
s3_secret: YOUR_AWS_S3_SECRET_ACCESS_KEY
s3_bucket: your.blog.bucket.com
```

Run 

```
jkl --s3
```

### Documentation

See the official [Jekyll wiki](https://github.com/mojombo/jekyll/wiki)
... just remember that you are using Go templates instead of Liquid templates.

