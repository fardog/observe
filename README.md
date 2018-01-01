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

## Running

Create a BigQuery table with the following schema; observe cannot yet create
its own:

| Field          | Type      | Attributes |
|----------------|-----------|------------|
| URL            | STRING    | NULLABLE   |
| RemoteAddr     | STRING    | REQUIRED   |
| Observed       | TIMESTAMP | REQUIRED   |
| Header         | RECORD    | REPEATED   |
| _Header_.Key   | STRING    | REQUIRED   |
| _Header_.Value | STRING    | REPEATED   |

```bash
go get -u github.com/fardog/observe/cmd/observe
observer -gcloud-project-id my-project -bigquery-table observe.observations
```

You may run `observer -help` to view all options.

## Usage

When running Observe exposes a single endpoint, which serve a 1x1 pixel
transparent GIF image: `/observe.gif`. Whenver that image is loaded, an entry
will be written to BigQuery with the following:

* The URL visited (taken from the Referrer headers)
* The IP address of the viewer
* The time observed
* A list of all header values the browser sent

You may also manually specify the URL visited, if you don't want to rely on the
referrer headers: `/observe.gif?referrer=https://cool.site/awesome-page/`

Observe respects the [Do Not Track][DNT] header, and will serve the gif to the
requester but not log or store their information. There is no option to disable
this behavior.

## TODO

* Attempt to create schema if not present
* Figure out if there's a decent way to distribute the schema with this app
* Document running on AppEngine

## License

[MIT](./LICENSE)

[DNT]: https://en.wikipedia.org/wiki/Do_Not_Track
