# GPGCLI

## Download

Executables are inside the `dist` folder

## Generating GPG keys

1. List all available keys

    ```
    gpg --list-keys
    ```

2. Create a new key
    The adviced configuration is:

    | Question    | Answer                                |
    | ---         | ---                                   |
    | kind of key | (1) RSA and RSA (default)             |
    | keysize     | 4096                                  |
    | *validity   | 0 = key does not expire (for staging) |

    *validity should be specified for production.

    ```
    gpg --full-generate-key
    ```

3. Export the keys

    ```
    gpg --armor --output <uid>.gpg.pub --export <uid>
    gpg --armor --output <uid>.gpg.pvt --export-secret-keys <uid>
    ```

## Usage

To encrypt a file,

```
gpgcli encrypt <file> --public <public-file-gpg> --output <outfile>
```

Omitting the `--output` will dump the result in base64

For when decrypting,

```
gpgcli encrypt <file> --public <public-file-gpg> --output <outfile>
```

## Troubleshooting
| #   | Error                                                   |                                                                                |
| --- | ---                                                     | ---                                                                            |
|  1  | `missing command`                                       | use either `help`, `encrypt`, or `decrypt`                                       |
|  2  | ``command `<command>` not supported``                   | similar to 1                                                                   |
|  3  | `missing file`                                          | missing file to be encrypted or decrypted                                      |
|  4  | ``missing `--public` argument while using encrypt``     | `encrypt` command should contain a `--public` which points the public GPG file |
|  5  | ``missing `--secret` argument while using decrypt``     | `decrypt` command should contain a `--secret` which points the secret GPG file |
|  6  | ``missing `--passphrase` argument while using decrypt`` | `decrypt` command should contain a `--passphrase` with the secret passphrase   |

## .onsave

`.onsave` is a plugin used by Visual Studio Code which executes when files are changed. Inside the `.onsave` file is the linux command on how to compile the application in multiple OS' and stored in the `dist` folder.
