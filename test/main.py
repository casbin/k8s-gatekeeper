import os
from framework.TestSuite import TestSuite
def main():
    testPath = os.path.dirname(__file__)
    workspacePath = os.path.abspath(testPath+"/..")
    passed=0
    failed=0
    dirs=os.listdir(workspacePath+"/example")
    for i in range(0,len(dirs)):
        testsuit=TestSuite(dirs[i])
        res=testsuit.run()
        passed+=res[0]
        failed+=res[1]
    print("[E2E Test]: %d passed, %d failed"%(passed,failed))
    if failed!=0:
        exit(-1)

if __name__ == '__main__':
    main()