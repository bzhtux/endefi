[![CI](https://github.com/bzhtux/endefi/actions/workflows/ci.yml/badge.svg)](https://github.com/bzhtux/endefi/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/bzhtux/endefi/graph/badge.svg?token=o1qcG7cV9A)](https://codecov.io/github/bzhtux/endefi)

# EnDeFi

Encrypt Decrypt Files

## Goal

Encrypt and decrypt files use a secret key.
Secret key can come from:

- env var (ENDEFI_SECRET_KEY)
- local yaml file ($HOME/.endefi/secret.yaml)
- bitwarden ?

## How to build EnDeFi ?

As it is a simple Golang application just run the command as below:

```shell
go build -o endefi main.go
```

And then move `endefi` binary to your `$PATH`. For example for me:

```shell
mv endefi ~/bin/
```

## Usage

### With ENV Vars

At this point ENDEFI only supports env vars to get external secret provider and secret key. To use ENDEFI with env vars, you have to generate a secret key first:

```shell
endefi keygen
2024/06/10 10:39:02 Starting EnDeFi
2024/06/10 10:39:02 Generate a new secret Key: ce8f829cc257b910ac11cf9fe57f3ea4d2811ff06b3b47e8f41149be56f234b7
```

And then setup env vars as below:

```shell
export ENDEFI_SECRET_PROVIDER=env
export ENDEFI_SECRET_KEY=ce8f829cc257b910ac11cf9fe57f3ea4d2811ff06b3b47e8f41149be56f234b7
```

### With Secret file

Create a secret file `$HOME/.endefi/secret.yaml`:

To use local secret file, you have to set env var `ENDEFI_SECRET_PROVIDER` to `local`:

```shell
export ENDEFI_SECRET_PROVIDER=local
```

```shell
mkdir ~/.endefi
touch ~/.endefi/secret.yaml
chmod 600 ~/.endefi/secret.yaml
```

Add your secret key to `$HOME/.endefi/secret.yaml`:

```yaml
key: ce8f829cc257b910ac11cf9fe57f3ea4d2811ff06b3b47e8f41149be56f234b7
provider: local
```

## Examples

I assume you already run the requirement above. First let's use a text file `/tmp/lorem.txt`:

```shell
cat /tmp/lorem.txt
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
```

### File

Run endefi to encrypt this file:

```shell
endefi encrypt -f /tmp/lorem.txt
2024/06/10 10:44:30 Starting EnDeFi
2024/06/10 10:44:30 Encrypt a new file: /tmp/lorem.txt
```

Let's see how `lorem.txt` looks like now:

```shell
cat /tmp/lorem.txt
y��20Vk�T��#}Ɓa�-h2,�4W��O��o6̃S�Dcr��� +�B_.�0ck} �31GS�`�y �)T7�&;iL �a&�s>A��KΣ�~jbU\»v��%Szy��(ܷ�;᯺N
                      s+�2�h.PÜ�6@:$e`'�ԧx"f�IGN��)b]NXr��
Ip1F
    ��B0�f+%   
```

And now decrypt this encrypted file:

```shell
endefi decrypt -f /tmp/lorem.txt
2024/06/10 10:45:44 Starting EnDeFi
2024/06/10 10:45:44 Decrypt a new file: /tmp/lorem.txt
```

Verify `lorem.txt` is now a plain text file:

```shell
cat /tmp/lorem.txt
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
```

### Directory

Let's create the following tree:

```shell
/tmp/demo/
├── dir1
│   ├── file11
│   └── file12
├── dir2
│   ├── file21
│   └── file22
├── file1
└── file2
```

Each file content its own name:

```shell
cat /tmp/demo/file1
file1
cat /tmp/demo/dir2/file21
file21
```

Let's encrypt `/tmp/demo` without recursive mode:

```shell
endefi encrypt -d /tmp/demo
2024/06/13 09:50:11 Starting EnDeFi
2024/06/13 09:50:11 Encrypt a new directory: /tmp/demo
2024/06/13 09:50:11 Encrypt a new file: /tmp/demo/file1
2024/06/13 09:50:11 Encrypt a new file: /tmp/demo/file2
```

Let's check `/tmp/demo/file1` content:

```shell
cat /tmp/demo/file1
4@9)�Zt3#O@�DZDzHA�*% 
```

Decrypt `/tmp/demo`:

```shell
endefi decrypt -d /tmp/demo
2024/06/13 09:51:26 Starting EnDeFi
2024/06/13 09:51:26 Decrypt a new directory: /tmp/demo
2024/06/13 09:51:26 Decrypt a new file: /tmp/demo/file1
2024/06/13 09:51:26 Decrypt a new file: /tmp/demo/file2
```

Verify:

```shell
cat /tmp/demo/file1
file1
```

Let's encrypt `/tmp/demo` with recursive mode:

```shell
endefi encrypt -d /tmp/demo -r
2024/06/13 08:43:40 Starting EnDeFi
2024/06/13 08:43:40 Encrypt all files recursively within /tmp/demo
2024/06/13 08:43:40 Encrypt a new file: /tmp/demo/dir1/file11
2024/06/13 08:43:40 Encrypt a new file: /tmp/demo/dir1/file12
2024/06/13 08:43:40 Encrypt a new file: /tmp/demo/dir2/file21
2024/06/13 08:43:40 Encrypt a new file: /tmp/demo/dir2/file22
2024/06/13 08:43:40 Encrypt a new file: /tmp/demo/file1
2024/06/13 08:43:40 Encrypt a new file: /tmp/demo/file2
```

Verify:

```shell
cat  /tmp/demo/dir2/file22
2εwY0�׊ ���bD'%  
```

Decrypt `/tmp/demo`:

```shell
endefi decrypt -d /tmp/demo -r
2024/06/13 08:44:42 Starting EnDeFi
2024/06/13 08:44:42 Decrypt all files recursively within /tmp/demo
2024/06/13 08:44:42 Decrypt a new file: /tmp/demo/dir1/file11
2024/06/13 08:44:42 Decrypt a new file: /tmp/demo/dir1/file12
2024/06/13 08:44:42 Decrypt a new file: /tmp/demo/dir2/file21
2024/06/13 08:44:42 Decrypt a new file: /tmp/demo/dir2/file22
2024/06/13 08:44:42 Decrypt a new file: /tmp/demo/file1
2024/06/13 08:44:42 Decrypt a new file: /tmp/demo/file2
```

Verify:

```shell
cat  /tmp/demo/dir2/file22
file22
```

## Roadmap

- [x] MVP#1 : Build a cli to encrypt and decrypt files with a secret key from harcoded value
- [x] MVP#2 : Use environment variable to get secret key.
- [x] MVP#3 : Use file to get secret key
- [x] MVP#4 : Encrypt all files within a directory (with or without recursive mode)
- [x] MVP#5 : Add unit tests
- [ ] MVP#6 : Add init to create a local secret file
- [ ] MVP#7 : crypt file with password and apply to local secret file
- [ ] MVP#8 : make local secret the default provider (rely on MVP#7)
