FROM centos

RUN yum install -y wget
RUN wget -O httpcase.tar.gz "https://github.com/roseboy/httpcase/releases/download/v1.0.9-beta/httpcase_1.0.9-beta_linux_x86_64.tar.gz"
RUN tar --remove-files -C /usr/local/bin -zxf httpcase.tar.gz