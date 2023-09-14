## Installation

`aftctl` is available to install from official releases as described below. I recommend that you install `aftctl` from only the official GitHub releases.

### For Mac
To download the latest release, run:

```sh
# for ARM systems, set ARCH to: `arm64`
ARCH=x86_64

curl -sLO "https://github.com/edgarsilva948/aftctl/releases/latest/download/aftctl_Darwin_$ARCH.tar.gz"

# (Optional) Verify checksum
curl -sL "https://github.com/edgarsilva948/aftctl/releases/latest/download/checksums.txt" | grep $ARCH | sha256sum --check

tar -xzf aftctl_Darwin_$ARCH.tar.gz -C /tmp && rm aftctl_Darwin_$ARCH.tar.gz

sudo mv /tmp/aftctl /usr/local/bin
```

### For Linux
To download the latest release, run:

```sh
# for ARM systems, set ARCH to: `arm64` or `i386`
ARCH=x86_64

curl -sLO "https://github.com/edgarsilva948/aftctl/releases/latest/download/aftctl_Linux_$ARCH.tar.gz"

# (Optional) Verify checksum
curl -sL "https://github.com/edgarsilva948/aftctl/releases/latest/download/checksums.txt" | grep $ARCH | sha256sum --check

tar -xzf aftctl_Linux_$ARCH.tar.gz -C /tmp && rm aftctl_Linux_$ARCH.tar.gz

sudo mv /tmp/aftctl /usr/local/bin
```

### For Windows

#### Direct download (latest release): [x86_64](https://github.com/edgarsilva948/aftctl/releases/download/v0.1.0/aftctl_Windows_x86_64.zip) - [ARM64](https://github.com/edgarsilva948/aftctl/releases/download/v0.1.0/aftctl_Windows_arm64.zip) - [i386](https://github.com/edgarsilva948/aftctl/releases/download/v0.1.0/aftctl_Windows_i386.zip)

Make sure to unzip the archive to a folder in the `PATH` variable. 

Optionally, verify the checksum: 

1. Download the checksum file: [latest](https://github.com/edgarsilva948/aftctl/releases/latest/download/checksums.txt)
2. Use Command Prompt to manually compare `CertUtil`'s output to the checksum file downloaded. 
```cmd
# Replace x86_64 with ARM64 or i386
CertUtil -hashfile aftctl_Windows_x86_64.zip SHA256
```
3. Using PowerShell to automate the verification using the `-eq` operator to get a `True` or `False` result:
```pwsh
# Replace x86_64 with ARM64 or i386
(Get-FileHash -Algorithm SHA256 .\aftctl_Windows_x86_64.zip).Hash -eq ((Get-Content .\checksums.txt) -match 'aftctl_Windows_x86_64.zip' -split ' ')[0]
```