# Gitfox Library Webhooks

Library webhooks allows for easy receiving and parsing of Gitfox Webhook Events

## Installation

Use go get.

```go
go get -u github.com/easysoft/gitfox-webhooks
```

Then import the package into your own code.

```go
import "github.com/easysoft/gitfox-webhooks"
```

## Usage

```go
package main

import (
 "net/http"

 "github.com/easysoft/gitfox-webhooks/gitfox"
)

func main() {
 hook, _ := gitfox.New(gitfox.Options.Secret("MyGitfoxSecret...?"))

 http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
  payload, err := hook.Parse(r, gitfox.BranchUpdatedEvent)
  if err != nil {
   if err == gitfox.ErrEventNotFound {
    // ok event wasn't one of the ones asked to be parsed
   }
  }
  switch payload.(type) {
  case gitfox.BranchUpdatedPayload:
   push := payload.(gitfox.BranchUpdatedPayload)
   // Do whatever you want from here...
   fmt.Printf("%+v", push)
 })
 http.ListenAndServe(":3000", nil)
}
```
