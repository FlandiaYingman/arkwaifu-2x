# Arkwaifu 2x

Arkwaifu 2x is a CLI tool for Arkwaifu. It supports the following features:

1. Enlarge the Arts with *real-esrgan(realesrgan-x4plus)*.
2. Enlarge the Arts with *real-esrgan(realesrgan-x4plus-anime)*.
3. Re-generate thumbnails for all arts.

This tool is created because the backend of Arkwaifu cannot handle hard works, like enlarging the arts with neutral
networks. So Arkwaifu 2x takes the tasks in local, reducing pressure of the backend.

![screenshot](https://github.com/FlandiaYingman/arkwaifu-2x/assets/9929037/bf23c573-e76e-4a96-9700-6326aa107450)

# Usage

> Prerequisites:
> - Windows 10 or above
> - Python: ^3.11
> - A `USER_TOKEN` to grant permission, see below

This tool is intended to run on Windows. However, it is possible to support other OS like macOS or Linux. I did not
do that because I do not have these environments to test this tool. If you wish to run this tool on other OS, feel free
to open a issue to tell me, and I will try my best to support these OS.

In order to install this tool, run

```commandline
py -m pip install arkwaifu-2x
```

In order to run this tool, set the `USER_TOKEN` first

```cmd
$env:USER_TOKEN="<YOUR USER_TOKEN> # PowerShell
SET USER_TOKEN=<YOUR USER_TOKEN> # CMD
```

Then, you can run the tool

```commandline
arkwaifu-2x
```

# What is `USER_TOKEN`?

`USER_TOKEN` is a token for trusted users to grant permission when they are submitting something to the website. I
designed this because I do not want everyone to be able to submit to the website, and that may cause security issues.
If you want to get a `USER_TOKEN`, just open a issue for that. I will be happy to allocate a `USER_TOKEN` for you. :)

# License

The source code of this project is licensed under the MIT License.

This project utilizes Real-ESRGAN and Real-ESRGAN-ncnn-vulkan, which are respectively licensed under the BSD-3-Clause
license and the MIT License.
