FROM debian:stable-slim
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y
RUN apt-get install apt-transport-https -y
RUN apt-get install apt-utils -y
RUN apt-get install gcc -y
RUN apt-get install g++ -y
RUN apt-get install nano -y
RUN apt-get install tar -y
RUN apt-get install file -y
RUN apt-get install bash -y
RUN apt-get install sudo -y
# RUN apt-get install openssl -y
RUN apt-get install git -y
# RUN apt-get install make -y
RUN apt-get install wget -y
RUN apt-get install curl -y
RUN apt-get install bc -y
# RUN apt-get install pv -y
# RUN apt-get install openssh-server -y
# RUN apt-get install openssh-client -y
# RUN apt-get install python3 -y
# RUN apt-get install python3-pip -y
# RUN apt-get install python3-dev -y
# RUN apt-get install python3-setuptools -y
# RUN apt-get install python -y
RUN apt-get install net-tools -y
RUN apt-get install iproute2 -y
RUN apt-get install iputils-ping -y
# RUN apt-get install golang-go -y

RUN apt-get install nmap -y
ENV TZ="US/Eastern"
ARG USERNAME="morphs"
ARG PASSWORD="asdf"
RUN useradd -m $USERNAME -p $PASSWORD -s "/bin/bash"
RUN mkdir -p /home/$USERNAME
RUN chown -R $USERNAME:$USERNAME /home/$USERNAME
RUN usermod -aG sudo $USERNAME
RUN echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
RUN echo "export PATH=$PATH:/usr/local/go/bin" | tee -a /home/$USERNAME/.bashrc
USER $USERNAME
WORKDIR /home/$USERNAME

# 05JUN2021 , apt installs golang-1.11
# https://golang.org/dl/
COPY ./go_install.sh /home/$USERNAME/go_install.sh
RUN sudo chmod +x /home/$USERNAME/go_install.sh
RUN sudo chown $USERNAME:$USERNAME /home/$USERNAME/go_install.sh
#RUN sudo tar -C /usr/local -xzf $ARCHIVE_NAME
RUN /home/$USERNAME/go_install.sh
# RUN pv /home/$USERNAME/go.tar.gz | sudo tar -C /usr/local -xz
# RUN sudo tar --checkpoint=1 --checkpoint-action=ttyout='%{%Y-%m-%d %H:%M:%S}t (%d sec): #%u, %T%*\r' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
# RUN sudo tar --checkpoint=1 --checkpoint-action=echo -C /usr/local -xzf /home/$USERNAME/go.tar.gz
# RUN /bin/bash -c 'export GO_TAR_KILOBYTES=$(printf "%.3f\n" $(echo "$GO_TAR_BYTES / 1000" | bc -l))'
# RUN /bin/bash -c 'echo $GO_TAR_BYTES'
# RUN /bin/bash -c 'GO_TAR_BYTES=$(stat --format="%s" go.tar.gz) && \
# GO_TAR_KILOBYTES=$(printf "%.3f\n" $(echo "$(stat --format="%s" go.tar.gz) / 1000" | bc -l)) \
# RUN /bin/bash -c 'export GO_TAR_KILOBYTES=$(printf "%.3f\n" $(echo "$(stat --format="%s" go.tar.gz) / 1000" | bc -l))'
# RUN export GO_TAR_KILOBYTES=$(printf "%.3f\n" $(echo "$(stat --format="%s" go.tar.gz) / 1000" | bc -l))
# RUN echo $GO_TAR_KILOBYTES
# RUN sudo tar --checkpoint=1 --checkpoint-action=exec='sh -c "GO_TAR_BYTES=$(stat --format=\"%s\" go.tar.gz) && \
# GO_TAR_KILOBYTES=$(printf "%.3f\n" $(echo "$GO_TAR_BYTES / 1000" | bc -l)) && \
# echo [$TAR_CHECKPOINT] of $GO_TAR_KILOBYTES kilobytes"' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
# RUN sudo tar --checkpoint=1 --checkpoint-action=dot -C /usr/local -xzf /home/$USERNAME/go.tar.gz
# RUN sudo tar --checkpoint=1 --checkpoint-action=exec='/bin/bash -c "cmd=$(echo ZXhwb3J0IEdPX1RBUl9LSUxPQllURVM9JChwcmludGYgIiUuM2ZcbiIgJChlY2hvICIkKHN0YXQgLS1mb3JtYXQ9IiVzIiAvaG9tZS9tb3JwaHMvZ28udGFyLmd6KSAvIDEwMDAiIHwgYmMgLWwpKSAmJiBlY2hvICRHT19UQVJfS0lMT0JZVEVT | base64 -d ; echo); eval $cmd && echo [$TAR_CHECKPOINT] of $GO_TAR_KILOBYTES kilobytes"' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
# base64Encode 'export GO_TAR_KILOBYTES=$(printf "%.3f\n" $(echo "$(stat --format="%s" /home/morphs/go.tar.gz) / 1000" | bc -l)) && echo Extracting [$TAR_CHECKPOINT] of $GO_TAR_KILOBYTES kilobytes /usr/local/go'
# RUN sudo tar --record-size=10K --checkpoint=100 --checkpoint-action=exec='/bin/bash -c "cmd=$(echo ZXhwb3J0IEdPX1RBUl9LSUxPQllURVM9JChwcmludGYgIiUuM2ZcbiIgJChlY2hvICIkKHN0YXQgLS1mb3JtYXQ9IiVzIiAvaG9tZS9tb3JwaHMvZ28udGFyLmd6KSAvIDEwMDAiIHwgYmMgLWwpKSAmJiBlY2hvIEV4dHJhY3RpbmcgWyRUQVJfQ0hFQ0tQT0lOVF0gb2YgJEdPX1RBUl9LSUxPQllURVMga2lsb2J5dGVzIC91c3IvbG9jYWwvZ28= | base64 -d ; echo); eval $cmd"' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
# https://askubuntu.com/questions/1094519/tar-checkpoint-action-exec-write-to-file
RUN sudo tar --checkpoint=100 --checkpoint-action=exec='/bin/bash -c "cmd=$(echo ZXhwb3J0IEdPX1RBUl9LSUxPQllURVM9JChwcmludGYgIiUuM2ZcbiIgJChlY2hvICIkKHN0YXQgLS1mb3JtYXQ9IiVzIiAvaG9tZS9tb3JwaHMvZ28udGFyLmd6KSAvIDEwMDAiIHwgYmMgLWwpKSAmJiBlY2hvIEV4dHJhY3RpbmcgWyRUQVJfQ0hFQ0tQT0lOVF0gb2YgJEdPX1RBUl9LSUxPQllURVMga2lsb2J5dGVzIC91c3IvbG9jYWwvZ28= | base64 -d ; echo); eval $cmd"' -C /usr/local -xzf /home/$USERNAME/go.tar.gz
# RUN git clone https://github.com/0187773933/MACAddressTracker.git
# RUN mkdir -p /home/$USERNAME/MACAddressTracker
COPY . /home/$USERNAME/MACAddressTracker
RUN sudo chown $USERNAME:$USERNAME /home/$USERNAME/MACAddressTracker
WORKDIR /home/$USERNAME/MACAddressTracker

