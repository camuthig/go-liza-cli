# Liza CLI

A commandline tool for interacting with BitBucket pull request updates. The main goal of the project is to provide a way for developers to quickly see what PRs they should take a look at based on new comments or new code without needing to use the BitBucket email notifications.

# App Password

To use Liza CLI you will need an app password configured in BitBucket. This password must have the following permissions

* Account read
* Repositories read
* Pull requests read

# Installation

## Pre-built Binary

Pre-build binaries can be found for most architectures of Windows, Linux, and MacOS. Visit the [releases](https://github.com/camuthig/go-liza-cli/releases)
page to find the latest release and download the binary.

## Building Manually

Requires Go 1.17.0+

1. Download the project from Github
1. Build the project using Golang

    `go build`
1. Symlink the generated `lizacli` binary into your `PATH`

    `ln -s lizacli <location on your path>`
1. Add your credentials

    `liza credentials`
2. Add the `update` command to your crontab to continuously pull latest updates from BitBucket

    `* * * * * liza update`

# Adding to oh-my-zsh Prompt

The way I use this tool currently is by having a notification appear in my oh-my-zsh prompt any time there are new notifications for me. This can be done using the following script.

```bash
# Format for BitBucket updates alert
parse_bb_alerts () {
        if [[ ! $(command -v liza) ]]
        then
            echo ""
            exit
        fi

        # This can be skipped by most users. It is helpful if you are symlinking and developing the project
        # and need a way to disable the prompt when bugs arise.
        if [ "$LIZA_ENABLED" != true ]
        then
            echo ""
            exit
        fi

        local COUNT

        COUNT=$(liza updates --count)

        if [[ $COUNT > 0 ]]
        then
                echo "%{$BLUE%} [BB]"
        else
                echo ""
        fi
}

# Add $(parse_bb_alerts) into your PROMPT property to include the `[BB]` icon whenever you have new updates.
```