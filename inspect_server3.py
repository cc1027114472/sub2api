import paramiko

ssh = paramiko.SSHClient()
ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
ssh.connect('154.23.162.32', port=22, username='root', password='2Q99EAjlnb', timeout=10)

stdin, stdout, stderr = ssh.exec_command('cat /root/sub2api/deploy/build-fork-image.sh')
print(stdout.read().decode('utf-8', errors='ignore'))

ssh.close()
