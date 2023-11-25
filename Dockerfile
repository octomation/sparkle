FROM scratch

LABEL author="OctoLab team <team@octolab.org>"
LABEL org.opencontainers.image.title="✨ Sparkle service"
LABEL org.opencontainers.image.description="The personal development framework and Personal Knowledge Management platform."
LABEL org.opencontainers.image.source="https://github.com/withsparkle/service"
LABEL org.opencontainers.image.licenses="AGPL-3.0-later"

COPY sparkle /sparkle
EXPOSE 3360 8080 8081 8890 8891

ENTRYPOINT ["/sparkle"]
CMD ["server", "run"]
