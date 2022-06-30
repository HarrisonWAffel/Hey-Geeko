#!/bin/bash

echo ""
echo "====> Make sure your hardware is properly setup, otherwise this probably won't work <===="
echo
echo "Ensure that you have a working speaker setup. To test if your speakers work, run the following commands"
echo ---
echo
echo "wget https://www2.cs.uic.edu/~i101/SoundFiles/PinkPanther30.wav # download a sample audio file "
echo "aplay PinkPanther30.wav"
echo
echo
echo "To test your microphone, you can run the following commands (If your speakers don't work, this test won't work)"
echo ---
echo
echo "arecord -D plughw:1 -c1 -r 16000 -f S16_LE test.wav # then say something into the microphone for 5 seconds"
echo "aplay test.wav"
echo
echo "Also make sure you are running python > 3.6 < 3.9, ideally 3.7.x"
echo
read -p "Do you have everything setup? [y/n]:" -n 1 -r
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    clear
    exit 1
fi
echo

BEGIN=$(pwd)
go build -o hey-geeko

sudo apt update
sudo apt-get update

# download deps
sudo apt-get install \
      libasound-dev \
      portaudio19-dev \
      libportaudio2 \
      libportaudiocpp0 \
      swig \
      libatlas-base-dev \
      libpcaudio-dev\
      kramdown \
      ronn \
      libsonic-dev \
      pkg-config \
      libbz2-dev \
      docker \
      jq -y

# You need to make sure that you are working with
# python 3.7.8
# you can run these if you're not:
#
# curl https://pyenv.run | bash
# pyenv download 3.7.8
# pyenv shell 3.7.8

# Solid TTS Library. needs docker
curl -fsSL https://get.docker.com -o docker.sh

## openTTS runs in a always-on server. Without this container TTS wont work!
# this pulls a 4GB(!!) image. However, this provides many voice options and is worth it.
sudo docker run -it -d -p 5500:5500 synesthesiam/opentts:en --cache

# voice2json should be installed in ~
cd

# clone voice-processing repo, voice2json
git clone https://github.com/synesthesiam/voice2json.git

# configure and build voice2json source
cd voice2json || exit
./configure || exit
make
make install
# add command to path
echo "export PATH=\$PATH:"$PWD"/voice2json/bin" >> ~/.bashrc
cd ..

# download english model todo; an env var should change this
PROFILE="en-us_kaldi-rhasspy"

voice2json --debug --profile "$PROFILE" download-profile

# todo; inject custom sentences into sentences.ini
# sentences.ini should be in this repo and just injected at configuration time.
# There should also be a command to update the file when the local version is updated

cd "$BEGIN"

./hey-geeko profile update-sentences || exit 69

# move the custom 'hey-geeko' wake-word model into the correct place
./hey-geeko wake-word --model hey-geeko.pb || exit 69

# train model on sentences.ini
voice2json --debug --profile "$PROFILE" train-profile

cd || exit
source .bashrc


