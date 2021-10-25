# Bootstrapping

The process for bootstrapping GLab is by far the most complex part of the
project. As noted in the overview, reproducability is an important trait to
maintain and therefore much thought has been put into creating a bootstrap
process which is idempotent. The bootstrap process is first described from the
local development standpoint where the stack is brought up by Vagrant on a local
development machine. A later section covers the additional steps required when
deploying the stack to a vCenter cluster.

## Development

Below is a high-level overview of the bootstrap process on a local development
machine:

1. The project repository is cloned to the local machine
2. An environment file is sourced to pull-down sensitive bootstrap parameters
3. Vagrant is invoked to bring up a 3-node version of the stack
4. A development container is run which provides a pre-configured environemnt
   for interacting with the local stack.

Each of these steps is further expounded upon below.

### Cloning

The GLab project utilizes a monorepo approach in which the entirety of the
project is contained within a single repository. While there are many advantages
and disadvantages to this approach - the boostrap process greatly benefits from
a monorepo format where everything is self-contained and accessible without
having to fetch additional repositories. In some cases auxiliary tools are
located in separate repositories to avoid cluttering the main repo and those
tools are cloned and utilized as needed.

The project repository contains all static configuration data needed in order
to bootstrap the stack. This data is stored in YAML configuration files and is
ultimately fed into the Consul KV store in order to make configuration data
widely accessible across the stack. Ansible also utilizes these configuration
files in order to generate the files necessary to bootstrap the stack. 

### Environment File

Sensitive data that should remain outside of change control is published in two 
places:

1. The AWS Parameter Store.
2. Vault KV store

The primary secret store for the stack is Hashicorp Vault, however, when
bootstrapping the stack the AWS Parameter Store is utilized as an outside
reference since Vault is not yet available. The AWS Parameter Store only
includes the information necessary to bootstrap the stack and all sensitive data
contained within is copied into Vault during the bootstrap process.

An environment file is included in the project repository for easily pulling
down the data into local environment variables. Subsequent processes rely on
these environment variables for performing the bootstrap process.

### Vagrant

Vagrant provides an abstraction layer for bringing up the nodes on a local 
development machine. It handles the process of starting the virtual machines,
handing off the Ignition file to Flatcar Linux, and then starting the bootstrap
process.