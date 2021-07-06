#!/bin/bash

exit_on_error() {
    exit_code=$1
    text=$2
    last_command=${@:2}
    if [ $exit_code -ne 0 ]; then
        >&2 echo -e "$text."
        exit $exit_code
    fi
}

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
BANG_PATH="$HOME/.bang/bin/@"
LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/christiangelone/bang/releases/latest)
LATEST_VERSION=$(echo $LATEST_RELEASE | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
BIN_URL="https://github.com/christiangelone/bang/releases/download/$LATEST_VERSION/bang_${OS}_${ARCH}"

echo -e "⬇️  Downloading bang."
curl -fsL -H 'Accept: application/octet-stream' $BIN_URL -o bang
exit_on_error $? "💀  Failed to download bang bootstrap" !!
chmod +x bang
exit_on_error $? "💀  Failed bootstrap executable" !!

echo -e "⚙️  Installing bang."
./bang install github.com/christiangelone/bang >/dev/null 2>&1
exit_on_error $? "💀  Failed to install bang" !!
rm ./bang
exit_on_error $? "💀  Error removing bootstrap" !!

if [[ "$PATH" != *"$BANG_PATH"* ]]; then
  case "$SHELL" in
  */bash*)
    if [[ -r "$HOME/.bash_profile" ]]; then
      SHELL_PROFILE="$HOME/.bash_profile"
    else
      SHELL_PROFILE="$HOME/.bashrc"
    fi
    ;;
  */zsh*)
    SHELL_PROFILE="$HOME/.zshrc"
    ;;
  *)
    SHELL_PROFILE="$HOME/.profile"
    ;;
  esac

  echo -e "✏️  Adding bang to your PATH in $SHELL_PROFILE"
  echo -e "    \033[0;34m•\033[0m Please run \033[1;33msource $SHELL_PROFILE\033[0m"
  echo -e "\nBANG_PATH=$BANG_PATH" >> $SHELL_PROFILE
  echo 'export PATH="$PATH:$BANG_PATH"' >> $SHELL_PROFILE
fi

echo -e "🎉 Bang installed!"
