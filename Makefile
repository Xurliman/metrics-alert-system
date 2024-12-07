agent_dir = cmd/agent/agent
server_dir = cmd/server/server
dsn = 'host=localhost port=5432 user=postgres password=kali dbname=metrics sslmode=disable'

test1:
	cd ~/Projects/metrics/cmd/agent && go build -o agent . && \
    cd ~/Projects/metrics/cmd/server && go build -o server . && \
    cd ~/Projects/metrics && \
    ./metricstest -test.v -test.run=^TestIteration2$ \
    -binary-path=${server_dir} -server-port=8088 -agent-binary-path=${agent_dir} -source-path=. -file-storage-path='host=localhost port=5432 user=postgres password=kali dbname=metrics sslmode=disable'

test2:
	cd ~/Projects/metrics/cmd/agent && go build -o agent . && \
	cd ~/Projects/metrics/cmd/server && go build -o server . && \
	cd ~/Projects/metrics && \
	./metricstest -test.v -test.run=^TestIteration1$ \
	-binary-path=${server_dir} -server-port=8088 -agent-binary-path=${agent_dir} -source-path=. -file-storage-path=${dsn}

#https://makefiletutorial.com/