#!/bin/bash

# function to check if a program is installed
function check_installed {
    if ! command -v "$1" &> /dev/null
    then
        # return false if the program isn't found
        return 1
    else
        # return true if the program is found
        return 0
    fi
}

# function to install a package
function install_package {
    echo "updating your package lists... this might take a sec!"
    sudo apt-get update
    echo "installing $1, hang tight..."
    sudo apt-get install -y "$1"
}

# check if our programming language, go, is installed
if check_installed "go"; then
    echo "yay! go is already installed."
else
    echo "uh-oh, looks like go isn't installed. let's fix that!"
    install_package "golang"
fi

# check if the mighty media swiss-army knife (ffmpeg) is installed
if check_installed "ffmpeg"; then
    echo "awesome! ffmpeg is already installed."
else
    echo "looks like ffmpeg is missing. installing it for you..."
    install_package "ffmpeg"
fi

# check for git
if check_installed "git"; then                                                                   
    echo "awesome! git is already installed."                                                    
else                                                                                                
    echo "looks like git is missing. installing it for you..."                                   
    install_package "git"

# verify that the current working directory is writable and usable
if [ ! -w "$PWD" ]; then
    echo "whoops! you don't have write permission for the current directory: $PWD"
    # no permission, no fun, no script execution
    exit 1
else
    echo "perfect! you can write here. let's move forward!"
fi

# clone the go program from github
if [ -d "auto_compressor" ]; then
    echo "hmm... looks like 'auto_compressor' directory already exists. skipping the clone."
else
    echo "cloning the auto_compressor repository from github..."
    git clone https://github.com/aetherbird/auto_compressor.git
    if [ $? -ne 0 ]; then
        echo "uh-oh, something went wrong while cloning the repository. try checking your internet connection!"
        exit 1  # since cloning failed, we exit
    fi
fi

echo
echo
echo
cat << "EOF"
          .--.
         |o_o |
         |:_/ |
        //   \ \
       (|     | )
      /'\_   _/`\
      \___)=(___/
EOF
echo
echo
echo
# wrap-up: instructions for running the go program!
echo "auto_compressor is ready"
echo "Here's how you can get started:"
echo "1. Navigate into the 'auto_compressor' directory:"
echo "   cd auto_compressor"
echo "2. Run the go program using the following command:"
echo "   go run auto_compressor.go <input_file.mp4> <desired_output_size_in_MB>"
echo "I.E:
echo "   go run auto_compressor.go ~/Videos/my_cool_video.mp4 300
echo "...and that's it! Happy compressing!"

