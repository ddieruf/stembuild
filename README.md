# Stembuild

The stembuild binary is used to build BOSH stemcells for **Windows 2012R2**,**Windows Server, version v1709**, **Windows Server, version 1803**, **Windows Server 2019** on **vSphere**. 

**Instructions**: See [here](https://bosh.io/docs/windows-stemcell-create/) for instructions to build Windows stemcells for vSphere.

## Installation
Download the latest stembuild from the [Releases page](https://github.com/cloudfoundry-incubator/stembuild/releases) that corresponds to the operating system of your local host and the stemcell version that you want to build

## Dependencies
[LGPO](https://www.microsoft.com/en-us/download/details.aspx?id=55319) must be downloaded in the same folder as your `stembuild`

## Current Commands
```
stembuild version <STEMCELL-VERSION>, Windows Stemcell Building Tool

Usage: stembuild <global options> <command> <command flags>

Commands:
  help		Describe commands and their syntax
  package	Create a BOSH Stemcell from a VMDK file or a provisioned vCenter VM
  construct	Provisions and syspreps an existing VM on vCenter, ready to be packaged into a stemcell

Global Options:
  -color	Colorize debug output
  -debug	Print lots of debugging information
  -v		Stembuild version (shorthand)
  -version	Show Stembuild version

```
## `stembuild construct`

This command provisions and syspreps an existing VM on vCenter. It prepares a VM to be used by `stembuild package`.

```
stembuild construct -vm-ip <IP of VM> -vm-username <vm username> -vm-password <vm password>  -vcenter-url <vCenter URL> -vcenter-username <vCenter username> -vcenter-password <vCenter password> -vm-inventory-path <vCenter VM inventory path>
```

### Requirements
- LGPO.zip in current working directory
- Running Windows VM with:
	- Up to date Operating System
	- Reachable by IP over port 5985
	- Username and password with Administrator privileges
	- vCenter URL, username and password
	- vCenter Inventory Path
- The `vm-ip`, `vm-username`, `vm-password`, `vcenter-url`, `vcenter-username`, `vcenter-password`, `vm-inventory-path` must be specified

```
Example:
	stembuild construct -vm-ip '10.0.0.5' -vm-username Admin -vm-password 'password' -vcenter-url vcenter.example.com -vcenter-username root -vcenter-password 'password' -vm-inventory-path '/datacenter/vm/folder/vm-name'

Flags:
  -vcenter-ca-certs string
    	filepath for custom ca certs
  -vcenter-password string
    	vCenter password
  -vcenter-url string
    	vCenter url
  -vcenter-username string
    	vCenter username
  -vm-inventory-path string
    	vCenter VM inventory path. (e.g: <datacenter>/vm/<vm-folder>/<vm-name>)
  -vm-ip string
    	IP of target machine
  -vm-password string
    	Password of target machine. Needs to be wrapped in single quotations.
  -vm-username string
    	Username of target machine
	
```

### Troubleshooting
After running `stembuild construct`, you may find yourself with a connection issue to the VM
- Confirm port 5985 is reachable via something like `nmap [vm-ip] -Pn`


## `stembuild package`

This command creates a BOSH Stemcell from a provisioned vCenter VM 

```
  stembuild package -vcenter-url <vCenter URL> -vcenter-username <vCenter username> -vcenter-password <vCenter password> -vm-inventory-path <vCenter VM inventory path>
```

*Requirements*:
- VM provisioned using the stembuild construct command
- Access to vCenter environment
- The `vcenter-url`, `vcenter-username`, `vcenter-password`, and `vm-inventory-path` flags must be specified.
- **NOTE**: The 'vm' keyword must be included between the datacenter name and folder name for the vm-inventory-path (e.g: <datacenter>/vm/<vm-folder>/<vm-name>)
 
```
Example:
 stembuild package -vcenter-url vcenter.example.com -vcenter-username root -vcenter-password 'password' -vm-inventory-path '/my-datacenter/vm/my-folder/my-vm'

Flags:
  -o string
    	Output directory (shorthand)
  -outputDir string
    	Output directory, default is the current working directory.
  -vcenter-ca-certs string
    	filepath for custom ca certs
  -vcenter-password string
    	vCenter password
  -vcenter-url string
    	vCenter url
  -vcenter-username string
    	vCenter username
  -vm-inventory-path string
    	vCenter VM inventory path. (e.g: <datacenter>/vm/<vm-folder>/<vm-name>)
```

## [DEPRECATED] Package a Windows Stemcell from a VMDK using `stembuild package`

This command converts a VMDK into a bosh-deployable Windows Stemcell 

The VMware 'ovftool' binary must be on your path or Fusion/Workstation must be installed (both include the 'ovftool').

```
stembuild package -vmdk <path-to-vmdk>
```

*Requirements*
- The VMware 'ovftool' binary must be on your path or Fusion/Workstation must be installed (both include the 'ovftool').
- The `vmdk` flag must be specified.  If the `output` flag is not specified the stemcell will be created in the current working directory.

```
Example:
	stembuild package -vmdk my-1803-vmdk.vmdk
	
	Will create an Windows 1803 stemcell using [vmdk] 'my-1803-vmdk.vmdk'
	The final stemcell will be found in the current working directory.

Flags:
  -o string
    	Output directory (shorthand)
  -outputDir string
    	Output directory, default is the current working directory.
  -vmdk string
    	VMDK file to create stemcell from

```

Process can take between 10 and 20 minutes. See Progress with `-debug` flag.

## Testing

Older tests were written using the default testing framework.  However, more recent code
has been test-driven with Ginkgo.  We recommend that any new code be test-driven using Ginkgo.
Below are steps to run the tests:

Make puts some files in `out` dir. To clean state of this dir:
```
make clean
```
To run only unit tests:
```
make units
```
To run integration tests:
```
make integration
```


## Vendoring

Vendoring for this project is done using `dep`. 
To sync all the dependecies run
```
dep ensure
```

To add a new dependency run 
```
dep ensure -add <git package url>
```
like 
```
dep ensure -add github.com/google/subcommands
```

To check if dependencies are in sync or not run
```
dep sync
```
The output should be nothing if there no out-of-sync dependencies.
