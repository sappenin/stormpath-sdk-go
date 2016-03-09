Go App Engine SDK for the [Stormpath](http://stormpath.com/) API

This is a fork of the excellent [stormpath-sdk-go library](https://github.com/sappenin/stormpath-sdk-go).  This version removes various dependencies and code-paths that don't play nicely 
with Google App Engine.  For example:

- **URL Fetch**: Google App Engine requires outbound HTTP requests to be made using the URL Fetch service.  In order to work, UrlFetch requires a request Context
 object that varies on every request.  To make this work, the library has been adjusted to remove its singleton HTTP Client and to allow callers to pass-in a 
 Context object when interacting with the SDK.
- **User-Agent Header**: In the original version of this library, the *User-Agent* header is signed as part of the SAuthc1 algorithm that StormPath employs to sign
 its requests.  However, the URL Fetch service appends various environment-related information to the User-Agent, but only once the request is inside of Google's infrastructure.
 For now, this library removes the User-Agent header from the signature computed for each request.  In the future, additional logic may be added to re-introduce this
 header into the digest calculation.
- **GET/DELETE Requests**: The [original version](https://github.com/jarias/stormpath-sdk-go) of this library errantly introduces a payload of double-quotes on all 
 HTTP GET and DELETE requests, which is also signed as part of the StormPath [SAuthc1](https://github.com/stormpath/stormpath-sdk-spec/blob/master/specifications/algorithms/sauthc1.md) 
 algorithm.  Inside of App Engine, these payloads are stripped by the URL Fetch service, which makes all GET and DELETE requests fail the SAuthc1 verification 
 inside of StormPath.  This client fixes this by using an empty payload for GET and DELETE requests.  See [here](https://github.com/jarias/stormpath-sdk-go/issues/23) 
 for the fix request in the original client.

Develop:

[![Build Status](https://travis-ci.org/sappenin/stormpath-sdk-go.svg?branch=develop)](https://travis-ci.org/sappenin/stormpath-sdk-go) [![codecov.io](http://codecov.io/github/sappenin/stormpath-sdk-go/coverage.svg?branch=develop)](http://codecov.io/github/sappenin/stormpath-sdk-go?branch=develop)

Master:

[![Build Status](https://travis-ci.org/sappenin/stormpath-sdk-go.svg?branch=master)](https://travis-ci.org/sappenin/stormpath-sdk-go) [![codecov.io](http://codecov.io/github/sappenin/stormpath-sdk-go/coverage.svg?branch=master)](http://codecov.io/github/sappenin/stormpath-sdk-go?branch=master)

# Usage

```go get github.com/sappenin/stormpath-sdk-go```

```go
import "github.com/sappenin/stormpath-sdk-go"
import "fmt"

//This would look for env variables first STORMPATH_API_KEY_ID and STORMPATH_API_KEY_SECRET if empty
//then it would look for os.Getenv("HOME") + "/.config/stormpath/apiKey.properties" for the credentials
credentials, _ := stormpath.NewDefaultCredentials()

// Caching is optional.  This line creates cache with a default expiration time of 5 minutes, and which 
// purges expired items every 30 seconds
var c *cache.Cache = cache.New(5 * time.Minute, 30 * time.Second)

// Init with Cache.  Pass nil instead for no caching.
stormpath.Init(credentials, stormpath.CacheableCache{Cache: c})

// r is of type *http.Request.  All client usage must be in the context of a request so that the 
// appengine.Context can be properly populated.
ctx := appengine.NewContext(r)

//Get the current tenant
tenant, _ := stormpath.CurrentTenant(ctx)

//Get the tenant applications
apps, _ := tenant.GetApplications(ctx, stormpath.MakeApplicationCriteria().NameEq("test app"))

//Get the first application
app := apps.Items[0]

//Authenticate a user against the app
account, _ := app.AuthenticateAccount(ctx, "username", "password")

fmt.Println(account)
```

Features:

* Cache via [go-cache](https://github.com/patrickmn/go-cache) implementation
* Almost 100% of the Stormpath API implemented
* Load credentials via properties file or env variables
* Requests are authenticated via Stormpath SAuthc1 algorithm

# Debugging

If you need to trace all requests done to stormpath you can enable debugging in the logs
by setting the environment variable STORMPATH_LOG_LEVEL=DEBUG the default level is ERROR.

# Contributing

Pull request are more than welcome, all pull requests should be from and directed to the ```develop``` branch **NOT** ```master```.

Please make sure you add tests ;)

Development requirements:

- Go 1.4+
- [Ginkgo](https://onsi.github.io/ginkgo/) ```go get github.com/onsi/ginkgo/ginkgo```
- [Gomega](http://onsi.github.io/gomega/) ```go get github.com/onsi/gomega```
- An [Stormpath](https://stormpath.com) account (for integration testing)

Running the test suite

Env variables:

```
export STORMPATH_API_KEY_ID=XXXX
export STORMPATH_API_KEY_SECRET=XXXX
```

```
ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover
```

I'm aiming at 85% test coverage not yet met but thats the goal.

# License

Copyright 2014, 2015 Julio Arias
Copyright 2016 Sappenin Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
