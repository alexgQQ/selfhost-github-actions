This is a repo for testing out some features with Github Actions, particularly with a focus on self hosting on k8s. This contains a dummy golang application that can be hosted with Google App Engine and is built/tested/deployed with Github Actions. In particular I wanted to make a deployment pipeline where every push action is tested, semver tags build a release and publishing deploys the specific release version.

## Requirements:

* [golang 1.16+](https://go.dev/doc/install)
* [reckoner](https://github.com/FairwindsOps/reckoner/releases) and by extension [helm](https://helm.sh/docs/intro/install/) and [kubectl](https://kubernetes.io/docs/tasks/tools/)
* [terraform 0.14+](https://developer.hashicorp.com/terraform/downloads)
* [gcloud cli](https://cloud.google.com/sdk/docs/install)
* [docker](https://docs.docker.com/get-docker/)

A GCP project with billing enabled, cli authentication setup and these services enabled:
* App Engine Admin API
* Cloud Build API

A Github repo. Feel free to fork this to make your own.

## Setup

Test/run the go app. It just takes GET to the root path and responds in json with the status, git sha and version.
```bash
cd app
go test -v .
go run main.go  # Endpoint listens on port 8080
```

### Google App Engine

Define our project and create the initial resources from the cli

```bash
gcloud config set project myproject-12345
gcloud app create
```

If not using the self hosted runners then an extra step is needed for GCP authentication. Create a key for the GAE default service account and make it available to the action runners on the repo.

```bash
gcloud iam service-accounts keys create key.json --iam-account myproject-12345@appspot.gserviceaccount.com
cat key.json | xclip -sel clip  # Copy file contents to clipboard
```

[Create a repo secret](https://docs.github.com/en/actions/security-guides/encrypted-secrets#creating-encrypted-secrets-for-a-repository) that contains the json file data from the previous step named `GCP_CREDENTIALS`. Make sure to enable the `auth` step in the deployment actions workflow too.

### GKE and Actions runner

Create the cluster for your specific project

```bash
cd k8s/terraform

# First time setup
echo "
gke_project          = "myproject-12345"
gke_zone             = "us-central1-a"
gke_region           = "us-central1"
" > terraform.tfvars
terraform init

terraform apply
```

That will create a k8s cluster named `primary-cluster`. When done go ahead and get the cluster credentials.
```bash
gcloud container clusters get-credentials --region us-central1-a primary-cluster
```

Jump over to the reckoner dir to build our custom image. May need to run `gcloud auth configure-docker` if you encounter auth errors.

```bash
cd ../reckoner
docker build -t us.gcr.io/myproject-12345/github-actions:latest .
docker push us.gcr.io/myproject-12345/github-actions:latest
```

Before applying the resources, generate a [Github PAT](https://github.com/settings/tokens/) with read and write access to the target repo. Copy and set the token as the `github_token` key in the reckoner course.yml file. Set the `repository` value on the runnerdeployment.yaml to the target Github repo.
```bash
reckoner plot course.yml --run-all
```

This will setup the runner controller and registers the target repo to use it for workflows.
