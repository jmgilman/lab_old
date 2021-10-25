# Bootstrap Secrets

The primary secret provider for the GLab stack is Hashicorp Vault. However, 
relying on it during the bootstrap process creates a chicken/egg scenario
because initially the service is not available to be used. To workaround this,
the AWS Parameter Store is utilized to provide an early method for obtaining
sensitive data required for the bootstrap process. 

The process for creating, setting, and deleting bootstrap secrets is done via
the CLI tool. The purpose behind this is to prevent vendor lock-in to one 
specific secret provider for the bootstrap process. The CLI tool by default uses
the AWS Parameter Store but can be further expanded to use other backends as 
needed.

Most bootstrap secrets persist even after the bootstrap process and can vary
from third-party API credentials to SSH provisioning certificates. To increase
security, credentials that are not used outside the stack are randomly generated
during the bootstrap process and persisted when Vault is configured.