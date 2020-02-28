# kube-shell
Simple docker image to deploy as a Pod or Daemonset in Kubernetes for interactive sessions

To build and push docker image:

```bash
docker build -t jonmoter/kube-shell:latest .
docker push jonmoter/kube-shell:latest
```

## What's in the image

To get an exhaustive list of everything in the image, read the `Dockerfile`. Here are
some highlights of what's present:

* `curl`, of course
* [httpie](https://httpie.org/), a prettier version of curl
* [grpcurl](https://github.com/fullstorydev/grpcurl), like curl for gRPC
* [httpstat](https://github.com/davecheney/httpstat), which displays HTTP connection timing info
* [hey](https://github.com/rakyll/hey), an HTTP load generator
* TCP tools like `tcpdump` and `tcptraceroute`
* DNS tools like `dig` and `host` and `dnstop`
* `jq` for JSON wrangling

`kubectl` is not installed by default, because it's a pretty large binary.
But there's an `install_kubectl` script that can be used to download it if you need it.

### Truth Service

There's a minimal golang server that runs and listens on port 4242. You can use that
if you want to run this docker imagine in a Kubernetes cluster and you want something
you can send HTTP requests to. The binary is `/usr/local/bin/truthserver`

## Deploying to Kubernetes

There's a DaemonSet definition in the `kubernetes` subdirectory. You can use that to
run an instance of `kube-shell` on each node in your cluster. That can be handy if you
want to try sending a request from a particular node or AZ.
