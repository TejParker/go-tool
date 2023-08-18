package template

var (
	Filebeat = `
# filebeat download address
# https://www.elastic.co/cn/downloads/past-releases/filebeat-7-9-3/
# input
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - ./*.log
# output
output.logstash:
  hosts: ["localhost:5044"]
`
)
