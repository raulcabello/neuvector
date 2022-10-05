FROM docker.io/neuvector/controller:5.0.3

COPY output/go_build_github_com_neuvector_neuvector_controller /controller

RUN cp /usr/local/bin/libpcre.so.3.13.1 /usr/local/bin/libpcre.so.1
ENV LD_LIBRARY_PATH /usr/local/bin:
ENV GOTRACEBACK=all

LABEL name="controller" \
      vendor="NeuVector Inc." \
      version=${NV_TAG} \
      release=${NV_TAG} \
      neuvector.image="neuvector/controller" \
      neuvector.role="controller" \
	  neuvector.rev="git.xxxx"


ENTRYPOINT ["/controller"]
