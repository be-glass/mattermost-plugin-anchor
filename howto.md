autopilot token on debug server: d5abgut31ifb7gkxx3a96oes3a

cd /home/boris/GolandProjects/mattermost-plugin-development/server
sudo service mysql stop
make run-server

cd /home/boris/GolandProjects/mattermost-plugin-development/webapp
make run

http://localhost:8065/
