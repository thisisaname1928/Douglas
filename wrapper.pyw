import webview
import socket
import time
import subprocess
import sys
import requests

host = "localhost"
port = 8080

if __name__ == '__main__':
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