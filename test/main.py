# Copyright 2022 The casbin Authors. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
