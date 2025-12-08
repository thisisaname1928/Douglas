import webview
import zipfile
import time
import subprocess
import sys
import requests
import os
import json

host = "localhost"
port = 8080

def update(path):
    try:
        f = zipfile.ZipFile(path, 'r')
        flist = f.namelist()
        for i in flist:
            f.extract(i, path="./")
        return True
    except:
        return False
    
def check4Update():
    try:
        f = open("./appVersion.json", 'r', encoding='utf-8')
        dat = json.load(f)

        if dat["shouldUpdate"]:
            print("updatingg...")
            r = update('update.zip')
            if r:
                dat["shouldUpdate"] = False
                r = json.dumps(dat)
                f = open("./appVersion.json", "w", encoding="utf-8")
                f.write(r)
                f.close()
            print("ok")
    except:
        print("not ok")

if __name__ == '__main__':
    check4Update()
    
    os._exit(0)
    proccess = subprocess.Popen(["goParsingDocx.exe"])

    while True:
        try:
            response = requests.get(f'http://{host}:{port}/check')
            if response.status_code == 200:
                break
        except:
            time.sleep(0.5)
    
    webview.settings['ALLOW_DOWNLOADS'] = True
    webview.create_window(title="Douglas", url=f'http://{host}:{port}/Home', width=1424, height=700)

    webview.start(icon="./app/icon.ico")

    if sys.platform == "win32":
        proccess.terminate()
    else:
        proccess.kill()

    proccess.wait()