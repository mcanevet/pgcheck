FROM scratch
EXPOSE 8080
ADD pgcheck /
CMD ["/pgcheck"]