# RUN echo "export PATH=$PATH:/usr/local/go/bin" | sudo tee -a /etc/environment
# RUN /usr/local/go/bin/go version
# RUN /usr/local/go/bin/go clean -cache -modcache -i -r
# RUN /usr/local/go/bin/go get all
RUN /usr/local/go/bin/go build -o macAddressTracker
RUN chmod +x ./macAddressTracker
RUN sudo cp ./macAddressTracker /usr/bin/
RUN mkdir -p ~/.config
RUN mkdir -p ~/.config/personal

# ENTRYPOINT [ "/bin/bash" ]
# RUN /home/$USERNAME/MACAddressTracker/build.sh


# Name too Long
# ENTRYPOINT [ "MAC_LOCATION_NAME=$MAC_LOCATION_NAME \
# MAC_CRON_STRING=$MAC_CRON_STRING MAC_SERVER_PORT=$MAC_SERVER_PORT \
# MAC_SAVED_RECORD_TOTAL=$MAC_SAVED_RECORD_TOTAL \
# MAC_NETWORK_HARDWARE_INTERFACE_NAME=$MAC_NETWORK_HARDWARE_INTERFACE_NAME \
# MAC_REDIS_HOST=$MAC_REDIS_HOST MAC_REDIS_PORT=$MAC_REDIS_PORT MAC_REDIS_DB=$MAC_REDIS_DB \
# MAC_REDIS_PASSWORD=$MAC_REDIS_PASSWORD MAC_REDIS_PREFIX=$MAC_REDIS_PREFIX \
# /usr/bin/macAddressTracker" ]

ENTRYPOINT [ "/usr/bin/macAddressTracker" ]