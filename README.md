# cke - CUE kubernetes render helper

Render helper program to dynamically output cue kubernetes manifests to
either stdout or a \_rendered directory.

To see available options run:
`just`

```
Available recipes:
    run     # Output test cue to stdout
    render  # Output test cue to _rendered/
    install # build and install go binary
```

## Install

`go install github.com/Kystverket/cke-cue-kubernetes-resource-exporter`

## Usage

Running `cke` will output any kubernetes manifest to stdout

```
: cke
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  namespace: mynamespace
---
apiVersion: ""
kind: Service
metadata:
  name: myservice
  namespace: mynamespace
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
```

Adding `-out files` parameter will output into a \_rendered directory

```
: cke -out files
Created: _rendered/mynamespace-deployment-myapp.yaml
Created: _rendered/mynamespace-service-myservice.yaml
Created: _rendered/deployment-myapp.yaml
```
