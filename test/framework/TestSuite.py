import subprocess
import os
import time
from typing import *


class TestSuite:
    def __init__(self, testName):
        '''
            params:\n
            testName: name of test folder in example/
        '''
        self.testName = testName
        self.webhookProcess = None
        self.logFileHandler = None
        self.testPath = os.path.dirname(__file__)
        self.workspacePath = os.path.abspath(self.testPath+"/../..")

    def setUp(self) -> None:
        '''
        set up external webhook server. Logs of output will be put into ../testlog
        '''
        #print("[E2E Test]:%s: setting up webhook server and log"%(self.name))
        
        #setup log file
        logFileName = self.testName+"-"+time.strftime("%Y-%m-%d-%H-%M-%S")+".log"
        self.logFileHandler = open(
            "%s/test/log/%s" % (self.workspacePath, logFileName), "w")
        

        #load model and policy 
        subprocess.Popen("kubectl apply -f %s/example/%s/model.yaml"%(self.workspacePath,self.testName),shell=True,stdout=self.logFileHandler, stderr=self.logFileHandler).wait()
        subprocess.Popen("kubectl apply -f %s/example/%s/policy.yaml"%(self.workspacePath,self.testName),shell=True,stdout=self.logFileHandler, stderr=self.logFileHandler).wait()
        #setup webhook 
        subprocess.Popen("kubectl apply -f %s/config/webhook_external.yaml"%(self.workspacePath),shell=True,stdout=self.logFileHandler, stderr=self.logFileHandler).wait()

        #start the webhook
        cmd = [
            "%s/test/build/main.exe" % (self.workspacePath),
        ]
        self.webhookProcess = subprocess.Popen(
            cmd, cwd=self.workspacePath, stdout=self.logFileHandler, stderr=self.logFileHandler)
        #print("[E2E Test]:%s: admission webhook started, pid %d"%(self.name,self.webhookProcess.pid))
        time.sleep(0.2)

    def tearDown(self) -> None:
        '''
        shut down external webhook server. 
        '''
        #print("[E2E Test]:%s: shutting webhook server and log"%(self.name))
        self.webhookProcess.kill()
        self.logFileHandler.close()

        #remove webhook 
        os.system("kubectl delete -f %s/config/webhook_external.yaml 1> /dev/null 2>/dev/null"%(self.workspacePath))
        os.system("kubectl delete -f %s/example/%s/model.yaml 1> /dev/null 2>/dev/null"%(self.workspacePath,self.testName))
        os.system("kubectl delete -f %s/example/%s/policy.yaml 1> /dev/null 2>/dev/null"%(self.workspacePath,self.testName))
        #remove model and policy 

    def test(self) -> Tuple[int, int]:
        '''
        test each testcase and collect the result\n
        return: Tuple[int,int],in which:\n
            1st value of tuple is the number of passed test\n
            2st value of tuple is the number of failed test\n
        '''
        success = 0
        fail = 0
        testCaseFiles=os.listdir("%s/example/%s/testcase"%(self.workspacePath,self.testName))
        for i in range(0, len(testCaseFiles)):
            time.sleep(0.2)
            yamlFileName = testCaseFiles[i]
            yamlFileAbsoluteName="%s/example/%s/testcase/%s"%(self.workspacePath,self.testName,yamlFileName)
            shouldSuccess = yamlFileName.startswith("approve")

            webhookRunning = self.webhookProcess.poll()
            if webhookRunning != None:
                # webhook server crashed, immediately failed
                print("[E2E Test]:UNTESTED WEBHOOK HAS CRASHED. Test suit %s, Test case %s" % (
                    self.testName, os.path.basename(yamlFileName)))
                fail += 1
                continue
            cmd = [
                "minikube",
                "kubectl",
                "--",
                "apply",
                "-f",
                yamlFileAbsoluteName,
                "--dry-run=server"
            ]
            res = subprocess.Popen(
                cmd, stdout=self.logFileHandler, stderr=self.logFileHandler)
            res.wait()
            # check whether webhook has crashed
            webhookRunning = self.webhookProcess.poll()
            if webhookRunning != None:
                print("[E2E Test]:FAILED WEBHOOK HAS CRASHED. Test suit %s, Test case %s" % (
                    self.testName, os.path.basename(yamlFileName)))
                fail+=1
            elif (shouldSuccess and res.returncode == 0) or (not shouldSuccess and res.returncode != 0):
                # passed
                print("[E2E Test]:PASSED Test suit %s, Test case %s" %
                      (self.testName, os.path.basename(yamlFileName)))
                success += 1
            else:
                # failed
                print("[E2E Test]:FAILED Test suit %s, Test case %s" %
                      (self.testName, os.path.basename(yamlFileName)))
                fail += 1

        return (success, fail)

    def run(self) -> Tuple[int, int]:
        '''
        run whole testsuite and collect the result\n
        return: Tuple[int,int],in which:\n
            1st value of tuple is the number of passed test\n
            2st value of tuple is the number of failed test\n
        '''
        self.setUp()
        res = self.test()
        self.tearDown()
        #print("[E2E Test]:%s: %d test passed, %d test failed"%(self.name,res[0],res[1]))
        return res
