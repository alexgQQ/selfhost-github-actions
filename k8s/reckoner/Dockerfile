FROM summerwind/actions-runner:latest

# gettext-base contains the envsubst cmd
RUN sudo apt update -y \
  && sudo apt install gettext-base \
  && sudo rm -rf /var/lib/apt/lists/*