# Terraform Provider for DTEKT.IO

## Requirements

To be able to test the provider two environment variables needs to be set:

* `DTEKT_API_TOKEN` - API token for the dtekt.io user
* `DTEKT_API_URL=http://127.0.0.1:8000` - optional API url. If empty provider will use the upstream Terraform provider

## Building provider

Run the following command to build the provider

```shell
go build -o terraform-provider-dtekt
```

## Testing provider

Under `./test` there is a sample test workspace that can be used to test provider after changes.

First, build and install the provider.

```shell
make install
```

Then run the following command to initialize the workspace and apply the sample configuration.

```shell
cd test && terraform init 
terraform plan
```

## Production deployment
Not implemented at the moment.
