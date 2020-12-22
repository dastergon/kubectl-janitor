# kubectl janitor

[![Build Status](https://github.com/dastergon/kubectl-janitor/workflows/ci/badge.svg
)](https://github.com/dastergon/kubectl-janitor/workflows/ci/)
[![Go Report Card](https://goreportcard.com/badge/dastergon/kubectl-janitor)](https://goreportcard.com/report/dastergon/kubectl-janitor)
[![LICENSE](https://img.shields.io/github/license/dastergon/kubectl-janitor.svg)](https://github.com/dastergon/kubectl-janitor/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/release-pre/dastergon/kubectl-janitor.svg)](https://github.com/dastergon/kubectl-janitor/releases)

`kubectl janitor` is a kubectl plugin that assists in finding objects in a problematic state in your Kubernetes cluster.

## Introduction
Troubleshooting Kubernetes clusters sometimes requires a [combination](https://learnk8s.io/troubleshooting-deployments) of `kubectl` commands and other command-line tools such as [jq](https://github.com/stedolan/jq) to do correlations around the issues that the various objects might have. Moreover, sometimes the supported options of the `--field-selector` flag might be [limited](https://github.com/kubernetes/kubernetes/issues/49387).

During troubleshooting scenarios, people need to identify the issues quickly without worrying about remembering all the different command combinations. The primary goal of this plugin is to collect some commonly executed kubectl command patterns to identify objects in a problematic state in the cluster and reduce the cognitive load for people troubleshooting.

## Installing

### Krew

You can install `kubectl janitor` using the [Krew](https://github.com/kubernetes-sigs/krew), the package manager for kubectl plugins.

Once you have Krew [installed](https://krew.sigs.k8s.io/docs/user-guide/setup/install/) run the following command:

    kubectl krew install janitor

### Releases

Check the [release](https://github.com/dastergon/kubectl-janitor/releases) page for the full list of pre-built assets.

#### Install

1. Download one of the releases that are compatible with your os/arch.
2. Unzip to get `kubectl-janitor`
3. Add it to your `PATH` or move it to a path already in in `PATH` (i.e., `/usr/local/bin`)

### Source

    go get -u github.com/dastergon/kubectl-janitor/cmd/kubectl-janitor

This command will download and compile `kubectl-janitor`.

## Usage

To get the full list of commands with examples:

    kubectl janitor

### Features

#### List Pods that are in a pending state (waiting to be scheduled)

    kubectl janitor pods unscheduled

#### List Pods in an unhealthy state

    kubectl janitor pods unhealthy

#### List Pods that are currently running but not ready for some reason

    kubectl janitor pods unready

#### List the current statuses of the Pods and their respective count

    kubectl janitor pods status

#### List Jobs that have failed to run and have restartPolicy: Never

    kubectl janitor jobs failed

#### List PesistentVolumes that are available for claim

    kubectl janitor pvs unclaimed

#### List PersistentVolumeClaims in a pending state (unbound)

    kubectl janitor pvcs pending

You can use the `-A` or `--all-namespaces` flag to search for objects in all namespaces.

You can use the `--no-headers` flag to avoid showing the column names.

## Cleanup
If you have installed the plugin via the `krew` command. You can remove the plugin by using the same tool:

    kubectl krew uninstall kubectl-janitor

Or, you can "uninstall" this plugin from kubectl by simply removing it from your PATH:

    rm /usr/local/bin/kubectl-janitor

## Author

Pavlos Ratis [@dastergon](https://twitter.com/dastergon).

## License

[Apache 2.0.](./LICENSE)
