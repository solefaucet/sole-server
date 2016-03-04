from fabric.api import env
from fabric.context_managers import cd
from fabric.operations import run

env.shell = '/bin/bash -l -c'
env.user = 'd'
env.roledefs.update({
    'staging': ['staging.solebtc.com'],
    'production': ['solebtc.com']
})

# Heaven will execute fab -R staging deploy:branch_name=master
def deploy(branch_name):
    print("Executing on %s as %s" % (env.host, env.user))
    
    run('rm -rf $GOPATH/src/github.com/freeusd/solebtc')
    run('go get -u github.com/freeusd/solebtc')
    run('go get -u bitbucket.org/liamstask/goose/cmd/goose')
    with cd('$GOPATH/src/github.com/freeusd/solebtc'):
        run('goose up')
        run('go install')
        run('supervisorctl restart solebtc')
