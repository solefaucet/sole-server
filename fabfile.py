from fabric.api import env
from fabric.context_managers import cd
from fabric.operations import run, local, put

env.shell = '/bin/bash -l -c'
env.user = 'd'
env.roledefs.update({
    'staging': ['staging.solebtc.com'],
    'production': ['solebtc.com']
})

# Heaven will execute fab -R staging deploy:branch_name=master
def deploy(branch_name):
    deployProduction(branch_name) if env.roles[0] == 'production' else deployStaging(branch_name)
    
def deployStaging(branch_name):
    printMessage("staging")

    codedir = '$GOPATH/src/github.com/freeusd/solebtc'
    run('rm -rf %s' % codedir)
    run('mkdir -p %s' % codedir)

    local('git archive --format=tar --output=/tmp/archive.tar %s' % branch_name)
    local('ls /tmp')
    put('/tmp/archive.tar', '~/')
    local('rm /tmp/archive.tar')
    run('mv archive.tar %s' % codedir)

    with cd(codedir):
        run('tar xf archive.tar')
        run('go build -o ~/solebtc')

        # database version control
        run("mysql -e 'create database if not exists solebtc_prod';")
        run('go get bitbucket.org/liamstask/goose/cmd/goose')
        run('goose -env production up')

    # restart solebtc service with supervisorctl
    run('supervisorctl restart solebtc')

def deployProduction(branch_name):
    printMessage("production")

    # TODO
    # scp executable file from staging to production, database up, restart service
    # mark current timestamp or commit as version number so we can rollback easily

def printMessage(server):
    print("Deploying to %s server at %s as %s" % (server, env.host, env.user))
