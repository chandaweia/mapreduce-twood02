# mapreduce tips

## Setting up EC2 Environment

 - Start a new VM with RedHat Enterprise Linux 8
    - We don't suggest using Ubuntu since it comes with an outdated version of Go.
    - Check the lecture video for details! 
    - Be sure to save your keypair!
 - Setup an Elastic IP for your VM so that its IP won't change each time you reboot it 
    - Follow these instructions: https://aws.amazon.com/premiumsupport/knowledge-center/ec2-associate-static-public-ip/
 - Install the VS Code Remote Development extension on your computer
    - Check the website under HW1 for details!
 - Use VS Code to make a remote connection to your VM
    - Use the elastic IP and set the path to your keypair. Your config should be something like:

```
Host distSysVM
    HostName XXX.XXX.XXX.XXX
    IdentityFile ~/.ssh/dist-sys-test.pem
    User ubuntu
```

 - If you have VS Code liveshare extension setup it will give you an error about not being able to install a script to allow browser links. To fix the error run the following in a terminal on your VM:

 ```
wget -O ~/vsls-reqs https://aka.ms/vsls-linux-prereq-script && chmod +x ~/vsls-reqs
sudo vsls-reqs
 ```

 ## Update Go
 If you run `go version` it will show an outdated version of Go. We want at least version 1.13.

```
wget https://golang.org/dl/go1.15.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.2.linux-amd64.tar.gz