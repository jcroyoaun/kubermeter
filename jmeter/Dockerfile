# Use an official Java 11 base image
FROM openjdk:11

# Set environment variables for JMeter version
ENV JMETER_VERSION=5.5 \
    JMETER_HOME=/opt/apache-jmeter-${JMETER_VERSION} \
    PATH=$PATH:/opt/apache-jmeter-${JMETER_VERSION}/bin

# Install prerequisites
RUN apt-get update && \
    apt-get install -y wget tar && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Download and unpack JMeter
RUN wget https://archive.apache.org/dist/jmeter/binaries/apache-jmeter-${JMETER_VERSION}.tgz -P /tmp/ && \
    tar -xzf /tmp/apache-jmeter-${JMETER_VERSION}.tgz -C /opt && \
    rm /tmp/apache-jmeter-${JMETER_VERSION}.tgz

COPY user.properties /opt/apache-jmeter-5.5/bin/

# Expose the default JMeter server port
EXPOSE 1099

# Set jmeter-server as the entry point
ENTRYPOINT ["/opt/apache-jmeter-5.5/bin/jmeter"]
