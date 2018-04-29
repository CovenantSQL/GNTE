FROM ubuntu
RUN apt-get update
RUN apt-get install -y build-essential
RUN apt-get install -y git wget iftop iproute2
RUN apt-get install -y netcat-openbsd tmux screen dstat mtr net-tools sendip tcpreplay netperf iperf iperf3 fping iputils-ping
RUN apt-get install -y tcpdump tcsh iptraf
RUN echo 'bind "\C-n":history-search-backward' >> ~/.bashrc
