# observe

Simple traffic analytics collection for static websites. Observe is designed for
a situation where:

* You have a static web site
* You don't control the site's server (say, GitHub Pages)
* You don't want Google Analytics
* But, you _do_ want some form of metrics.

Observe is designed to be runnable in the free tier of Google Cloud; right now
it only supports BigQuery as a storage backend, and is intended to be run under
AppEngine. It also has a standalone server which can be run anywhere.

**Warning:** This is alpha software; it's not actually run anywhere yet.

## Running

Create a BigQuery table with [provided schema][schema]; observe cannot yet
create its own:

```bash
# get the code
go get -u github.com/fardog/observe/cmd/observe
# create the schema
bq mk --table --description 'Observed visits' \
  --time_partitioning_field=Observed \
  --time_partitioning_expiration 7776000 \
  <gcloud_project_name>:observe.observations \
  $GOPATH/src/github.com/fardog/observe/bq-schema.json
# run the observer
observe -gcloud-project-id <gcloud_project_name> -bigquery-table observe.observations
```

You may run `observe -help` to view all options.

## Usage

When running Observe exposes a single endpoint, which serve a 1x1 pixel
transparent GIF image: `/observe.gif`. Whenver that image is loaded, an entry
will be written to BigQuery with the following:

* The URL visited (taken from the Referrer headers)
* The anonymized IP address of the viewer (IPv4 to 20 bits, IPv6 to 32 bits)
* The time observed
* A list of all header values the browser sent

You may also manually specify the URL visited, if you don't want to rely on the
referrer headers: `/observe.gif?referrer=https://cool.site/awesome-page/`

Observe respects the [Do Not Track][DNT] header, and will serve the gif to the
requester but not log or store their information. There is no option to disable
this behavior.

## TODO

* Attempt to create schema if not present
* Document running on AppEngine
* Tests

## License

[MIT](./LICENSE)

[DNT]: https://en.wikipedia.org/wiki/Do_Not_Track
[schema]: ./bq-schema.json
