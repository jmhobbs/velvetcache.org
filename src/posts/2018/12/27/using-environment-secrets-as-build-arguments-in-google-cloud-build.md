---
category:
- Geek
creator: admin
date: 2018-12-27
tags:
- Docker
- GCP
- Google Cloud Platform
- guide
- How To
title: Using environment secrets as build arguments in Google Cloud Build
type: post
permalink: /2018/12/27/using-environment-secrets-as-build-arguments-in-google-cloud-build/
wp_id: "2836"
summary: >
  Cloud Build is a nice tool for continuous building of your Docker images. Using
  an environment secret in the build is a gap in the docs, here is how I did it.
---

[Google Cloud Build](https://cloud.google.com/cloud-build/) is a pretty nice tool for building your docker images continually, and cloud-build-local is pretty great for working on your images in dev.  All around, a nice piece of kit to have in a Kubernetes shop.

The docs are pretty good, but one thing that I've recently dealt with did not show up in my searching; how to use an environment secret as a build argument to Docker.  So here's how I found to do it.

First, we will follow the [encrypted secrets guide](https://cloud.google.com/cloud-build/docs/securing-builds/use-encrypted-secrets-credentials#encrypting_an_environment_variable_using_the_cryptokey) to get a secret wrapped up by KMS.

```shell
$ # Create a keyring and encryption key
$ gcloud kms keyrings create tinkering --location=global
$ gcloud  kms keys create cloud-build-demo \
  --keyring=tinkering \
  --purpose=encryption \
  --location=global
$ Now encrypt our secret string
$ echo -n "This is the super secret secret." | gcloud kms encrypt \
  --plaintext-file=- \
  --ciphertext-file=- \
  --location=global \
  --keyring=tinkering \
  --key=cloud-build-demo | base64
CiQATajs0GI7M6ZFM68Qu+GbJTfJ/d3tqqLcHz69RY1AaHkzV20SSQDt7E4V65imqbOnq8DvieiaglxjEztxWQCwrr2Mtu+xwT6tko6FHB+NNauyos6X1nnh5x217Cwx5QbX3h0YtjOJ15I4dnHDM+I=
```

Next, we will create a super simple Dockerfile to show how it is used.

```dockerfile
FROM busybox
ARG THE_SECRET
RUN echo "::${THE_SECRET}::"
```

Last, we set up the cloudbuild.yaml.  In the documentation demo files they use a shell entrypoint to access the environment variable.

```yaml
# Note: You need a shell to resolve environment variables with $$
- name: 'gcr.io/cloud-builders/docker'
  entrypoint: 'bash'
  args: ['-c', 'docker login --username=[MY-USER] --password=$$PASSWORD']
  secretEnv: ['PASSWORD']
```

However, it would be nicer to not have to stringify our whole Docker build command.

Luckily, using `--build-arg` [without a value falls through to the environment variable of the same name.](https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg)

```shell
$ export HTTP_PROXY=http://10.20.30.2:1234
$ docker build --build-arg HTTP_PROXY .
```

So, we can just use it directly:

```yaml
steps:
- id: docker
  name: 'gcr.io/cloud-builders/docker'
  args: ['build', '--build-arg', 'THE_SECRET', '.']
  secretEnv: ['THE_SECRET']
secrets:
- kmsKeyName: projects/hobbs-tinkering/locations/global/keyRings/tinkering/cryptoKeys/cloud-build-demo
  secretEnv:
    THE_SECRET: CiQATajs0GI7M6ZFM68Qu+GbJTfJ/d3tqqLcHz69RY1AaHkzV20SSQDt7E4V65imqbOnq8DvieiaglxjEztxWQCwrr2Mtu+xwT6tko6FHB+NNauyos6X1nnh5x217Cwx5QbX3h0YtjOJ15I4dnHDM+I=
```

Testing locally, it happily runs:

<!-- todo: mark line 20 -->
```shell
$ cloud-build-local --dryrun=false .
Using default tag: latest
latest: Pulling from cloud-builders/metadata
Digest: sha256:bcdb85e67ab9719c6441cb80fe9e8badc6d5ab0ab8bc73ee67adc0112233d20c
Status: Image is up to date for gcr.io/cloud-builders/metadata:latest
2018/12/23 13:18:29 Started spoofed metadata server
2018/12/23 13:18:29 Build id = localbuild_9cef1240-3a68-4ec3-a273-f49cd018316d
2018/12/23 13:18:29 status changed to "BUILD"
BUILD
Starting Step #0 - "docker"
Step #0 - "docker": Already have image (with digest): gcr.io/cloud-builders/docker
Step #0 - "docker": Sending build context to Docker daemon   5.12kB
Step #0 - "docker": Step 1/3 : FROM busybox
Step #0 - "docker":  ---> 59788edf1f3e
Step #0 - "docker": Step 2/3 : ARG THE_SECRET
Step #0 - "docker":  ---> Using cache
Step #0 - "docker":  ---> f289a756b157
Step #0 - "docker": Step 3/3 : RUN echo "::${THE_SECRET}::"
Step #0 - "docker":  ---> Running in 0e90f8f4f349
Step #0 - "docker": ::This is the super secret secret.::
Step #0 - "docker": Removing intermediate container 0e90f8f4f349
Step #0 - "docker":  ---> 75d19dee1d47
Step #0 - "docker": Successfully built 75d19dee1d47
Finished Step #0 - "docker"
2018/12/23 13:18:35 status changed to "DONE"
DONE
```

It is worth noting that using build args for secrets is not recommended.  Anyone with the image can see what the argument passed in was.

<!-- todo: mark line 3 -->
```shell
$ docker history 75d19dee1d47
IMAGE               CREATED             CREATED BY                                      SIZE                COMMENT
75d19dee1d47        4 days ago          |1 THE_SECRET=This is the super secret secre…¦   0B
f289a756b157        4 days ago          /bin/sh -c #(nop)  ARG THE_SECRET               0B
59788edf1f3e        2 months ago        /bin/sh -c #(nop)  CMD ["sh"]                   0B
<missing>           2 months ago        /bin/sh -c #(nop) ADD file:63eebd629a5f7558c…¦   1.15MB
```

Docker 18.09, added [build secrets](https://medium.com/@tonistiigi/build-secrets-and-ssh-forwarding-in-docker-18-09-ae8161d066) for a better solution, but GCB is still running Docker 17.12, so we will have to wait for that update.

A gist of the code is available at: [https://gist.github.com/jmhobbs/a572b47048eb42803bcb2102ac57a8df](https://gist.github.com/jmhobbs/a572b47048eb42803bcb2102ac57a8df)

