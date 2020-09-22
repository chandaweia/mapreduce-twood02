# mapreduce tips

## Setting up EC2 Environment

 - Start a new VM with RedHat Enterprise Linux 8
    - We don't suggest using Ubuntu 18.04 since it comes with an outdated version of Go.
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
    User ec2-user
    # or User ubuntu if using Ubuntu 20
```

Use the terminal in VS Code to complete the following setup steps:


### Install go / dev tools

On RedHat we can install Go and lots of other languages/tools using:

```
sudo yum groupinstall -y 'Development Tools'
sudo yum install -y go 
```

Be sure to configure your GOPATH settings as suggested in the HW2 docs.

### Clone your repo:

Now you can get your repo and use Open Folder in VS Code to display it in your workspace:

```
git clone clone https://github.com/gwDistSys20/YOUR_REPO_URL
```

### Test VS Code Go integrations
Open a file like `src/main/mrmaster.go` so VS Code will ask you to install the Go Analysis tools. Once the tools are installed you should be able to navigate the codebase easily, such as right clicking on `mr.MakeMaster(...)` and selecting Go To Definition.

> **WARNING:** Since your `$GOPATH` is the root of your repo, VS Code might install the packages for all of the analysis tools into your workspace! Be careful not to commit all of these files to your repository.  You should setup a `.gitignore` file to appropriately ignore all of these packages and other temporary files produced by running the code.

## Test the sequential Map Reduce

```
cd hw2-mapreduce-your-team-name
source ./setEnv.sh
cd src/main
go build -buildmode=plugin ../mrapps/wc.go
rm mr-out*
go run mrsequential.go wc.so pg*.txt
head mr-out-0

# should give output like:
A 509
ABOUT 2
ACT 8
...
```

## Test the Master and Worker
To test that the master and worker are running correctly, we can start them in two separate terminals.

The code in this repository has been modified so that the master will start its RPC server and the worker will  issue two RPC requests. First it will use the simple Example RPC, and then it will register itself as a worker and get back a list of all the files that need to be processed.

```
# Terminal 1 -- start first
source setenv.sh
cd src/main
go build -buildmode=plugin ../mrapps/wc.go
go run mrmaster.go pg*.txt

# Terminal 2 -- start after master
source setenv.sh
cd src/main
go run mrworker.go wc.so

# worker output:
reply.Y 100
Got a list of files [pg-being_ernest.txt pg-dorian_gray.txt pg-frankenstein.txt pg-grimm.txt pg-huckleberry_finn.txt pg-metamorphosis.txt pg-sherlock_holmes.txt pg-tom_sawyer.txt]

```

If you get errors about loading plugins in the worker process, then you probably modified a file in `mr/` and didn't recompile!