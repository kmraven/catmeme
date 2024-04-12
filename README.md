# `catmeme` : Instead of `cat` commamd!?
`catmeme` is a joke command to display ASCII art of cat meme on console.
Like the [sl](https://github.com/mtoyoda/sl) command, any keystrokes are disabled during execution, so you must settle down to watch the cat meme.

## Preview
https://github.com/kmraven/catmeme/assets/128337097/dbed3245-c820-47fa-b8b2-e3dad3a9af32

### `-c` option
With the -c option, the cat meme is displayed in color.  

https://github.com/kmraven/catmeme/assets/128337097/060e1872-f15f-44d4-a4b0-2520ea200c17

### `-t` option
By specifying a number after the -t option, the cat meme is displayed for that number of seconds.
The default display time is 3 seconds.  

https://github.com/kmraven/catmeme/assets/128337097/249011cf-42bc-4a8c-9846-66d36c54ed26

## How to use
From the [release page](https://github.com/kmraven/catmeme/releases), select the URL of the program that matches your environment and download it.
For example, Mac(arm) user can download with the following commands.
```
% curl -LO https://github.com/kmraven/catmeme/releases/download/v0.0.0/catmeme_Darwin_arm64.tar.gz
% tar xvf catmeme_Darwin_arm64.tar.gz
% ./catmeme [options]
```

## Development Notes
#### To convert video to frames
```
% python --version
Python 3.9.16
% python -m venv venv
% source ./venv/bin/activate
% pip install --upgrade pip && pip install -r requirements.txt
% python videoConv.py [video_files]
% deactivate
```

#### To release
```
% export GITHUB_TOKEN="<MY_TOKEN_HERE>"
% git merge develop
% git tag <NEW_TAG_HERE>
% git push origin main
% git push origin main --tags
% goreleaser --rm-dist
```
