# `catmeme` : Instead of the `cat` commamd!?
`catmeme` is a joke command that displays ASCII art of a cat meme on the console.
Just like the [sl](https://github.com/mtoyoda/sl) command, all key inputs are disabled while it's running, so you have no choice but to sit back and enjoy the cat meme calmly.

## Preview
The resolution of the cat meme displayed depends on the console's text size.

https://github.com/kmraven/catmeme/assets/128337097/dbed3245-c820-47fa-b8b2-e3dad3a9af32

### `-c` option
Adding the -c option displays the cat meme in color.  

https://github.com/kmraven/catmeme/assets/128337097/060e1872-f15f-44d4-a4b0-2520ea200c17

### `-t` option
When you specify a number after the -t option, the cat meme will be displayed for that duration in seconds. The default display time is 3 seconds.  

https://github.com/kmraven/catmeme/assets/128337097/249011cf-42bc-4a8c-9846-66d36c54ed26

## How to use
Please select the URL for the program that suits your environment from the [release page](https://github.com/kmraven/catmeme/releases) and download it.
For example, if you are using Mac (arm), you can download it using the following command.
```
% curl -LO https://github.com/kmraven/catmeme/releases/download/v0.0.0/catmeme_Darwin_arm64.tar.gz
% tar xvf catmeme_Darwin_arm64.tar.gz
% ./catmeme [options]
```

## For prank makers
To provide a one-time cat meme prank, insert the following script into the target's zshrc file. This script leaves no traces, so the victim won't be able to trace where the cat meme came from.
(This script was made by [Kei](https://github.com/Motifman). Thanks!)
```
function catmeme() {curl -sLo ~/catmeme_Darwin_arm64.tar.gz https://github.com/kmraven/catmeme/releases/download/v0.0.0/catmeme_Darwin_arm64.tar.gz && tar xfz ~/catmeme_Darwin_arm64.tar.gz -C ~/ && ~/catmeme -c -t 10 && sed -i ".aonaon" -e '/#CATMEME/d' ~/.zshrc && rm -f ~/catmeme* ~/.*.aonaon && source ~/.zshrc && unset -f catmeme && unalias cat;} #CATMEME
alias cat="catmeme" #CATMEME
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
