# Upgrading guide

## v1 to v2

The v2 release adds support for Go modules, removes web framework
integrations from the main repository, and supports library configuration
through environment variables. The following breaking changes occurred as a part
of this release:

### Removed `Configuration.Endpoint`

The `Endpoint` configuration option was deprecated as a part of the v1.4.0
release in November 2018. It was replaced with `Endpoints`, which includes
options for configuring both event and session delivery.

```diff+go
- config.Endpoint = "https://notify.myserver.example.com"
+ config.Endpoints = {
+ 	Notify: "https://notify.myserver.example.com",
+ 	Sessions: "https://sessions.myserver.example.com"
+ }
```
